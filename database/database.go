package database

import (
	"Panahon/logger"
	"github.com/influxdata/influxdb/client/v2"
)

type dbClient interface {
	Query(q client.Query) (*client.Response, error)
}

var influxClient dbClient

func QueryAll(offset string) (*client.Response, error) {
	var q client.Query

	if offset != "" {
		q = client.NewQuery(
			"SELECT * FROM meas ORDER BY DESC LIMIT "+offset,
			"test",
			"s")
		logger.Info.Println("Getting last " + offset + " entries")
	} else {
		q = client.NewQuery("SELECT * FROM meas", "test", "s")
		logger.Info.Println("Calling route /api/get")
	}

	response, err := influxClient.Query(q)
	return response, err
}

func QueryInterval(low string, high string) (*client.Response, error) {
	logger.Info.Println(low, high)
	q := client.NewQuery(
		"SELECT * FROM meas WHERE time < "+
			high+"s and time >"+low+"s", "test", "s")

	response, err := influxClient.Query(q)
	return response, err
}

func Init(influxConn dbClient) {
	influxClient = influxConn
}
