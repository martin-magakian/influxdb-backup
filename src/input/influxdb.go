package input

import (
	"github.com/influxdb/influxdb/client/v2"
	"time"
	"github.com/op/go-logging"
	"fmt"
	"common"
	"encoding/json"
)
var log = logging.MustGetLogger("main")


func NewInflux09(addr string, user string, pass string, db string) (Input, error) {
	var influx Influx09Input
	var err error
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     addr,
		Username: user,
		Password: pass,
	})
	if err != nil {
		return &influx,err
	}
	influx.Db = db
	influx.Client = c

	return &influx,  err
}


type Influx09Input struct {
	Addr string
	Client client.Client
	Db string
}

func (influx *Influx09Input) GetSeriesList() ([]string, error) {
	var err error
	var out []string
	q := client.NewQuery(`show series`, "stats", "ns")
	res  ,err  := influx.Client.Query(q)
	if res.Error() != nil {
		err = res.Error()
		return out, err
	}

	for _, v := range res.Results[0].Series {
		out = append(out,v.Name)

	}

	return out, err
}

func (influx *Influx09Input) GetFieldRangeByName(name string, start time.Time, end time.Time) ([]common.Field, error) {
	var datapoints []common.Field
	var err error
	q := client.NewQuery(
		fmt.Sprintf(`SELECT * FROM "%s" WHERE time > %d AND time < %d`,
			name,
			start.UnixNano(),
			end.UnixNano()),
		influx.Db,
		`ns`,
	)
	res, err := influx.Client.Query(q)
	if res.Error() != nil {
		err = res.Error()
		return datapoints, err
	}
	for _, result := range res.Results {
		for _, ser := range result.Series {
			columns := ser.Columns
			for _, value := range ser.Values {
				var f common.Field
				f.Name = ser.Name
				f.Tags = ser.Tags
				f.Values = make(map[string]interface{})
				for j, value := range value {
					// influxdb lib is a bit funny
					if w, ok := value.(json.Number); ok {
						// check if it is an int or float
						a, _ := w.Float64()

						if a == float64(int64(a)) {
							f.Values[ columns[j] ], _ = w.Int64()
						} else {
							f.Values[ columns[j] ] = a
						}
					} else {
						f.Values[ columns[j] ] = value
					}
				}
				datapoints = append(datapoints,f)
			}
		 }
	}
	return datapoints, err
}
