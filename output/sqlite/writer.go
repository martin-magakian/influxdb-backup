package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/efigence/influxdb-backup/common"
	"strings"
	"sync"
	"sync/atomic"
)

type dbWriters struct {
	sync.RWMutex
	nosync   bool
	path     string
	writeCh  map[string]chan *common.Field
	shutdown sync.WaitGroup
	writes   uint64
}

// once started writer's path/sync mode cant be changed, copy any relevant parameters at the creation time
func (s *SQLiteOut) newWriter() (w dbWriters) {
	w.nosync = s.nosync
	w.path = s.path
	w.writeCh = make(map[string]chan *common.Field)
	return w
}

func (w *dbWriters) NewChannel(name string) (chan *common.Field, error) {
	var err error
	w.Lock()
	defer w.Unlock()
	// request for new one itself is unsynchronized so if it already exists, just give it to whoever asked
	if ch, ok := w.writeCh[name]; ok {
		return ch, err
	} else {
		log.Debug(`creating channel for DB, key %s`, name)
		w.writeCh[name] = make(chan *common.Field, 16)
		path := []string{w.path, name + `.sqlite`}
		err := w.RunWriter(w.writeCh[name], path, w.nosync)
		if err == nil {
			return w.writeCh[name], err
		} else {
			delete(w.writeCh, name)
			return make(chan *common.Field), err
		}
	}
}

func (w *dbWriters) GetRouteFor(r string) (chan *common.Field, error) {
	var err error
	w.RLock()
	if ch, ok := w.writeCh[r]; ok {
		w.RUnlock()
		return ch, err
	} else {
		w.RUnlock()
		log.Debug("Creating route for %s", r)
		return w.NewChannel(r)
	}
}
func (w *dbWriters) Shutdown() {
	w.Lock()
	defer w.Unlock()
	log.Debug("sending stop signal to workers")
	for _, ch := range w.writeCh {
		close(ch)
	}
	log.Debug("waiting for workers to finish")
	w.shutdown.Wait()
}

func (w *dbWriters) RunWriter(req chan *common.Field, path []string, nosync bool) (err error) {
	// short-circuit if error
	db, err := sqliteOpen(path, nosync)
	log.Debug("Running writer for %+v", path)
	if err != nil {
		return err
	}
	w.shutdown.Add(1)
	go func() {
		defer w.shutdown.Done()
		// shutdown indicator
		WriterLoop(db, req, &w.writes)
		//cleanup
		log.Debug("writer for %+v finished, flushing", path)
		_, err = db.Exec("PRAGMA  synchronous = FULL")
		log.Debug("writer for %+v exiting", path)
		db.Close()

	}()
	return err
}

func WriterLoop(db *sql.DB, req chan *common.Field, reqCtr *uint64) {
	iter := 0
	db.Exec("Begin")
	for field := range req {
		tableName := quoteTableName(field.Name)
		l := len(field.Values)
		keys := make([]string, l, l)
		values := make([]interface{}, l, l)
		params := make([]string, l, l)
		i := 0
		for k, v := range field.Values {
			keys[i] = k
			values[i] = v
			params[i] = `?`
			i++
		}
		query := "INSERT INTO " +
			tableName +
			"(" + strings.Join(keys, `,`) + ")" +
			" VALUES (" + strings.Join(params, `,`) + ")"

		_, err := db.Exec(query, values...)
		// I haven't found way to directly extract SQLite errors so we will have to rely on error strings;/
		// dynamically create tables
		if err != nil && strings.Contains(err.Error(), "no such table") {
			_, err = db.Exec(`CREATE TABLE "` + tableName + `" ( time INTEGER );`)
			_, err = db.Exec(query, values...)
		}
		// dynamically create fields...
		for err != nil && strings.Contains(err.Error(), "no column named") {
			out := strings.Split(err.Error(), "has no column named ")
			_, err = db.Exec(`ALTER TABLE "` + tableName + `" ADD COLUMN "` + out[1] + `" BLOB`)
			_, err = db.Exec(query, values...)
			log.Debug("ERR: %+v", err)
		}
		if err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
		iter++
		atomic.AddUint64(reqCtr, 1)
		if (iter % 10000) == 9999 {
			_, err = db.Exec("COMMIT")
		}
	}
	db.Exec("COMMIT")
	log.Debug("Writer extited after %d iterations", iter)
}
