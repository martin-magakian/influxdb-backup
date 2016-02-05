package common

import (
	"github.com/op/go-logging"
	"net/http"
	_ "net/http/pprof"
)

var log = logging.MustGetLogger("main")

func RunDebug(addr string) {
	go func() {
		log.Warning("starting debug server on %s $+v", addr, http.ListenAndServe(addr, nil))
	}()
}
