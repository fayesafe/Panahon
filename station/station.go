package station

import (
	"log"
	"time"

	"github.com/d2r2/go-dht"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/kidoman/embd"
	"github.com/kidoman/embd/sensor/bmp180"
	"github.com/kidoman/embd/sensor/watersensor"

	_ "github.com/kidoman/embd/host/rpi"

	"Panahon/database"
	"Panahon/logger"
)

type dht22Result struct {
	Temperature float32
	Humidity    float32
}

type bmp180Result struct {
	Temperature float64
	Pressure    int
}

type rainResult struct {
	Rain int
}

type ldrResult struct {
	Sun int
}

type dbClient interface {
	QueryAll(offset string) (*client.Response, error)
	QueryInterval(low string, high string) (*client.Response, error)
	QueryAverage(col string, interval string, offset string, end string) (*client.Response, error)
	QueryMax(col string, interval string, offset string, end string) (*client.Response, error)
}

func StartReadRoutine(influx *database.DBClient, DHTPin int, LDRPin int, RainPin int) {
	for {
		logger.Info.Println("Measurements started...")
		go ReadSensors(influx, DHTPin, LDRPin, RainPin)
		logger.Info.Printf("Going to sleep for %d minutes", 10)
		time.Sleep(10 * time.Minute)
	}
}

func ReadSensors(influx *database.DBClient, DHTPin int, LDRPin int, RainPin int) {
	sensorResults := make(chan interface{}, 4)
	fields := map[string]interface{}{}
	tags := map[string]string{}

	go readDHT22(sensorResults, DHTPin)
	go readBMP180(sensorResults, 1)
	go readRain(sensorResults, RainPin)
	go readLDR(sensorResults, LDRPin)

	// Create a point and add to batch
	for result := range sensorResults {
		switch result.(type) {
		case *dht22Result:
			tmp, _ := result.(dht22Result)
			if val, ok := fields["temperature"]; ok {
				fields["temperature"] = (tmp.Temperature + val.(float32)) / 2
			} else {
				fields["temperature"] = tmp.Temperature
			}
			fields["humidity"] = tmp.Humidity
		case *bmp180Result:
			tmp, _ := result.(bmp180Result)
			if val, ok := fields["temperature"]; ok {
				fields["temperature"] = (float32(tmp.Temperature) + val.(float32)) / 2
			} else {
				fields["temperature"] = float32(tmp.Temperature)
			}
			fields["pressure"] = tmp.Pressure
		case *rainResult:
			tmp, _ := result.(rainResult)
			fields["rain"] = tmp.Rain
		case *ldrResult:
			tmp, _ := result.(ldrResult)
			fields["sun"] = tmp.Sun
		}
	}

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  influx.Database,
		Precision: "ms",
	})
	pt, err := client.NewPoint(influx.Series, tags, fields, time.Now())
	if err != nil {
		logger.Error.Println(err)
	}

	// write point to db
	bp.AddPoint(pt)
}

func readDHT22(sensorResults chan interface{}, port int) {
	temp, hum, retried, err := dht.ReadDHTxxWithRetry(dht.DHT22, port, true, 10)
	if err != nil {
		log.Fatal(err)
	}
	logger.Info.Printf("Temperature= %v°C, Humidity= %v%% (retried %d times), (DHT22)\n", temp, hum, retried)

	sensorResults <- dht22Result{temp, hum}
}

func readLDR(sensorResults chan interface{}, port int) {
	if err := embd.InitGPIO(); err != nil {
		logger.Error.Panicln(err)
	}
	defer embd.CloseGPIO()

	ldr, err := embd.NewDigitalPin(port)
	if err != nil {
		logger.Error.Panicln(err)
	}
	defer ldr.Close()

	if err := ldr.SetDirection(embd.Out); err != nil {
		logger.Error.Panicln(err)
	}
	if err := ldr.Write(embd.Low); err != nil {
		logger.Error.Panicln(err)
	}

	time.Sleep(100 * time.Millisecond)

	if err := ldr.SetDirection(embd.In); err != nil {
		logger.Error.Panicln(err)
	}

	count := 0
	j, err := ldr.Read()
	if err != nil {
		logger.Error.Panicln(err)
	}
	for j == embd.Low {
		count++
		j, err = ldr.Read()
		if err != nil {
			logger.Error.Panicln(err)
		}
	}
	logger.Info.Printf("Sun Value=%d", count)
	sensorResults <- ldrResult{count}
}

func readRain(sensorResults chan interface{}, port int) {
	if err := embd.InitGPIO(); err != nil {
		logger.Error.Panicln(err)
	}
	defer embd.CloseGPIO()

	pin, err := embd.NewDigitalPin(port)
	if err != nil {
		logger.Error.Panicln(err)
	}
	defer pin.Close()

	fluidSensor := watersensor.New(pin)

	dry, err := fluidSensor.IsWet()
	if err != nil {
		logger.Error.Panicln(err)
	}

	logger.Info.Printf("Rain=%t\n", !dry)
	if dry {
		sensorResults <- rainResult{0}
	} else {
		sensorResults <- rainResult{1}
	}
}

func readBMP180(sensorResults chan interface{}, port byte) {
	if err := embd.InitI2C(); err != nil {
		logger.Error.Panicln(err)
	}
	defer embd.CloseI2C()

	bus := embd.NewI2CBus(port)

	baro := bmp180.New(bus)
	defer baro.Close()

	temp, err := baro.Temperature()
	if err != nil {
		panic(err)
	}
	logger.Info.Printf("Temperature=%v (BMP180)\n", temp)

	pressure, err := baro.Pressure()
	if err != nil {
		logger.Error.Panicln(err)
	}

	logger.Info.Printf("Pressure=%v hPa\n", pressure)
	sensorResults <- bmp180Result{temp, pressure}
}
