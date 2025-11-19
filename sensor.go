package veml6031

import (
	"fmt"

	"github.com/swdee/go-i2c"
)

const (
	// I2C slave address for sensor model codes
	AddressVEML6031X00  = 0x29
	AddressVEML60311X00 = 0x10
	// Sensor ID code for ID_L
	SensorID = 0x01
)

// ProximityPersistance defines the ambient persistance types
type AmbientPersistance uint8

const (
	AmbientPersistance1 AmbientPersistance = 1
	AmbientPersistance2 AmbientPersistance = 2
	AmbientPersistance4 AmbientPersistance = 4
	AmbientPersistance8 AmbientPersistance = 8
)

// GainFactor defines the gain factor type
type GainFactor uint8

const (
	Gain1X        GainFactor = ALS_GAIN_1X
	Gain2X        GainFactor = ALS_GAIN_2X
	GainTwoThirds GainFactor = ALS_GAIN_TWO_THIRDS
	GainHalf      GainFactor = ALS_GAIN_HALF
)

// PhotoDiodeSize defines the photo diode effective size
type PhotoDiodeSize uint8

const (
	PD1 PhotoDiodeSize = PD_DIV4_14
	PD4 PhotoDiodeSize = PD_DIV4_44
)

// registerConf caches the register configuration settings so these can be
// used in the Lux light calculation
type registerConf struct {
	// interval is the ALS_IT_* value
	interval interval
	// gain is the ALS_GAIN_* value
	gain gain
	// pd is the PD_DIV4_* value
	pd pd
}

// Sensor defines the sensor device
type Sensor struct {
	// i2c bus connection
	i2c *i2c.Options
	// vals are the cached register configuration values
	vals registerConf
}

// NewSensor returns a driver instance for the given sensor at the I2C device
// and address
func NewSensor(dev string, addr uint8) (*Sensor, error) {

	i2c, err := i2c.New(addr, dev)

	if err != nil {
		return nil, fmt.Errorf("i2c bus error: %w", err)
	}

	s := &Sensor{
		i2c: i2c,
	}

	check := s.i2c.GetAddr()

	if check == 0 {
		return nil, fmt.Errorf("I2C device is not initiated")
	}

	id, err := s.GetID()

	if err != nil {
		return nil, fmt.Errorf("error getting sensor ID: %w", err)
	}

	if id != SensorID {
		return nil, fmt.Errorf("unexpected sensor ID value: 0x%02X, wanted 0x%02X", id, SensorID)
	}

	if err := s.init(); err != nil {
		return nil, fmt.Errorf("error initializing sensor: %w", err)
	}

	return s, nil
}

// init initializes the sensor at start up
func (s *Sensor) init() error {

	// clear SD bit in ALS_CONF0 (bandgap + LDO ON)
	if err := s.bitMask(ALS_CONF0, SD_MASK, SD_ON); err != nil {
		return fmt.Errorf("error clearing SD bit: %w", err)
	}

	// clear ALS_IR_SD bit in ALS_CONF1 (ALS + IR channels ON)
	if err := s.bitMask(ALS_CONF1, ALS_IR_SD_MASK, ALS_IR_SD_ON); err != nil {
		return fmt.Errorf("error clearing ALS_IR_SD bit: %w", err)
	}

	// set ALS_CAL = 1 as required by datasheet
	if err := s.bitMask(ALS_CONF1, ALS_CAL_MASK, ALS_CAL_ENABLE); err != nil {
		return fmt.Errorf("error setting ALS_CAL bit: %w", err)
	}

	// set sensor to 100ms
	if err := s.SetAmbientIntegrationTime(100); err != nil {
		return fmt.Errorf("error setting ambient integration time: %w", err)
	}

	// set ALS high threshold window to maximum
	if err := s.SetALSHighThreshold(0xFF); err != nil {
		return fmt.Errorf("error setting ALS high threshold: %w", err)
	}

	// set ALS low threshold window setting to minimum
	if err := s.SetALSLowThreshold(0x00); err != nil {
		return fmt.Errorf("error setting ALS low threshold: %w", err)
	}

	// set gain factor
	if err := s.SetALSGain(Gain1X); err != nil {
		return fmt.Errorf("error setting ALS gain: %w", err)
	}

	if err := s.SetPDSize(PD1); err != nil {
		return fmt.Errorf("error setting PD size: %w", err)
	}

	// reset interrupt
	if _, err := s.GetAmbientInterrupt(); err != nil {
		return fmt.Errorf("error resetting interrupt: %w", err)
	}

	return nil
}

// GetID gets the sensors ID
func (s *Sensor) GetID() (uint8, error) {
	return s.readCommandLower(ID_L)
}

// readCommand writes command to sensor and reads the response
func (s *Sensor) readCommand(commandCode byte) (uint16, error) {

	readBuf := make([]byte, 2)

	if _, _, err := s.i2c.WriteThenReadBytes([]byte{commandCode}, readBuf); err != nil {
		return 0, err
	}

	// combine the two bytes into a 16-bit value
	return uint16(readBuf[1])<<8 | uint16(readBuf[0]), nil
}

// readCommandLower reads the lower byte for the given command code address
func (s *Sensor) readCommandLower(commandCode byte) (byte, error) {

	commandValue, err := s.readCommand(commandCode)

	if err != nil {
		return 0, err
	}

	return byte(commandValue & 0xFF), nil
}

// readCommandUpper reads the upper byte for the given command code address
func (s *Sensor) readCommandUpper(commandCode byte) (byte, error) {

	commandValue, err := s.readCommand(commandCode)

	if err != nil {
		return 0, err
	}

	return byte(commandValue >> 8), nil
}

// writeCommand writes a 16-bit value to the given command code location
func (s *Sensor) writeCommand(commandCode byte, value uint16) error {

	buf := []byte{commandCode, byte(value & 0xFF), byte(value >> 8)}

	if _, err := s.i2c.WriteBytes(buf); err != nil {
		return err
	}

	return nil
}

// writeCommandLower writes to the lower byte without affecting the upper byte
// for the given command code address
func (s *Sensor) writeCommandLower(commandCode byte, newValue byte) error {

	commandValue, err := s.readCommand(commandCode)

	if err != nil {
		return err
	}

	commandValue &= 0xFF00           // Remove lower 8 bits
	commandValue |= uint16(newValue) // Mask in

	return s.writeCommand(commandCode, commandValue)
}

// writeCommandUpper writew to the upper byte without affecting the lower byte
// for the given command code address
func (s *Sensor) writeCommandUpper(commandCode byte, newValue byte) error {

	commandValue, err := s.readCommand(commandCode)

	if err != nil {
		return err
	}

	commandValue &= 0x00FF                // Remove upper 8 bits
	commandValue |= uint16(newValue) << 8 // Mask in

	return s.writeCommand(commandCode, commandValue)
}

// bitMask reads a value from a register, masks it, then writes it back
func (s *Sensor) bitMask(commandAddress byte, mask byte, thing byte) error {

	var registerContents byte
	var err error

	registerContents, err = s.readCommandLower(commandAddress)

	if err != nil {
		return err
	}

	// zero-out the portions of the register we're interested in
	registerContents &= mask

	// mask in new thing
	registerContents |= thing

	// change contents
	err = s.writeCommandLower(commandAddress, registerContents)

	return err
}

// SetALSHighThreshold is the value the ambient light sensor (ALS) must go
// above to trigger an interrupt
func (s *Sensor) SetALSHighThreshold(threshold uint16) error {
	return s.writeCommand(ALS_WH_L, threshold)
}

// SetALSLowThreshold is the value the ambient light sensor (ALS) must go
// below to trigger an interrupt
func (s *Sensor) SetALSLowThreshold(threshold uint16) error {
	return s.writeCommand(ALS_WL_L, threshold)
}

// SetAmbientIntegrationTime sets the integration time for the ambient light
// sensor in the number of milliseconds.
func (s *Sensor) SetAmbientIntegrationTime(timeValue uint16) error {

	if timeValue >= 400 {
		timeValue = uint16(ALS_IT_400MS)
	} else if timeValue >= 200 {
		timeValue = uint16(ALS_IT_200MS)
	} else if timeValue >= 100 {
		timeValue = uint16(ALS_IT_100MS)
	} else if timeValue >= 50 {
		timeValue = uint16(ALS_IT_50MS)
	} else if timeValue >= 25 {
		timeValue = uint16(ALS_IT_25MS)
	} else if timeValue >= 12 {
		timeValue = uint16(ALS_IT_12MS5)
	} else if timeValue >= 6 {
		timeValue = uint16(ALS_IT_6MS25)
	} else {
		timeValue = uint16(ALS_IT_3MS125)
	}

	err := s.bitMask(ALS_CONF0, ALS_IT_MASK, byte(timeValue))

	if err == nil {
		s.vals.interval = interval(timeValue)
	}

	return err
}

// GetAmbient reads the ambient light value. Values range from 0 to 65535
// where 0 is dark and 65535 is a bright light source.
func (s *Sensor) GetAmbient() (uint16, error) {
	return s.readCommand(ALS_DATA_L)
}

// GetIR reads the infrared light value.  The higher the value the more IR
// light is present
func (s *Sensor) GetIR() (uint16, error) {
	return s.readCommand(IR_DATA_L)
}

// EnableInterrupt enables the sensor interrupt.  When no event has occurred yet
// the interrupt is held HIGH.  Once a threshold has been triggered the interrupt
// goes low.   The interrupt needs to be reset by calling GetAmbientInterrupt()
// to trigger again.
func (s *Sensor) EnableInterrupt() error {
	return s.bitMask(ALS_CONF0, ALS_INT_EN_MASK, ALS_INT_EN_ENABLE)
}

// DisableInterrupt disables the interrupt functionality
func (s *Sensor) DisableInterrupt() error {
	return s.bitMask(ALS_CONF0, ALS_INT_EN_MASK, ALS_INT_EN_DISABLE)
}

// SetAmbientInterruptPersistance sets the Ambient interrupt persistance value
// The ALS persistence function (ALS_PERS, 1, 2, 4, 8) helps to avoid
// false trigger of the ALS INT. It defines the amount of consecutive hits
// needed in order for a ALS interrupt event to be triggered.
// valid values are  ALS_PERS_[1,2,4,8]
func (s *Sensor) SetAmbientInterruptPersistance(val AmbientPersistance) error {

	var persValue uint8

	if val == AmbientPersistance1 {
		persValue = ALS_PERS_1
	} else if val == AmbientPersistance2 {
		persValue = ALS_PERS_2
	} else if val == AmbientPersistance4 {
		persValue = ALS_PERS_4
	} else {
		// AmbientPersistance8
		persValue = ALS_PERS_8
	}

	return s.bitMask(ALS_CONF1, ALS_PERS_MASK, persValue)
}

// GetAmbientInterrupt reads the ambient light (ALS) interrupt register value
// and returns a zero value if no interrupt event has occurred yet.  Or it
// returns ALS_IF_H (High) or ALS_IF_L (Low) to indicate which ambient threshold
// was triggered.  By called this function the interrupt gets reset and will
// trigger again on next threshold event.
func (s *Sensor) GetAmbientInterrupt() (uint8, error) {
	return s.readCommandLower(ALS_INT)
}

// SetALSGain sets the gain factor.  valid values
func (s *Sensor) SetALSGain(val GainFactor) error {

	err := s.bitMask(ALS_CONF1, ALS_GAIN_MASK, uint8(val))

	if err == nil {
		s.vals.gain = gain(val)
	}

	return err
}

// SetPDSize sets the photo diode effective size
func (s *Sensor) SetPDSize(val PhotoDiodeSize) error {
	err := s.bitMask(ALS_CONF1, PD_DIV4_MASK, uint8(val))

	if err == nil {
		s.vals.pd = pd(val)
	}

	return err
}

// GetLux returns the lux light level for the current ALS sensor reading
func (s *Sensor) GetLux() (float32, error) {

	als, err := s.GetAmbient()

	if err != nil {
		return 0.0, err
	}

	// check if map key values exist
	intervals, ok := luxVals[s.vals.interval]

	if !ok {
		return 0.0, fmt.Errorf("unknown Interval in lux map")
	}

	pds, ok := intervals[s.vals.pd]

	if !ok {
		return 0.0, fmt.Errorf("unknown PD in lux map")
	}

	gainVal, ok := pds[s.vals.gain]

	if !ok {
		return 0.0, fmt.Errorf("unknown Gain in lux map")
	}

	// calculate lux
	lux := float32(als) * gainVal

	return lux, nil
}
