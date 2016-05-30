package station

import (
	"time"

	"panahon/logger"
)

func TestRoutine() {
	for {
		logger.Info.Println("--------- CONCURRENCY TEST ALERT ---------")
		time.Sleep(5000 * time.Millisecond)
	}
}
