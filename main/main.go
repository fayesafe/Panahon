package main

import (
	"os"

	"Panahon/database"
	"Panahon/logger"
	"Panahon/server"
	"Panahon/station"
	"github.com/influxdata/influxdb/client/v2"
)

// main is the Main Function of the Program
func main() {
	logger.Init(os.Stdout, os.Stdout, os.Stderr)
	logger.Info.Println("Logger initialized")

	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if err != nil {
		logger.Error.Println(err)
		os.Exit(1)
	}

	database.Init(influxClient)
	logger.Info.Println("InfluxDB client initialized")

	go station.TestRoutine()
	server.StartServer()
}
