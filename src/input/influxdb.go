package input

import (
	"github.com/influxdb/influxdb/client/v2"
	"github.com/op/go-logging"
//	"fmt"
)
var log = logging.MustGetLogger("main")


func NewInflux09(addr string, user string, pass string) (Influx09Input, error) {
	var influx Influx09Input
	var err error
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     addr,
		Username: user,
		Password: pass,
	})
	if err != nil {
		return influx,err
	}
	influx.Client = c

	return influx,  err
}


type Influx09Input struct {
	Addr string
	Client client.Client
}

func (influx Influx09Input) GetSeriesList() ([]string, error) {
	var err error
	var out []string
	query := client.NewQuery(`show series`, "stats", "ns")
	res  ,err  := influx.Client.Query(query)
	if res.Error() != nil {
		err = res.Error()
		return out, err
	}

	for _, v := range res.Results[0].Series {
		out = append(out,v.Name)

	}

	return out, err
}
