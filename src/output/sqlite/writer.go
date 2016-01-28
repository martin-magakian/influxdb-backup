package sqlite

import (
	"sync"
	"github.com/efigence/influxdb-backup/src/common"
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
		err := RunWriter(w.writeCh[name],w.path)
		if (err == nil) {
			return w.writeCh[name],err
		} else {
			delete (w.writeCh,name)
			return make(chan *common.Field),err
		}
	}


}

func RunWriter (req chan *common.Field, path string) (err error) {
	return err
}
