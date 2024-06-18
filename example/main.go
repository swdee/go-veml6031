package main

import (
	"flag"
	"github.com/swdee/go-veml6031"
	"log"
	"strconv"
	"time"
)

func main() {

	// disable logging timestamps
	log.SetFlags(0)

	// read in cli flags
	i2cbus := flag.String("b", "/dev/i2c-0", "Path to I2C bus to use")
	addr := flag.String("a", "0x29", "Hex address of sensor on I2C bus, default 0x29")
	flag.Parse()

	var useAddr uint8

	if *addr != "" {
		intValue, err := strconv.ParseInt(*addr, 0, 64)

		if err != nil || intValue > 255 {
			log.Fatalf("Error casting sensor hex address: %v", err)
		}

		useAddr = uint8(intValue)
	}

	sensor, err := veml6031.NewSensor(*i2cbus, useAddr)

	if err != nil {
		log.Fatalf("Error connecting to sensor: %v\n", err)
	}

	configureInterrupt(sensor)
	readValues(sensor)
}

func readValues(sensor *veml6031.Sensor) {

	c := 0

	for {
		time.Sleep(1 * time.Second)

		// read ambient light data
		ambient, err := sensor.GetAmbient()

		if err != nil {
			log.Printf("Failed to read ambient light: %v", err)
			continue
		}

		// read IR light data
		ir, err := sensor.GetIR()

		if err != nil {
			log.Printf("Failed to read IR light: %v", err)
			continue
		}

		lux, err := sensor.GetLux()

		if err != nil {
			log.Printf("Failed to read Lux value: %v", err)
			continue
		}

		if c >= 10 {
			// read interrupt
			checkInterrupt(sensor)
			c = 0
		}

		c++
		log.Printf("Count=%d, Ambient Light: %d, IR: %d, Lux: %.2f\n", c, ambient, ir, lux)
	}
}

// checkInterrupt is a quick way of checking the sensor register every 10 seconds
// to read its state.   This is done so setup of GPIO pin to the read the
// interrupt is not needed in this example.
func checkInterrupt(sensor *veml6031.Sensor) {

	intVal, err := sensor.GetAmbientInterrupt()

	if err != nil {
		log.Printf("error reading ambient interrupt value: %w", err)
		return
	}

	switch intVal {
	case veml6031.ALS_IF_H:
		log.Printf("Interrupt: HIGH threshold triggered")
	case veml6031.ALS_IF_L:
		log.Printf("Interrupt: LOW threshold triggered")
	default:
		log.Printf("Interrupt: no event")
	}
}

func configureInterrupt(sensor *veml6031.Sensor) {

	// NOTE: the default state of the interrupt is to be HIGH.  When the
	// threshold event occurs the interrupt will go LOW.   The interrupt needs
	// to be reset after triggering by calling GetAmbientInterrupt() before
	// it will trigger again

	// if sensor sees a value higher than this, interrupt pin will go LOW
	err := sensor.SetALSHighThreshold(2000)

	if err != nil {
		log.Fatalf("Failed to set proximity high threshold: %v", err)
	}

	// if sensor sees a value lower than this, interrupt pin will go LOW
	err = sensor.SetALSLowThreshold(500)

	if err != nil {
		log.Fatalf("Failed to set proximity low threshold: %v", err)
	}

	// set the number of consecutive ALS values required to trigger the interrupt
	err = sensor.SetAmbientInterruptPersistance(veml6031.AmbientPersistance4)

	if err != nil {
		log.Fatalf("Failed to set ambient persistance: %v", err)
	}

	err = sensor.EnableInterrupt()

	if err != nil {
		log.Fatalf("Failed to enable interrupt mode: %v", err)
	}
}
