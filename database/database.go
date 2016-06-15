package database

import (
	"fmt"

	"Panahon/logger"

	"github.com/influxdata/influxdb/client/v2"
)

type dbClient struct {
	Client   client.Client
	Database string
	Series   string
}

// QueryAll queries all entries of Influx DB, or last n entries
func (influx dbClient) QueryAll(offset string) (*client.Response, error) {
	var q client.Query

	if offset != "" {
		query := fmt.Sprintf(
			"SELECT * FROM %s ORDER BY DESC LIMIT %s",
			influx.Series,
			offset)
		q = client.NewQuery(query, influx.Database, "ms")
		logger.Info.Printf("Getting last %s entries", offset)
	} else {
		query := fmt.Sprintf("SELECT * FROM %s", influx.Series)
		q = client.NewQuery(query, influx.Database, "ms")
		logger.Info.Println("Calling route /api/get")
	}

	return influx.Client.Query(q)
}

// QueryInterval queries Influx DB for an interval of time
func (influx dbClient) QueryInterval(low string, high string) (*client.Response, error) {
	query := fmt.Sprintf(
		"SELECT * FROM %s WHERE time < %sms AND time >= %sms",
		influx.Series,
		high,
		low)

	q := client.NewQuery(query, influx.Database, "ms")
	logger.Info.Printf(
		"Getting entries from timestamp %s to %s", low, high)
	return influx.Client.Query(q)
}

// QueryAverage queries the average for given cols on a given interval
func (influx dbClient) QueryAverage(
	col string,
	interval string,
	offset string,
	end string) (*client.Response, error) {
	query := fmt.Sprintf(
		"SELECT mean(%s) FROM %s WHERE time >= %sms AND time < %sms GROUP BY time(%s)",
		col,
		influx.Series,
		offset,
		end,
		interval)

	q := client.NewQuery(query, influx.Database, "ms")
	logger.Info.Printf("Getting average from %s of col %s", offset, col)
	return influx.Client.Query(q)
}

func (influx dbClient) QueryMax(
	col string,
	interval string,
	offset string,
	end string) (*client.Response, error) {
	query := fmt.Sprintf(
		"SELECT max(%s) FROM %s WHERE time >= %sms AND time < %sms GROUP BY time(%s)",
		col,
		influx.Series,
		offset,
		end,
		interval)

	q := client.NewQuery(query, influx.Database, "ms")
	logger.Info.Printf("Getting max value from %s to %s of col %s",
		offset, end, col)
	return influx.Client.Query(q)
}

// Init of database client
func Init(database string, series string, server string, port string) *dbClient {
	databaseConn := new(dbClient)
	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://" + server + ":" + port,
	})
	if err != nil {
		logger.Error.Fatalln(err)
	}
	databaseConn.Client = influxClient
	databaseConn.Database = database
	databaseConn.Series = series
	return databaseConn
}
