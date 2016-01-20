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


type DummyIn struct{}

func NewDummy(addr string, user string, pass string, db string) (Input, error) {
	var dummy DummyIn
	var err error
	return &dummy,err
}

func (d *DummyIn) GetSeriesList() ([]string, error) {
	var err error
	return []string{
		"dc1.host1.cpu.user",
		"dc2.host2.cpu.idle",
	}, err
}

func (d *DummyIn) GetFieldRangeByName(name string, start time.Time, end time.Time) ([]common.Field, error) {
	var err error
	var startField common.Field
	var midField common.Field
	var endField common.Field

	startTs := start.UnixNano()
	endTs := end.UnixNano()

	if startTs - endTs < 10 || endTs - startTs < 10 {

	}

	startField.Name = name
	startField.Tags = make(map[string]string)
	startField.Values = make(map[string]interface{})
	startField.Values[`time`] = startTs

	// same TS, emit only one record
	if startTs == endTs {
		return []common.Field{startField}, err
	}
	endField.Name = name
	endField.Tags = make(map[string]string)
	endField.Values = make(map[string]interface{})
	endField.Values[`time`] = endTs

	// very small spread, emit only start and end
	if startTs - endTs < time.Second.Nanoseconds() || endTs - startTs < time.Second.Nanoseconds() {
		return []common.Field{startField, endField}, err
	}

	midField.Name = name
	midField.Tags = make(map[string]string)
	midField.Values = make(map[string]interface{})
	midField.Values[`time`] = ( startTs + endTs ) / 2


	return []common.Field{startField, midField, endField}, err
}
