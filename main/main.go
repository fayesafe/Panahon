package main

import (
	"os"

	"Panahon/database"
	"Panahon/logger"
	"Panahon/server"
	"Panahon/station"
	"github.com/BurntSushi/toml"
)

type Config struct {
	AppPort        string      `toml:"app_port"`
	DB             Database    `toml:"database"`
	App            Application `toml:"app"`
	WeatherStation Station     `toml:"station"`
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

type Station struct {
	Rain   int
	DHT22  int
	LDR    int
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

	influxClient := database.Init(
		config.DB.RootDB,
		config.DB.Series,
		config.DB.Server,
		config.DB.Port,
	)
	logger.Info.Printf(
		"InfluxDB client initialized on %s:%s, DB: %s, Series: %s",
		config.DB.Server,
		config.DB.Port,
		config.DB.RootDB,
		config.DB.Series,
	)

	go station.StartReadRoutine(influxClient, config.WeatherStation.DHT22, config.WeatherStation.LDR, config.WeatherStation.Rain)
	server.StartServer(influxClient, config.AppPort, config.App.Path)
}
