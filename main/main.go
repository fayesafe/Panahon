package main

import (
	"os"

	"panahon/logger"
	"panahon/server"
	"panahon/station"
)

func main() {
	logger.Init(os.Stdout, os.Stdout, os.Stderr)
	logger.Info.Println("Logger initialized")
	go station.TestRoutine()
	server.StartServer()
}
