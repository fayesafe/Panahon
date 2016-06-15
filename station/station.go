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

func read_ldr(port int) {
    if err := embd.InitGPIO(); err != nil {
        panic(err)
    }
    defer embd.CloseGPIO()

    ldr,err := embd.NewDigitalPin(port)
    if err != nil {
        panic(err)
    }
    defer ldr.Close()

    for i := 0; i < 200; i++ {
        if err := ldr.SetDirection(embd.Out); err != nil {
            panic(err)
        }
        if err := ldr.Write(embd.Low); err != nil {
            panic(err)
        }

        time.Sleep(100 * time.Millisecond)

        if err := ldr.SetDirection(embd.In); err != nil {
            panic(err)
        }

        count := 0
        j,_ := ldr.Read()
        for j == embd.Low {
            count += 1
            j,_ = ldr.Read()
        }

        fmt.Printf("LDR value is: %v\n", count)
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
