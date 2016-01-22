package sqlite

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSQLite(t *testing.T) {
	sql, err := NewSQLite(`t-data/sqlite`)
	Convey("Create DB", t, func() {
		So(err, ShouldEqual, nil)
	})
	err = sql.SaveSeriesList([]string{"test"})
	Convey("SaveSeriesList", t, func() {
		So(err, ShouldEqual, nil)
	})
}
