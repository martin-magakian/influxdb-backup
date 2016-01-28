package sqlite

import (
	"sync"
	"github.com/efigence/influxdb-backup/src/common"
	"strings"
	"fmt"
)

type writers struct {
	sync.RWMutex
	nosync bool
	path string
	writeCh map[string]chan *common.Field

}


// once started writer's path/sync mode cant be changed, copy any relevant parameters at the creation time
func (s *SQLiteOut) newWriter() (w writers) {
	w.nosync = s.nosync
	w.path = s.path
	return w
}


func (w *writers) NewChannel(name string) (chan *common.Field, error) {
	var err error
	w.Lock()
	defer w.Unlock()
	// request for new one itself is unsynchronized so if it already exists, just give it to whoever asked
	if ch, ok := w.writeCh[name]; ok {
		return ch,err
	} else {
		w.writeCh[name] = make(chan *common.Field,16384)
		path := []string{w.path, name + `.sqlite` }
		err := RunWriter(w.writeCh[name],path,w.nosync)
		if (err == nil) {
			return w.writeCh[name],err
		} else {
			delete (w.writeCh,name)
			return make(chan *common.Field),err
		}
	}
}

func (w *writers) GetRouteFor(r string) (chan *common.Field,error) {
	var err error
	w.RLock()
	if ch, ok := w.writeCh[r]; ok {
		w.RUnlock()
		return ch,err
	} else {
		w.RUnlock()
		return w.NewChannel(r)
	}
}

func RunWriter (req chan *common.Field, path []string,nosync bool) (err error) {
	s, err := sqliteOpen(path ,nosync)
	if err != nil {return err}
	for field := range req {
		tableName := quoteTableName(field.Name)
		l := len(field.Values)
		keys := make([]string,l,l)
		values := make([]interface{},l,l)
		params := make([]string,l,l)
		i:=0
		for k,v := range field.Values {
			keys[i] = k
			values[i] = v
			params[i] = `?`
			i++
		}
		query := "INSERT INTO " +
			tableName +
			"(" + strings.Join(keys,`,`) + ")" +
			" VALUES (" + strings.Join(params,`,`) + ")"

		_, err := s.Exec(query, values...)
		// I haven't found way to directly extract SQLite errors so we will have to rely on error strings;/
		// dynamically create tables
		if err != nil && strings.Contains(err.Error(), "no such table") {
			_, err = s.Exec("CREATE TABLE " + tableName +"( time INTEGER );")
			_, err = s.Exec(query, values...)
		}
		// dynamically create fields...
		if err != nil && strings.Contains(err.Error(), "no column named") {
			out := strings.Split(err.Error(), "has no column named ")
			if (len(out) < 0) {panic(fmt.Sprintf("cant parse error %s", err)) }
			for(strings.Contains(err.Error(), "no column named")) {
				_, err = s.Exec("ALTER TABLE " + tableName + " ADD COLUMN " + out[1] + " BLOB")
				_, err = s.Exec(query, values...)
			}
		}
		if err != nil {
			panic(fmt.Sprintf("%+v",err))
		}

	}
	_, err = s.Exec("PRAGMA  synchronous = FULL")
	return err
}
