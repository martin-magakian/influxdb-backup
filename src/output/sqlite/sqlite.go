package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"path/filepath"
	"os"
	"crypto/sha256"
	"hash"
	"errors"
	"fmt"
	"strings"
	"github.com/efigence/influxdb-backup/src/common"

)

type SQLiteOut struct {
	path string
	spine int // number of directories
	leaf int // number of files per dir
	hash hash.Hash
	nosync bool

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
	if out.nosync {
		_, err = s.Exec("PRAGMA  synchronous = 0")
		if (err != nil ) { return err }
	}
	_, err = s.Exec("CREATE TABLE IF NOT EXISTS series( name TEXT UNIQUE , file TEXT )")
	if (err != nil ) { return err }
	_, err = s.Exec("BEGIN")
	if (err != nil ) { return err }
	for _, name := range series {
		_, err = s.Exec("INSERT OR IGNORE INTO series(name, file) VALUES(?,?)",name,SeriesNameGen(name))
		if (err != nil ) { return err }
	}
	_, err = s.Exec("COMMIT")
	if (err != nil ) { return err }
	err = s.Close()

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

func (out *SQLiteOut) SaveFields(prefix string) (error) {
	path := filepath.Join(out.path, prefix + ".sqlite")
	s, err := sql.Open("sqlite3", path)
	if (err != nil ) { return errors.New("open: " + err.Error() + "; dbfile: " + path) }
	// sqlite is lazy and it only checks on first access; force it
	_, err = s.Exec("BEGIN;COMMIT")
	if (err != nil ) { return errors.New("init: " + err.Error() + "; dbfile: " + path) }
	if out.nosync {
		_, err = s.Exec("PRAGMA  synchronous = 0")
		if (err != nil ) { return errors.New("desync: " + err.Error()) }
	}
	_, err = s.Exec("CREATE TABLE IF NOT EXISTS asd ( time INT, tags TEXT, a TEXT )",`asd`)
	if (err != nil ) { return errors.New("CT: " + err.Error()) }
	rows, err := s.Query("PRAGMA table_info( ? )",`asd`)
	if (err != nil ) { return errors.New("Pragma: " + err.Error()) }
	s.Close()
	return errors.New(fmt.Sprintf("%+v %s %s %s",rows, err,  prefix, path))

}

func quoteTableName (in string)(out string) {
	//fixme
	return in
}
