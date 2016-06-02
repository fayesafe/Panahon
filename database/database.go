package database

import (
	"Panahon/logger"
	"github.com/influxdata/influxdb/client/v2"
)

type dbClient interface {
	Query(q client.Query) (*client.Response, error)
}

// Global variable of the database client
var influxClient dbClient

// QueryAll queries all entries of Influx DB, or last n entries
func QueryAll(offset string) (*client.Response, error) {
	var q client.Query

	if offset != "" {
		q = client.NewQuery(
			"SELECT * FROM meas ORDER BY DESC LIMIT "+offset,
			"test",
			"s")
		logger.Info.Printf("Getting last %s entries", offset)
	} else {
		q = client.NewQuery("SELECT * FROM meas", "test", "s")
		logger.Info.Println("Calling route /api/get")
	}

	response, err := influxClient.Query(q)
	return response, err
}

// QueryInterval queries Influx DB for an interval of time
func QueryInterval(low string, high string) (*client.Response, error) {
	q := client.NewQuery(
		"SELECT * FROM meas WHERE time < "+
			high+"s and time >"+low+"s", "test", "s")
	logger.Info.Printf(
		"Getting entries from timestamp %s to %s", low, high)
	response, err := influxClient.Query(q)
	return response, err
}

// Init of database client
func Init(influxConn dbClient) {
	influxClient = influxConn
}
