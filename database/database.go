package database

import (
	"fmt"

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
		query := fmt.Sprintf(
			"SELECT * FROM %s ORDER BY DESC LIMIT %s",
			influxSeries,
			offset)
		q = client.NewQuery(query, influxDatabase, "s")
		logger.Info.Printf("Getting last %s entries", offset)
	} else {
		query := fmt.Sprintf("SELECT * FROM %s", influxSeries)
		q = client.NewQuery(query, influxDatabase, "s")
		logger.Info.Println("Calling route /api/get")
	}

	return influxClient.Query(q)
}

// QueryInterval queries Influx DB for an interval of time
func QueryInterval(low string, high string) (*client.Response, error) {
	query := fmt.Sprintf(
		"SELECT * FROM %s WHERE time < %ss AND time > %ss",
		influxSeries,
		high,
		low)

	q := client.NewQuery(query, influxDatabase, "s")
	logger.Info.Printf(
		"Getting entries from timestamp %s to %s", low, high)
	return influxClient.Query(q)
}

// QueryAverage queries the average for given cols on a given interval
func QueryAverage(
	col string,
	interval string,
	offset string) (*client.Response, error) {
	query := fmt.Sprintf(
		"SELECT mean(%s) FROM %s WHERE time > %ss GROUP BY time(%s)",
		col,
		influxSeries,
		offset,
		interval)
	
	q := client.NewQuery(query, influxDatabase, "s")
	logger.Info.Println("Getting average from %s of col %s", offset, col)
	return influxClient.Query(q)
}

// Init of database client
func Init(influxConn dbClient, database string, series string) {
	influxClient = influxConn
	influxDatabase = database
	influxSeries = series
}
