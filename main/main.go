package main

import (
	"os"

	"Panahon/database"
	"Panahon/logger"
	"Panahon/server"
	"Panahon/station"
	"github.com/BurntSushi/toml"
	"github.com/influxdata/influxdb/client/v2"
)

type Config struct {
	AppPort string      `toml:"app_port"`
	DB      Database    `toml:"database"`
	App     Application `toml:"app"`
}

type Database struct {
	Server string
	Port   string
	RootDB string `toml:"root_db"`
	Series string
}

type Application struct {
	Path string
}

func parseConfig(configPath string, config *Config) error {
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		return err
	}
	return nil
}

// main is the Main Function of the Program
func main() {
	logger.Init(os.Stdout, os.Stdout, os.Stderr)
	logger.Info.Println("Logger initialized")

	config := new(Config)
	err := parseConfig("./config.toml", config)
	if err != nil {
		logger.Error.Fatalf("Error while parsing Config File: %s", err)
	}

	logger.Info.Println("Config file parsed, config set")

	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://" + config.DB.Server + ":" + config.DB.Port,
	})
	if err != nil {
		logger.Error.Fatalln(err)
	}

	database.Init(influxClient, config.DB.RootDB, config.DB.Series)
	logger.Info.Printf(
		"InfluxDB client initialized on %s:%s, DB: %s, Series: %s",
		config.DB.Server,
		config.DB.Port,
		config.DB.RootDB,
		config.DB.Series,
	)

	go station.TestRoutine()
	server.StartServer(config.AppPort, config.App.Path)
}
