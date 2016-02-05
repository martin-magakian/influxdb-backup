package input

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)


func TestDummy(t *testing.T) {
	// FIXME testdb with pregenerated data
	i, err := NewDummy("http://localhost:8086", "root","root","stats")
	series, err := i.GetSeriesList()
	Convey("GetSeriesList", t, func() {
		So(err, ShouldEqual, nil)
		So(len(series), ShouldBeGreaterThan,0)
	})
	tsStart := time.Now().Add(-1 * time.Hour)
	points, err := i.GetFieldRangeByName(series[0],tsStart,time.Now())
	Convey("GetFirstPoint", t, func() {
		So(err, ShouldEqual, nil)
		So(points, ShouldNotEqual, nil)
		So(points[0].Name, ShouldNotEqual,"")
		So(points[0].Values["time"], ShouldBeGreaterThan,0)

	})
}
