package crypt

import (
	"fmt"
	"github.com/ivelsantos/cryptor/models"
	"github.com/ivelsantos/cryptor/services/crypt/values"
)

func GetFuncValue(algo models.Algor, funcName string, args any) (float64, error) {
	argsSlice := args.([]any)
	switch funcName {

	case "@Mean":
		val, err := GetCloseMean(algo, argsSlice)
		if err != nil {
			return 0, err
		}
		return val, nil

	default:
		return 0, fmt.Errorf("Function %s does not exists", "@"+funcName)
	}
}

func GetCryptValue(algo models.Algor, key string) (float64, error) {
	switch key {
	case "Price":
		price, err := values.GetPrice(algo.BaseAsset + algo.QuoteAsset)
		if err != nil {
			return 0, err
		}
		return price, nil
	default:
		return 0, fmt.Errorf("Value not found!")
	}
}
