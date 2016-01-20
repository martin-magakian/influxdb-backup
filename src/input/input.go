package input

import (
	"common"
	"time"
//	"github.com/influxdb/influxdb/client/v2"
)


type Input interface {
	GetSeriesList() ([]string, error)
	GetFieldRangeByName(name string, start time.Time, end time.Time) ([]common.Field,error)
//	GetFieldRangeByName(name string, start time.Time, end time.Time) (string,error)
}


func New(inputName string, inputAddr string) {

}
