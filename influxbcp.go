package main

import (
	//	"fmt"
	"github.com/codegangsta/cli"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/op/go-logging"
	//	"gopkg.in/yaml.v2"
	"github.com/efigence/influxdb-backup/common"
	"github.com/efigence/influxdb-backup/input"
	"github.com/efigence/influxdb-backup/output"
	"os"
	"sync"
	"time"
	//	"strings"
)

var version string
var log = logging.MustGetLogger("main")
var stdout_log_format = logging.MustStringFormatter("%{color:bold}%{time:2006-01-02T15:04:05.9999Z-07:00}%{color:reset}%{color} [%{level:.1s}] %{color:reset}%{shortpkg}[%{longfunc}] %{message}")

const (
	MyDB     = "stats"
	username = "root"
	password = "root"
)

type Config struct {
	SourceType      string
	SourceAddr      string
	DebugAddr       string
	DestinationType string
	DestinationAddr string
}

func main() {
	var cfg Config
	stderrBackend := logging.NewLogBackend(os.Stderr, "", 0)
	stderrFormatter := logging.NewBackendFormatter(stderrBackend, stdout_log_format)
	logging.SetBackend(stderrFormatter)
	logging.SetFormatter(stdout_log_format)

	log.Info("Starting app")
	log.Debug("version: %s", version)
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "src-type",
			Value:       "influx09",
			Usage:       "Source type [influx09]",
			Destination: &cfg.SourceType,
		},
		cli.StringFlag{
			Name:        "src-addr",
			Value:       "http://localhost:8086",
			Usage:       "Source addr",
			Destination: &cfg.SourceAddr,
		},
		cli.StringFlag{
			Name:        "dst-type",
			Value:       "influx09",
			Usage:       "Destination type [influx09]",
			Destination: &cfg.SourceType,
		},
		cli.StringFlag{
			Name:        "dst-addr",
			Value:       "http://localhost:8086",
			Usage:       "Destination addr",
			Destination: &cfg.SourceAddr,
		},
		cli.StringFlag{
			Name:        "debug-addr",
			Value:       "none",
			Usage:       "Listen address of pprof in form of http://localhost:12345",
			Destination: &cfg.DebugAddr,
		},
	}
	app.Run(os.Args)
	c, _ := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: username,
		Password: password,
	})
	_ = c
	if cfg.DebugAddr != "none" {
		common.RunDebug(cfg.DebugAddr)
	} else {
		log.Error("wtf")
		common.RunDebug("localhost:6060")
	}
	log.Notice("Source type: %s, addr: %s", cfg.SourceType, cfg.SourceAddr)
	influxIn, err := input.NewInflux09("http://localhost:8086", "root", "root", "_internal")
	if err != nil {
		log.Error("input failed: %s", err)
		os.Exit(1)
	}
	sqliteOut, err := output.New(`sqlite`, `t-data/sqlite`)
	if err != nil {
		log.Error("output failed: %s", err)
		os.Exit(1)
	}
	series, err := influxIn.GetSeriesList()
	var wg sync.WaitGroup
	for _, ser := range series {
		fields, err := influxIn.GetFieldRangeByName(ser, time.Now().Add(-1*time.Hour), time.Now())
		ch := make(chan *common.Field, 1)
		err = sqliteOut.Run([]chan *common.Field{ch})
		if err != nil {
			log.Error("running writer failed: %s", err)
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, ev := range fields {
				ch <- &ev
			}
			close(ch)

		}()
	}
	wg.Wait()
	err = sqliteOut.SaveSeriesList(series)
	if err != nil {
		log.Error("saving series list failed: %s", err)
		os.Exit(1)
	}
	sqliteOut.Shutdown()
	log.Notice("Written %d records total", sqliteOut.GetTotalWrites())

	//	log.Info("v: %+v", series)

}
