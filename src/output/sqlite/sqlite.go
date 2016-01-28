package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"path/filepath"
	"os"
	"hash/fnv"
	"errors"
	"fmt"
//	"strings"
	"github.com/efigence/influxdb-backup/src/common"

)

type SQLiteOut struct {
	spineMask uint64  // mask, MSB bits
	leafMask uint64   // mask, LSB bits
	path string
	leafBits uint8 // number of directories
	spineBits uint8 // number of files per dir
	nosync bool
	workers *writers
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
	out.leafBits = 2
	out.spineBits = 9
	// spine uses MSB bits
	out.spineMask = ^uint64(0) << (64-out.spineBits)
	// leaf uses LSB bits
	out.leafMask = powOf2(out.leafBits) - 1
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
		_, err = s.Exec("INSERT OR IGNORE INTO series(name, file) VALUES(?,?)",name,out.SeriesNameGen(name))
		if (err != nil ) { return err }
	}
	_, err = s.Exec("COMMIT")
	if (err != nil ) { return err }
	err = s.Close()

	return err
}


// Generate shortened series name
// does not have to be unique, just unique enough that tens of thousands of series wont land in same sqlite DB
func (out *SQLiteOut) SeriesNameGen(seriesName string) string {
	hasher := fnv.New64a()
	hasher.Write([]byte(seriesName))
	hash := hasher.Sum64()

	leaf := hash & out.leafMask
	spine := (hash & out.spineMask)  >> (64 - out.spineBits)

	return fmt.Sprintf("%x/%x.sqlite",spine, leaf)
}

func (out *SQLiteOut) SaveFields(prefix string) (error) {
	s, err := sqliteOpen( []string{out.path, prefix + ".sqlite"} , false )
	_, err = s.Exec("CREATE TABLE IF NOT EXISTS asd ( time INT, tags TEXT, a TEXT )",`asd`)
	if (err != nil ) { return errors.New("CT: " + err.Error()) }
	rows, err := s.Query("PRAGMA table_info( ? )",`asd`)
	if (err != nil ) { return errors.New("Pragma: " + err.Error()) }
	s.Close()
	return errors.New(fmt.Sprintf("%+v %s %s %s",rows, err,  prefix))

}

func quoteTableName (in string)(out string) {
	//fixme
	return in
}

func sqliteOpen(pathComponents []string, nosync bool) (s *sql.DB, err error) {
	path :=  filepath.Join(pathComponents...)
	dir, _ := filepath.Split(path)
	err = os.MkdirAll(dir,0755)
	if (err != nil ) { return s, errors.New("mkdir: " + err.Error() + "; path: " + dir ) }
	s, err = sql.Open("sqlite3",path)
	if (err != nil ) { return s, errors.New("sql open: " + err.Error() + "; dbfile: " + path ) }
	// sqlite is lazy and it only checks on first access; force it
	_, err = s.Exec("BEGIN;COMMIT")
	if (err != nil ) { return s, errors.New("sql init: " + err.Error() + "; dbfile: " + path ) }
	if nosync {
		_, err = s.Exec("PRAGMA  synchronous = 0")
		if (err != nil ) { return s, errors.New("sql desync: " + err.Error()) }
	}

	return s,err
}


// there is no generic pow() in golang stdlib()
func powOf2(b uint8) uint64 {
  var result uint64 = 1;
  var a uint64 = 2
  for 0 != b {
    if 0 != (b & 1) {
      result *= a;

    }
    b >>= 1;
    a *= a;
  }

  return result;
}
