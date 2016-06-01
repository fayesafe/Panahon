package main

import (
	"os"

	"Panahon/logger"
	"Panahon/server"
	"Panahon/station"
)

// main is the Main Function of the Program
func main() {
	logger.Init(os.Stdout, os.Stdout, os.Stderr)
	logger.Info.Println("Logger initialized")
	go station.TestRoutine()
	server.StartServer()
}
