package main

import (
	//	"fmt"
	"github.com/codegangsta/cli"
	"github.com/influxdb/influxdb/client/v2"
	"github.com/op/go-logging"
	//	"gopkg.in/yaml.v2"
	"input"
	"os"
	"output"
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
	}
	app.Run(os.Args)
	c, _ := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: username,
		Password: password,
	})
	_ = c
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
	log.Info("saving series list %+v", series)
	err = sqliteOut.SaveSeriesList(series)
	if err != nil {
		log.Error("saving series list failed: %s", err)
		os.Exit(1)
	}

	//	log.Info("v: %+v", series)

}
