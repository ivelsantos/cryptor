package crypt

import (
	"fmt"
	"github.com/ivelsantos/cryptor/models"
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
