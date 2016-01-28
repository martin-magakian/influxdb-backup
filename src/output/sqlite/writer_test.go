package sqlite

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/efigence/influxdb-backup/src/common"
	"testing"
	"time"
)

func TestWriter(t *testing.T) {
	var f common.Field
	f.Name  = "zupa"
	f.Values = make(map[string]interface{})
	f.Values[`time`] = time.Now()
	f.Values[`cake`] = "lie"
	f.Values[`other`] = time.Now()

	writeCh := make(chan *common.Field, 1)
	writeCh <- &f
	close(writeCh) // close channel so writer exits
	err := RunWriter(writeCh, []string{`t-data`,`point-writer-test.sqlite`}, false)
	Convey("WritePoint", t, func() {
		So(err,ShouldBeNil)
	})
}
