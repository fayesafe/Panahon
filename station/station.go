package station

import (
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

type dbClient interface {
	QueryAll(offset string) (*client.Response, error)
	QueryInterval(low string, high string) (*client.Response, error)
	QueryAverage(col string, interval string, offset string, end string) (*client.Response, error)
	QueryMax(col string, interval string, offset string, end string) (*client.Response, error)
}

func StartReadRoutine(influx *database.DBClient, dhtPin int, ldrPin int, rainPin int) {
	// init GPIOs
	if err := embd.InitGPIO(); err != nil {
		logger.Error.Panicln(err)
	}
	defer embd.CloseGPIO()
	// init I2C
	if err := embd.InitI2C(); err != nil {
		logger.Error.Panicln(err)
	}
	defer embd.CloseI2C()
	// init rainSensor
	pin, err := embd.NewDigitalPin(rainPin)
	if err != nil {
		logger.Error.Panicln(err)
	}
	defer pin.Close()
	rainSensor := watersensor.New(pin)
	// init LDR
	ldr, err := embd.NewDigitalPin(ldrPin)
	if err != nil {
		logger.Error.Panicln(err)
	}
	defer ldr.Close()
	// init bmp180
	bmp := bmp180.New(embd.NewI2CBus(1))
	defer bmp.Close()

	for {
		logger.Info.Println("Measurements started...")
		go ReadSensors(influx, dhtPin, ldr, rainSensor, bmp)
		logger.Info.Printf("Going to sleep for %d minutes", 10)
		time.Sleep(10 * time.Minute)
	}
}

func ReadSensors(influx *database.DBClient, dhtPin int, ldr embd.DigitalPin, rainSensor *watersensor.WaterSensor, bmp *bmp180.BMP180) {
	fields := map[string]interface{}{}
	tags := map[string]string{}

	if wet, err := readRainSensor(rainSensor); err == nil {
		fields["rain"] = wet
	}
	if temp, hum, err := readDHT22(dhtPin); err == nil {
		fields["temperature"] = temp
		fields["humidity"] = hum
	}
	if temp, pressure, err := readBMP180(bmp); err == nil {
		if val, ok := fields["temperature"]; ok {
			fields["temperature"] = (temp + val.(float32)) / 2
		} else {
			fields["temperature"] = temp
		}
		fields["pressure"] = pressure
	}
	if sun, err := readLDR(ldr); err == nil {
		fields["sun"] = sun
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

func readRainSensor(rainSensor *watersensor.WaterSensor) (bool, error) {
	wet, err := rainSensor.IsWet()

	if err != nil {
		logger.Error.Println(err)
	} else {
		logger.Info.Printf("Rain=%t (Rain Sensor)\n", wet)
	}

	return wet, err
}

func readDHT22(port int) (float32, float32, error) {
	temp, hum, retried, err := dht.ReadDHTxxWithRetry(dht.DHT22, port, true, 10)

	if err != nil {
		logger.Error.Println(err)
	} else {
		logger.Info.Printf("Temperature=%v°C, Humidity=%v%% (retried %d times) (DHT22)\n", temp, hum, retried)
	}

	return temp, hum, err
}

func readLDR(ldr embd.DigitalPin) (int, error) {
	if err := ldr.SetDirection(embd.Out); err != nil {
		logger.Error.Println(err)
		return 0, err
	}
	if err := ldr.Write(embd.Low); err != nil {
		logger.Error.Println(err)
		return 0, err
	}

	time.Sleep(100 * time.Millisecond)

	if err := ldr.SetDirection(embd.In); err != nil {
		logger.Error.Println(err)
		return 0, err
	}

	count := 0
	j, err := ldr.Read()
	if err != nil {
		logger.Error.Println(err)
		return 0, err
	}
	for j == embd.Low {
		count++
		j, err = ldr.Read()
		if err != nil {
			logger.Error.Println(err)
			return 0, err
		}
	}

	logger.Info.Printf("Sun Value=%d (LDR)", count)
	return count, nil
}

func readBMP180(bmp *bmp180.BMP180) (float32, int, error) {
	temp, err := bmp.Temperature()
	if err != nil {
		logger.Error.Println(err)
		return 0, 0, err
	}
	logger.Info.Printf("Temperature=%v°C (BMP180)\n", temp)

	pressure, err := bmp.Pressure()
	if err != nil {
		logger.Error.Println(err)
		return float32(temp), 0, err
	}

	logger.Info.Printf("Pressure=%vhPa\n", pressure)
	return float32(temp), pressure, nil
}
