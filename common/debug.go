package common

import (
	_ "net/http/pprof"
	"net/http"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("main")

func RunDebug (addr string) {
	go func() {
		log.Warning("starting debug server on %s $+v", addr, http.ListenAndServe(addr, nil))
	}()
}
