package database

import (
	"Panahon/logger"
	"github.com/influxdata/influxdb/client/v2"
)

type dbClient interface {
	Query(q client.Query) (*client.Response, error)
}

// Global variable of the database client
var (
	influxClient   dbClient
	influxDatabase string
	influxSeries   string
)

// QueryAll queries all entries of Influx DB, or last n entries
func QueryAll(offset string) (*client.Response, error) {
	var q client.Query

	if offset != "" {
		q = client.NewQuery(
			"SELECT * FROM "+influxSeries+" ORDER BY DESC LIMIT "+offset,
			influxDatabase,
			"s")
		logger.Info.Printf("Getting last %s entries", offset)
	} else {
		q = client.NewQuery(
			"SELECT * FROM "+influxSeries,
			influxDatabase,
			"s")
		logger.Info.Println("Calling route /api/get")
	}

	response, err := influxClient.Query(q)
	return response, err
}

// QueryInterval queries Influx DB for an interval of time
func QueryInterval(low string, high string) (*client.Response, error) {
	q := client.NewQuery(
		"SELECT * FROM "+influxSeries+" WHERE time < "+
			high+"s and time >"+low+"s", influxDatabase, "s")
	logger.Info.Printf(
		"Getting entries from timestamp %s to %s", low, high)
	response, err := influxClient.Query(q)
	return response, err
}

// Init of database client
func Init(influxConn dbClient, database string, series string) {
	influxClient = influxConn
	influxDatabase = database
	influxSeries = series
}
