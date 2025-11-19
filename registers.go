package veml6031

const (
	// ALS_CONF_0
	ALS_IT_MASK   = ^uint8((1 << 6) | (1 << 5) | (1 << 4))
	ALS_IT_3MS125 = 0                              // 000
	ALS_IT_6MS25  = 1 << 4                         // 001
	ALS_IT_12MS5  = 1 << 5                         // 010
	ALS_IT_25MS   = (1 << 5) | (1 << 4)            // 011
	ALS_IT_50MS   = 1 << 6                         // 100
	ALS_IT_100MS  = (1 << 6) | (1 << 4)            // 101
	ALS_IT_200MS  = (1 << 6) | (1 << 5)            // 110
	ALS_IT_400MS  = (1 << 6) | (1 << 5) | (1 << 4) // 111

	ALS_AF_MASK    = ^uint8(1 << 3)
	ALS_AF_DISABLE = 0
	ALS_AF_ENABLE  = 1 << 3

	ALS_TRIG_MASK    = ^uint8(1 << 2)
	ALS_TRIG_TRIGGER = 1 << 2

	ALS_INT_EN_MASK    = ^uint8(1 << 1)
	ALS_INT_EN_DISABLE = 0
	ALS_INT_EN_ENABLE  = 1 << 1

	SD_MASK     = ^uint8(1 << 0)
	SD_ON       = 0
	SD_SHUTDOWN = 1 << 0

	// ALS_CONF_1
	ALS_IR_SD_MASK     = ^uint8(1 << 7)
	ALS_IR_SD_ON       = 0
	ALS_IR_SD_SHUTDOWN = 1 << 7

	PD_DIV4_MASK = ^uint8(1 << 6)
	PD_DIV4_44   = 0
	PD_DIV4_14   = 1 << 6

	ALS_GAIN_MASK       = ^uint8((1 << 4) | (1 << 3))
	ALS_GAIN_1X         = 0
	ALS_GAIN_2X         = 1 << 3
	ALS_GAIN_TWO_THIRDS = 1 << 4
	ALS_GAIN_HALF       = (1 << 4) | (1 << 3)

	ALS_PERS_MASK = ^uint8((1 << 2) | (1 << 1))
	ALS_PERS_1    = 0
	ALS_PERS_2    = 1 << 1
	ALS_PERS_4    = 1 << 2
	ALS_PERS_8    = (1 << 2) | (1 << 1)

	ALS_CAL_MASK   = ^uint8(1 << 0)
	ALS_CAL_ENABLE = 1 << 0

	ALS_IF_H = (1 << 1)
	ALS_IF_L = (1 << 2)
)
