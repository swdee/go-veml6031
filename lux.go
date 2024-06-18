package veml6031

type interval uint16
type pd uint8
type gain uint8

var (
	// lux level multipliers from datasheet
	luxVals = map[interval]map[pd]map[gain]float32{
		ALS_IT_3MS125: {
			PD_DIV4_14: {
				ALS_GAIN_2X:         1.7408,
				ALS_GAIN_1X:         3.4816,
				ALS_GAIN_TWO_THIRDS: 5.2752,
				ALS_GAIN_HALF:       6.9632,
			},
			PD_DIV4_44: {
				ALS_GAIN_2X:         0.4352,
				ALS_GAIN_1X:         0.8704,
				ALS_GAIN_TWO_THIRDS: 1.3188,
				ALS_GAIN_HALF:       1.7408,
			},
		},
		ALS_IT_6MS25: {
			PD_DIV4_14: {
				ALS_GAIN_2X:         0.8704,
				ALS_GAIN_1X:         1.7408,
				ALS_GAIN_TWO_THIRDS: 2.6376,
				ALS_GAIN_HALF:       3.4816,
			},
			PD_DIV4_44: {
				ALS_GAIN_2X:         0.2176,
				ALS_GAIN_1X:         0.4352,
				ALS_GAIN_TWO_THIRDS: 0.6594,
				ALS_GAIN_HALF:       0.8704,
			},
		},
		ALS_IT_12MS5: {
			PD_DIV4_14: {
				ALS_GAIN_2X:         0.4352,
				ALS_GAIN_1X:         0.8704,
				ALS_GAIN_TWO_THIRDS: 1.3188,
				ALS_GAIN_HALF:       1.7408,
			},
			PD_DIV4_44: {
				ALS_GAIN_2X:         0.1088,
				ALS_GAIN_1X:         0.2176,
				ALS_GAIN_TWO_THIRDS: 0.3297,
				ALS_GAIN_HALF:       0.4352,
			},
		},
		ALS_IT_25MS: {
			PD_DIV4_14: {
				ALS_GAIN_2X:         0.2176,
				ALS_GAIN_1X:         0.4352,
				ALS_GAIN_TWO_THIRDS: 0.6594,
				ALS_GAIN_HALF:       0.8704,
			},
			PD_DIV4_44: {
				ALS_GAIN_2X:         0.0544,
				ALS_GAIN_1X:         0.1088,
				ALS_GAIN_TWO_THIRDS: 0.1648,
				ALS_GAIN_HALF:       0.2176,
			},
		},
		ALS_IT_50MS: {
			PD_DIV4_14: {
				ALS_GAIN_2X:         0.1088,
				ALS_GAIN_1X:         0.2176,
				ALS_GAIN_TWO_THIRDS: 0.3297,
				ALS_GAIN_HALF:       0.4352,
			},
			PD_DIV4_44: {
				ALS_GAIN_2X:         0.0272,
				ALS_GAIN_1X:         0.0544,
				ALS_GAIN_TWO_THIRDS: 0.0824,
				ALS_GAIN_HALF:       0.1088,
			},
		},
		ALS_IT_100MS: {
			PD_DIV4_14: {
				ALS_GAIN_2X:         0.0544,
				ALS_GAIN_1X:         0.1088,
				ALS_GAIN_TWO_THIRDS: 0.1648,
				ALS_GAIN_HALF:       0.2176,
			},
			PD_DIV4_44: {
				ALS_GAIN_2X:         0.0136,
				ALS_GAIN_1X:         0.0272,
				ALS_GAIN_TWO_THIRDS: 0.0412,
				ALS_GAIN_HALF:       0.0544,
			},
		},
		ALS_IT_200MS: {
			PD_DIV4_14: {
				ALS_GAIN_2X:         0.0272,
				ALS_GAIN_1X:         0.0544,
				ALS_GAIN_TWO_THIRDS: 0.0824,
				ALS_GAIN_HALF:       0.0188,
			},
			PD_DIV4_44: {
				ALS_GAIN_2X:         0.0068,
				ALS_GAIN_1X:         0.0136,
				ALS_GAIN_TWO_THIRDS: 0.0206,
				ALS_GAIN_HALF:       0.0272,
			},
		},
		ALS_IT_400MS: {
			PD_DIV4_14: {
				ALS_GAIN_2X:         0.0136,
				ALS_GAIN_1X:         0.0272,
				ALS_GAIN_TWO_THIRDS: 0.0412,
				ALS_GAIN_HALF:       0.0544,
			},
			PD_DIV4_44: {
				ALS_GAIN_2X:         0.0034,
				ALS_GAIN_1X:         0.0068,
				ALS_GAIN_TWO_THIRDS: 0.0103,
				ALS_GAIN_HALF:       0.0136,
			},
		},
	}
)
