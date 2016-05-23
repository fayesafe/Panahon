package station

import (
    "fmt"
    "log"
    "time"

    "github.com/d2r2/go-dht"
    "github.com/kidoman/embd"
    "github.com/kidoman/embd/convertors/mcp3008"
    "github.com/kidoman/embd/sensor/bmp180"
    "github.com/kidoman/embd/sensor/watersensor"

    _ "github.com/kidoman/embd/host/rpi"

    "Panahon/logger"
)

func TestRoutine() {
    for {
        logger.Info.Println("-- CONCURRENCY TEST ALERT --")
        time.Sleep(5000 * time.Millisecond)
    }
}

func read_ldr() {
    if err := embd.InitSPI(); err != nil {
        panic(err)
    }
    defer embd.CloseSPI()

    spiBus := embd.NewSPIBus(embd.SPIMode0, channel, speed, bpw, delay)
    defer spiBus.Close()

    adc := mcp3008.New(mcp3008.SingleMode, spiBus)

    for i := 0; i < 20; i++ {
        time.Sleep(1 * time.Second)
        val, err := adc.AnalogValueAt(0)
        if err != nil {
            panic(err)
        }
        fmt.Printf("analog value is: %v\n", val)
    }
}

func read_dht22(port int) {
    t,h,retried,err := dht.ReadDHTxxWithRetry(dht.DHT22, port, true, 10)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("T= %v°C, H= %v%% (retried %d times)\n", t, h, retried)
}

func read_rain(port int) {
	if err := embd.InitGPIO(); err != nil {
		panic(err)
	}
	defer embd.CloseGPIO()

	pin, err := embd.NewDigitalPin(port)
	if err != nil {
		panic(err)
	}
	defer pin.Close()

    fluidSensor := watersensor.New(pin)

	dry, err := fluidSensor.IsWet()
	if err != nil {
		panic(err)
	}

    fmt.Printf("Rain= %t\n", !dry)
}

func read_bmp180(port byte) {
    if err := embd.InitI2C(); err != nil {
        panic(err)
    }
    defer embd.CloseI2C()

    bus := embd.NewI2CBus(port)

    baro := bmp180.New(bus)
    defer baro.Close()

    temp, err := baro.Temperature()
    if err != nil {
        panic(err)
    }

    fmt.Printf("Temp is %v °C\n", temp)
    pressure, err := baro.Pressure()
    if err != nil {
        panic(err)
    }

    fmt.Printf("Pressure is %v hPa\n", pressure)
    altitude, err := baro.Altitude()
    if err != nil {
        panic(err)
    }

    fmt.Printf("Altitude is %v\n", altitude)
}
