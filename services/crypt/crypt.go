package crypt

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ivelsantos/cryptor/models"
	"github.com/ivelsantos/cryptor/services/crypt/values"
)

func GetFuncValue(algo models.Algor, funcName string, args string) (float64, error) {
	argsSlice := strings.Split(args, ",")
	argsSlice[len(argsSlice)-1] = strings.ReplaceAll(argsSlice[len(argsSlice)-1], ")", "")
	arguments := make(map[string]string)
	for _, arg := range argsSlice {
		sls := strings.Split(arg, "=")
		if len(sls) != 2 {
			return 0, fmt.Errorf("Wrong argument format: %s", arg)
		}
		arguments[strings.Trim(sls[0], " ")] = strings.Trim(sls[1], " ")
	}

	switch funcName {

	case "@Max":
		val, err := GetMaxValue(algo, arguments)
		if err != nil {
			return 0, err
		}
		return val, nil

	case "@Min":
		val, err := GetMinValue(algo, arguments)
		if err != nil {
			return 0, err
		}
		return val, nil

	case "@Range":
		val, err := GetRangeValue(algo, arguments)
		if err != nil {
			return 0, err
		}
		return val, nil

	case "@Mean":
		val, err := GetMeanValue(algo, arguments)
		if err != nil {
			return 0, err
		}
		return val, nil

	case "@Median":
		val, err := GetMedianValue(algo, arguments)
		if err != nil {
			return 0, err
		}
		return val, nil

	case "@Std":
		val, err := GetStdValue(algo, arguments)
		if err != nil {
			return 0, err
		}
		return val, nil

	case "@Var":
		val, err := GetVarValue(algo, arguments)
		if err != nil {
			return 0, err
		}
		return val, nil

	case "@Ema":
		val, err := GetEmaValue(algo, arguments)
		if err != nil {
			return 0, err
		}
		return val, nil

	default:
		return 0, fmt.Errorf("Function %s does not exists", funcName)
	}
}

func GetCryptValue(algo models.Algor, key string, index any) (float64, error) {
	n := index.(int)

	switch key {
	case "@Price":
		if algo.State == "backtesting" {
			value, err := strconv.ParseFloat(models.Backtesting_Data[n].Close, 64)
			if err != nil {
				return 0.0, fmt.Errorf("Error parsing close value: %v", err)
			}
			return value, nil
		}

		price, err := values.GetPrice(algo.BaseAsset + algo.QuoteAsset)
		if err != nil {
			return 0, err
		}
		return price, nil
	default:
		return 0, fmt.Errorf("Value not found!")
	}
}
