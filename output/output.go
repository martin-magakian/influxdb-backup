package output

import (
	"github.com/efigence/influxdb-backup/output/sqlite"
	"github.com/efigence/influxdb-backup/common"
)


var outputs =  map[string]func(args []string) (common.Output, error)  {
	"sqlite": sqlite.New,
}


func New(outputType string, args ...string) (common.Output, error) {
	return outputs[outputType](args)

}
