package main

import (
	"fmt"
	"github.com/influxdb/influxdb/client/v2"
	"github.com/op/go-logging"
	"os"
	"strings"
)

var version string
var log = logging.MustGetLogger("main")
var stdout_log_format = logging.MustStringFormatter("%{color:bold}%{time:2006-01-02T15:04:05.9999Z-07:00}%{color:reset}%{color} [%{level:.1s}] %{color:reset}%{shortpkg}[%{longfunc}] %{message}")

const (
	MyDB     = "stats"
	username = "root"
	password = "root"
)

func main() {
	stderrBackend := logging.NewLogBackend(os.Stderr, "", 0)
	stderrFormatter := logging.NewBackendFormatter(stderrBackend, stdout_log_format)
	logging.SetBackend(stderrFormatter)
	logging.SetFormatter(stdout_log_format)

	log.Info("Starting app")
	log.Debug("version: %s", version)
	if !strings.ContainsRune(version, '-') {
		log.Warning("once you tag your commit with name your version number will be prettier")
	}
	c, _ := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: username,
		Password: password,
	})
	//q := client.NewQuery(`SELECT * FROM "dc1.ghroth_non_3dart_com.cpu.0.cpu.system" WHERE time > now() - 1h `, "stats", "ns")
	q := client.NewQuery(`show series`, "stats", "ns")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		fmt.Printf("V: %+v\n", response.Results)
	}
	log.Info("v: %+v", c)

}
