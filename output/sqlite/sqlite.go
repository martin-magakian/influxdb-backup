package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"hash/fnv"
	"os"
	"path/filepath"
	"sync"
	//	"strings"
	"github.com/efigence/influxdb-backup/common"
	"github.com/op/go-logging"
)

type SQLiteOut struct {
	spineMask uint64 // mask, MSB bits
	leafMask  uint64 // mask, LSB bits
	path      string
	leafBits  uint8 // number of directories
	spineBits uint8 // number of files per dir
	nosync    bool
	writers   *dbWriters
	routers   sync.WaitGroup
}

var log = logging.MustGetLogger("main")

func New(args []string) (common.Output, error) {
	return NewSQLite(args[0])
}

func NewSQLite(path string) (common.Output, error) {
	var err error
	var out SQLiteOut
	out.Init(path)
	return &out, err
}

func (out *SQLiteOut) Init(path string) {
	var mode os.FileMode
	mode = 0744
	out.path = path
	out.leafBits = 2
	out.spineBits = 4
	log.Info("Initializing SQLite store in %s using %d dirs and %d files per dir", path, powOf2(out.spineBits), powOf2(out.leafBits))
	out.spineMask = ^uint64(0) << (64 - out.spineBits)
	out.leafMask = powOf2(out.leafBits) - 1
	os.MkdirAll(path, mode)
	worker := out.newWriter()
	out.writers = &worker

}

// start writing data
func (out *SQLiteOut) Run(in []chan *common.Field) (err error) {
	for i, ch := range in {
		out.routers.Add(1)
		log.Debug("Running router gor %d", i)
		go out.route(ch)
	}
	return err

}

// stop writing and close data
func (out *SQLiteOut) Shutdown() (err error) {
	out.routers.Wait()
	out.writers.Shutdown()
	return err
}

func (out *SQLiteOut) GetTotalWrites() uint64 {
	return out.writers.writes
}

func (out *SQLiteOut) SaveSeriesList(series []string) (err error) {
	db, err := sql.Open("sqlite3", filepath.Join(out.path, "series.sqlite"))
	if err != nil {
		return err
	}
	if out.nosync {
		_, err = db.Exec("PRAGMA  synchronous = 0")
		if err != nil {
			return err
		}
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS series( name TEXT UNIQUE , file TEXT )")
	if err != nil {
		return err
	}
	_, err = db.Exec("BEGIN")
	if err != nil {
		return err
	}
	for _, name := range series {
		_, err = db.Exec("INSERT OR IGNORE INTO series(name, file) VALUES(?,?)", name, out.SeriesNameGen(name))
		if err != nil {
			return err
		}
	}
	_, err = db.Exec("COMMIT")
	if err != nil {
		return err
	}
	err = db.Close()

	return err
}

// Generate shortened series name
// does not have to be unique, just unique enough that tens of thousands of series wont land in same sqlite DB
func (out *SQLiteOut) SeriesNameGen(seriesName string) string {
	hasher := fnv.New64a()
	hasher.Write([]byte(seriesName))
	hash := hasher.Sum64()

	leaf := hash & out.leafMask
	spine := (hash & out.spineMask) >> (64 - out.spineBits)

	return fmt.Sprintf("%x/%x.sqlite", spine, leaf)
}

func (out *SQLiteOut) SaveFields(prefix string) error {
	db, err := sqliteOpen([]string{out.path, prefix + ".sqlite"}, false)
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS asd ( time INT, tags TEXT, a TEXT )", `asd`)
	if err != nil {
		return errors.New("CT: " + err.Error())
	}
	rows, err := db.Query("PRAGMA table_info( ? )", `asd`)
	if err != nil {
		return errors.New("Pragma: " + err.Error())
	}
	db.Close()
	return errors.New(fmt.Sprintf("%+v %s %s %s", rows, err, prefix))

}

func quoteTableName(in string) (out string) {
	//fixme
	return in
}

func sqliteOpen(pathComponents []string, nosync bool) (db *sql.DB, err error) {
	path := filepath.Join(pathComponents...)
	dir, _ := filepath.Split(path)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return db, errors.New("mkdir: " + err.Error() + "; path: " + dir)
	}
	db, err = sql.Open("sqlite3", path)
	if err != nil {
		return db, errors.New("sql open: " + err.Error() + "; dbfile: " + path)
	}
	// sqlite is lazy and it only checks on first access; force it
	_, err = db.Exec("BEGIN;COMMIT")
	if err != nil {
		return db, errors.New("sql init: " + err.Error() + "; dbfile: " + path)
	}
	if nosync {
		_, err = db.Exec("PRAGMA  synchronous = OFF")
		if err != nil {
			return db, errors.New("sql desync: " + err.Error())
		}
	}

	return db, err
}

// there is no generic pow() in golang stdlib()
func powOf2(b uint8) uint64 {
	var result uint64 = 1
	var a uint64 = 2
	for 0 != b {
		if 0 != (b & 1) {
			result *= a

		}
		b >>= 1
		a *= a
	}

	return result
}
