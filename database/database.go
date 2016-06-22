package database

import (
	"fmt"
	"time"

	"Panahon/logger"

	"github.com/influxdata/influxdb/client/v2"
)

type DBClient struct {
	Client   client.Client
	Database string
	Series   string
}

func (influx DBClient) SaveData(fields map[string]interface{}, tags map[string]string) {
	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  influx.Database,
		Precision: "ms",
	})
	pt, err := client.NewPoint(influx.Series, tags, fields, time.Now())
	if err != nil {
		logger.Error.Println(err)
		return
	}

	// write point to db
	bp.AddPoint(pt)
	influx.Client.Write(bp)
	logger.Info.Printf("Data written to db: %v", fields)
}

func (influx DBClient) QueryData(series string, startDate string, endDate string) (*client.Response, error) {
	var q client.Query

	query := fmt.Sprintf(
		"SELECT * FROM %s WHERE %s <= time AND time < %s", series, startDate, endDate)
	q = client.NewQuery(query, influx.Database, "ms")
	logger.Info.Printf("Getting data of %s from %s to %s", series, startDate, endDate)
	logger.Info.Println(query)

	return influx.Client.Query(q)
}

func (influx DBClient) QuerySensorData(last string) (*client.Response, error) {
	var q client.Query

	query := fmt.Sprintf(
		"SELECT * FROM sensors ORDER BY DESC LIMIT %s", last)
	q = client.NewQuery(query, influx.Database, "ms")
	logger.Info.Printf("Getting last %s entries", last)

	return influx.Client.Query(q)
}

// QueryAll queries all entries of Influx DB, or last n entries
func (influx DBClient) QueryAll(offset string) (*client.Response, error) {
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
func (influx DBClient) QueryInterval(low string, high string) (*client.Response, error) {
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
func (influx DBClient) QueryAverage(
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

func (influx DBClient) QueryMax(
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

func (influx DBClient) QueryMin(
	col string,
	interval string,
	offset string,
	end string) (*client.Response, error) {
	query := fmt.Sprintf(
		"SELECT min(%s) FROM %s WHERE time >= %sms AND time < %sms GROUP BY time(%s)",
		col,
		influx.Series,
		offset,
		end,
		interval)

	q := client.NewQuery(query, influx.Database, "ms")
	logger.Info.Printf("Getting min value from %s to %s of col %s",
		offset, end, col)
	return influx.Client.Query(q)
}

// Init of database client
func Init(database string, series string, server string, port string) *DBClient {
	databaseConn := new(DBClient)
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
