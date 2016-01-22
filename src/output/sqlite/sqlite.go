package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"path/filepath"
	"os"
	"crypto/sha256"
	"fmt"
	"strings"
	"github.com/efigence/influxdb-backup/src/common"

)

type SQLiteOut struct {
	path string

}

func New(args []string) (common.Output,error) {
	return NewSQLite(args[0])
}

func NewSQLite(path string) (common.Output,error) {
	var err error
	var out SQLiteOut
	var mode os.FileMode
	mode = 0744
	os.MkdirAll(path, mode)
	out.path = path
	return &out,err

}

func (out *SQLiteOut) SaveSeriesList(series []string) (err error) {
	s, err := sql.Open("sqlite3", filepath.Join(out.path, "series.sqlite"))
	if (err != nil ) { return err }
	_, err = s.Exec("CREATE TABLE IF NOT EXISTS series( name text UNIQUE , file text )")
	if (err != nil ) { return err }
	_, err = s.Exec("BEGIN")
	if (err != nil ) { return err }
	for _, name := range series {
		_, err = s.Exec("INSERT OR IGNORE INTO series(name, file) VALUES(?,?)",name,SeriesNameGen(name))
		if (err != nil ) { return err }
	}
	_, err = s.Exec("COMMIT")
	if (err != nil ) { return err }
	return err
}


// Generate shortened series name
// does not have to be unique, just unique enough that tens of thousands of series wont land in same sqlite DB
func SeriesNameGen(seriesName string) string {
	s := strings.Split(seriesName, `.`)
	prefix := ``
	if len(s[0]) > 5 {
		prefix = s[0][:5]
	} else {
		prefix = s[0]
	}
	if len(s) > 3 {
		if len(s[1]) > 5 {
			prefix = prefix + `.` + s[1][:5]
		} else {
			prefix = prefix + `.` + s[1]

		}
	}
	hash := sha256.Sum256([]byte(seriesName))
	// 256 buckets + name-based prefix should be fine for most users
	return prefix + fmt.Sprintf("-%x", hash)[:3]
}
