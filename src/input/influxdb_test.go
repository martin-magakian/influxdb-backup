package input

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
	"fmt"
//	"strconv"
)


func TestInflux09(t *testing.T) {
	// FIXME testdb with pregenerated data
	i, err := NewInflux09("http://localhost:8086", "root","root","_internal")
	if err != nil {
		SkipConvey("Influx09 GetSeriesList" + fmt.Sprintf("%s",err), t, func() {})
		return
	}
	series, err := i.GetSeriesList()
	Convey("GetSeriesList", t, func() {
		So(err, ShouldEqual, nil)
		So(len(series), ShouldBeGreaterThan,0)
	})
	tsStart := time.Now().Add(-100000 * time.Hour)
	points, err := i.GetFieldRangeByName(series[0],tsStart,time.Now())
	Convey("GetFirstPoint", t, func() {
		So(err, ShouldEqual, nil)
		So(points, ShouldNotEqual, nil)
		So(points[0].Name, ShouldNotEqual,"")
		So(points[0].Values["time"], ShouldBeGreaterThan,0)
	})
}
