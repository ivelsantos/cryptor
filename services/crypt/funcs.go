package crypt

import (
	"fmt"
	"strconv"
)

func (cp *Crypt) GetCloseMean(args []any) (float64, error) {
	if len(args) != 2 {
		return 0, fmt.Errorf("Usage: @Mean(window f64)")
	}

	window := int(args[0].(float64))
	var val float64
	for i := len(cp.klines); i > (len(cp.klines) - window); i-- {
		v, err := strconv.ParseFloat(cp.klines[i-1].Close, 64)
		if err != nil {
			return 0, err
		}
		val += v
	}

	mean := val / float64(window)
	// fmt.Printf("window: %v\tmean: %v\tprice: %v\n\n", window, mean, cp.vals["@Price"])
	return mean, nil
}
