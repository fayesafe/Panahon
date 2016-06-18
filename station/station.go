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

type PinLayout struct {
	DHTPin  int
	LDRPin  int
	RainPin int
}

type dbClient interface {
	QueryAll(offset string) (*client.Response, error)
	QueryInterval(low string, high string) (*client.Response, error)
	QueryAverage(col string, interval string, offset string, end string) (*client.Response, error)
	QueryMax(col string, interval string, offset string, end string) (*client.Response, error)
}

func TestRoutine() {
	for {
		logger.Info.Println("-- CONCURRENCY TEST ALERT --")
		time.Sleep(5000 * time.Millisecond)
	}
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
	temp, humidity := readDHT22(DHTPin)
	temp2, pressure := readBMP180(1)
	rain := readRain(RainPin)
	sun := readLDR(LDRPin)

	// Create a point and add to batch
	tags := map[string]string{}
	fields := map[string]interface{}{
		"temperature": (temp + float32(temp2)) / 2,
		"pressure":    pressure,
		"Humidity":    humidity,
		"Rain":        rain,
		"Sun":         sun,
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

func readLDR(port int) int {
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
	return count
}

func readDHT22(port int) (float32, float32) {
	temp, hum, retried, err := dht.ReadDHTxxWithRetry(dht.DHT22, port, true, 10)
	if err != nil {
		log.Fatal(err)
	}
	logger.Info.Printf("Temperature= %v°C, Humidity= %v%% (retried %d times), (DHT22)\n", temp, hum, retried)
	return temp, hum
}

func readRain(port int) bool {
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
	return !dry
}

func readBMP180(port byte) (float64, int) {
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
	return temp, pressure
}
