package crypt

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ivelsantos/cryptor/models"
)

func GetFuncValue(algo models.Algor, funcName string, index any, args string) (float64, error) {
	n := index.(int)

	argsSlice := strings.Split(args, ",")
	argsSlice[len(argsSlice)-1] = strings.ReplaceAll(argsSlice[len(argsSlice)-1], ")", "")
	arguments := make(map[string]string)
	for _, arg := range argsSlice {
		sls := strings.Split(arg, "=")
		if len(sls) < 2 {
			continue
		}
		arguments[strings.Trim(sls[0], " ")] = strings.Trim(sls[1], " ")
	}
	if algo.State == "backtesting" {
		arguments["backindex"] = strconv.Itoa(n)
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

	case "@Price":
		val, err := GetPrice(algo, arguments)
		if err != nil {
			return 0, err
		}
		return val, nil

	default:
		return 0, fmt.Errorf("Function %s does not exists", funcName)
	}
}
