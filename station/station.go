package station

import (
    "log"
    "time"

    "github.com/d2r2/go-dht"
    "github.com/kidoman/embd"
    "github.com/kidoman/embd/sensor/bmp180"
    "github.com/kidoman/embd/sensor/watersensor"
    "github.com/influxdata/influxdb/client/v2"

    "Panahon/logger"
    "Panahon/database"
)

type PinLayout struct {
    DHTPin  int
    LDRPin  int
    RainPin int
}

type StationResult struct {
    Temperature float32
    Humidity    float32
    Pressure    int
    IsWet       bool
    SunValue    bool
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
        temp, hum := readDHT22(DHTPin)
        var sunny bool
        if readLDR(LDRPin) < 100 {
            sunny = true
        } else {
            sunny = false
        }
        result := StationResult{
            Temperature : temp,
            Humidity : hum,
            Pressure : readBMP180(1),
            IsWet : readRain(RainPin),
            SunValue : sunny,
        }

        bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
            Database:  influx.Database,
            Precision: "ms",
        })

        // Create a point and add to batch
        tags := map[string]string{}
        fields := map[string]interface{}{
            "temperature": result.Temperature,
            "pressure": result.Pressure,
            "Humidity":   result.Humidity,
            "Rain" : result.IsWet,
            "Sun" : result.SunValue,
        }
        pt, err := client.NewPoint(influx.Series, tags, fields, time.Now())
        if err != nil {
            logger.Error.Println(err)
        }
        bp.AddPoint(pt)

        logger.Info.Printf("Going to sleep for %d minutes", 10)
        time.Sleep(10 * time.Minute)
    }
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
        count += 1
        j, err = ldr.Read()
        if err != nil {
            logger.Error.Panicln(err)
        }
    }
    logger.Info.Printf("Sun Value: %d", count)
    return count
}

func readDHT22(port int) (float32, float32) {
    temp, hum, retried, err := dht.ReadDHTxxWithRetry(dht.DHT22, port, true, 10)
    if err != nil {
        log.Fatal(err)
    }
    logger.Info.Printf("T= %v°C, H= %v%% (retried %d times)\n", temp, hum, retried)
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

    logger.Info.Printf("Rain= %t\n", !dry)
    return !dry
}

func readBMP180(port byte) int {
    if err := embd.InitI2C(); err != nil {
        logger.Error.Panicln(err)
    }
    defer embd.CloseI2C()

    bus := embd.NewI2CBus(port)

    baro := bmp180.New(bus)
    defer baro.Close()

    pressure, err := baro.Pressure()
    if err != nil {
        logger.Error.Panicln(err)
    }

    logger.Info.Printf("Pressure is %v hPa\n", pressure)
    return pressure
}
