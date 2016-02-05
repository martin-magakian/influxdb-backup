package sqlite

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/efigence/influxdb-backup/common"
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
	var sql SQLiteOut
	sql.Init("t-data/router-test")

	err := sql.workers.RunWriter(writeCh, []string{`t-data`,`point-writer-test.sqlite`}, false)
	Convey("WritePoint", t, func() {
		So(err,ShouldBeNil)
	})
	close(writeCh) // close channel so writer exits

	sql.workers.Shutdown()
}

func TestQuoting(t *testing.T) {
	var f common.Field
	f.Name  = "bad-name-test"
	f.Values = make(map[string]interface{})
	f.Values[`time`] = time.Now()
	f.Values[`cake`] = "lie"
	f.Values[`some-long-name`] = time.Now()
	f.Values[`name with spaces`] = time.Now()
	f.Values[`ga@#&*$H*&GD&*!@GTbage`] = time.Now()

	writeCh := make(chan *common.Field, 10)
	writeCh <- &f
	writeCh <- &f
	writeCh <- &f
	var sql SQLiteOut
	sql.Init("t-data/quote-test")

	err := sql.workers.RunWriter(writeCh, []string{`t-data`,`quoted-writer-test.sqlite`}, false)
	Convey("WriteGarbageFields", t, func() {
		So(err,ShouldBeNil)
	})
	close(writeCh) // close channel so writer exits
	sql.workers.Shutdown()
}
