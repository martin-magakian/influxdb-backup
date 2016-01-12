package input

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"fmt"
)

func TestInflux09(t *testing.T) {
	i, err := NewInflux09("http://localhost:8086", "root","root")
	if err != nil {
		SkipConvey("Influx09 GetSeriesList" + fmt.Sprintf("%s",err), t, func() {})
		return
	}
	series, err := i.GetSeriesList()
	Convey("GetSeriesList", t, func() {
		So(err, ShouldEqual, nil)
		So(len(series), ShouldBeGreaterThan,0)
	})
}
