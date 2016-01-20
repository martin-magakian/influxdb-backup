package output

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"path/filepath"
	"os"

)

type SQLiteOut struct {
	path string
}

func NewSQLite(path string) (Output,error) {
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
	_, err = s.Exec("CREATE TABLE IF NOT EXISTS series( name text UNIQUE )")
	if (err != nil ) { return err }
	_, err = s.Exec("BEGIN")
	if (err != nil ) { return err }
	for _, name := range series {
		s.Exec("INSERT OR IGNORE INTO series(name) VALUES(?)",name)
	}
	_, err = s.Exec("COMMIT")
	if (err != nil ) { return err }
	return err
}


// Generate shortened series name
// does not have to be unique, just unique enough that thousands of series wont land in same sqlite DB
func SeriesNameGen(seriesName string) string {
	return `s`
}
