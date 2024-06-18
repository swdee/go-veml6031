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


## Reference

* [VEML6031 Datasheet](https://www.vishay.com/docs/80007/veml6031x00.pdf)  
* [Application Notes](https://www.vishay.com/docs/80201/designingveml6031x00.pdf)