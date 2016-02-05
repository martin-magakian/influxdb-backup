package sqlite

import (
	"github.com/efigence/influxdb-backup/common"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestRouter(t *testing.T) {
	var f common.Field
	var sql SQLiteOut
	sql.Init("t-data/router-test")
	f.Name = "zupad"
	f.Values = make(map[string]interface{})
	f.Values[`time`] = time.Now()
	f.Values[`data`] = "lie"
	f.Values[`other_thing`] = time.Now()
	Convey("Create writer", t, func() {
		So(sql.writers, ShouldNotBeNil)
	})
	ch := make(chan *common.Field, 128)
	sql.routers.Add(1)
	go sql.route(ch)
	ch <- &f
	ch <- &f
	ch <- &f
	ch <- &f
	close(ch)
	sql.writers.Shutdown()
}
