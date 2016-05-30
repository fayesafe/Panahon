package main

import (
	"os"

	"Panahon/logger"
	"Panahon/server"
	"Panahon/station"
)

func main() {
	logger.Init(os.Stdout, os.Stdout, os.Stderr)
	logger.Info.Println("Logger initialized")
	go station.TestRoutine()
	server.StartServer()
}
