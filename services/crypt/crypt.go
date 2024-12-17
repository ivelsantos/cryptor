package crypt

import (
	"fmt"
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

	default:
		return 0, fmt.Errorf("Function %s does not exists", funcName)
	}
}

func GetCryptValue(algo models.Algor, key string) (float64, error) {
	switch key {
	case "@Price":
		price, err := values.GetPrice(algo.BaseAsset + algo.QuoteAsset)
		if err != nil {
			return 0, err
		}
		return price, nil
	default:
		return 0, fmt.Errorf("Value not found!")
	}
}
