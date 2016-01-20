package output

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSQLite(t *testing.T) {
	sql, err := NewSQLite(`t-data/sqlite`)
	Convey("N", t, func() {
		So(err, ShouldEqual, nil)
	})
	err = sql.SaveSeriesList([]string{"test"})
	Convey("N", t, func() {
		So(err, ShouldEqual, nil)
	})
}
