# go-veml6031

go-veml6031 is an I2C driver for the Vishay [VEML6031](https://www.vishay.com/docs/80007/veml6031x00.pdf)
ambient light sensor.



## Usage

To use in your Go project, get the library.

```
go get github.com/swdee/go-veml6031
```



## Example

To read Lux light level from sensor.

```
// initialise sensor
sensor, _ := veml6031.NewSensor("/dev/i2c-0", "0x29")

// read lux light level
lux, _ := sensor.GetLux()

fmt.Printf("Lux = %.2f\n", lux)
```

Note: Error handling has been skipped for brevity.

For reading Ambient Light, Infrared Light, and setting Interrupts see the more 
[complete example here](example/main.go).


## Interrupt Pin

The `INT` pin on sensor should be wired to a GPIO pin of your microprocessor.
Its normal state is High and when it goes Low this signals the threshold values
have been triggered.

Upon receiving a Low signal you need to clear the interrupt state by calling
`GetAmbientInterrupt()` which returns the register value indicating if the 
`ALS_IF_L` (Low) or `ALS_IF_H` (High) threshold value was triggered.

Once `GetAmbientInterrupt()` has been called the `INT` pin resets to High.


## Reference

* [VEML6031 Datasheet](https://www.vishay.com/docs/80007/veml6031x00.pdf)  
* [Application Notes](https://www.vishay.com/docs/80201/designingveml6031x00.pdf)