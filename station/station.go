package station

import (
	"sync"
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

type Sensors struct {
	dhtPin     int
	ldr        embd.DigitalPin
	rainPin    embd.DigitalPin
	rainSensor *watersensor.WaterSensor
	bmp        *bmp180.BMP180
}

var instance *Sensors
var once sync.Once
var mutexRead sync.Mutex

func InitSensors(dhtPin int, ldrPin int, rainPin int) *Sensors {
    sensors := new(Sensors)
    var err error
	// init GPIOs
	if err := embd.InitGPIO(); err != nil {
		logger.Error.Panicln(err)
	}
	// init I2C
	if err := embd.InitI2C(); err != nil {
		logger.Error.Panicln(err)
	}
	// init rainSensor
	sensors.rainPin, err = embd.NewDigitalPin(rainPin)
	if err != nil {
		logger.Error.Panicln(err)
	}
	sensors.rainSensor = watersensor.New(sensors.rainPin)
	// init LDR
	sensors.ldr, err = embd.NewDigitalPin(ldrPin)
	if err != nil {
		logger.Error.Panicln(err)
	}
	// init bmp180
	sensors.bmp = bmp180.New(embd.NewI2CBus(1))
	// DHT22
	sensors.dhtPin = dhtPin
	logger.Info.Println("Sensors initialized")

    return sensors
}

func (sensors Sensors) Close() {
	embd.CloseGPIO()
	embd.CloseI2C()
	sensors.rainPin.Close()
	sensors.ldr.Close()
	sensors.bmp.Close()
	logger.Info.Println("Sensors closed")
}
func (sensors Sensors) RunReadRoutine(influx database.DBClient, interval time.Duration) {
	for {
		logger.Info.Println("Measurements started...")
		go sensors.Read(influx)
		logger.Info.Printf("Going to sleep for %d minutes", interval)
		time.Sleep(interval * time.Minute)
	}
}

func (sensors Sensors) Read(influx database.DBClient) {
	mutexRead.Lock()
	fields := map[string]interface{}{}
	tags := map[string]string{}

	if wet, err := readRainSensor(sensors.rainSensor); err == nil {
        if wet {
            fields["rain"] = 1
        } else {
            fields["rain"] = 0
        }
	}
	if temp, hum, err := readDHT22(sensors.dhtPin); err == nil {
		fields["temperature"] = temp
		fields["humidity"] = hum
	}
	if temp, pressure, err := readBMP180(sensors.bmp); err == nil {
		if val, ok := fields["temperature"]; ok {
			fields["temperature"] = (temp + val.(float32)) / 2
		} else {
			fields["temperature"] = temp
		}
		fields["pressure"] = pressure
	}
	if sun, err := readLDR(sensors.ldr); err == nil {
		fields["sun"] = sun
	}

	logger.Info.Println("Measurements finished")

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
	influx.Client.Write(bp)
	logger.Info.Println("Measurements written to db")
	mutexRead.Unlock()
}

func readRainSensor(rainSensor *watersensor.WaterSensor) (bool, error) {
	dry, err := rainSensor.IsWet()

	if err != nil {
		logger.Error.Println(err)
	} else {
		logger.Info.Printf("Rain=%t (Rain Sensor)\n", !dry)
	}

	return !dry, err
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
