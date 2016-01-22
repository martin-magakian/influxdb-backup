package output

import (
	"github.com/efigence/influxdb-backup/src/output/sqlite"
	"github.com/efigence/influxdb-backup/src/common"
)


var outputs =  map[string]func(args []string) (common.Output, error)  {
	"sqlite": sqlite.New,
}


func New(outputType string, args ...string) (common.Output, error) {
	return outputs[outputType](args)

}
