package input

import (
//	"common"
	"time"
)


type Input interface {
	GetSeriesList() []string
	GetFieldRangeByName(name string, start time.Time, end time.Time)
}


func New(inputName string, inputAddr string) {

}
