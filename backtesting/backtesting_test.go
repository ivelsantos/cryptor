package backtesting_test

import (
	"testing"

	"github.com/ivelsantos/cryptor/backtesting"
	"github.com/ivelsantos/cryptor/models"
)

func TestBacktesting(t *testing.T) {
	err := models.InitDB("../algor.db")
	if err != nil {
		t.Error(err)
	}

	algos, err := models.GetAllAlgos()
	if err != nil {
		t.Errorf("Failed to get algos: %v", err)
		return
	}

	err = backtesting.BackTesting(algos[0], 2)
	if err != nil {
		t.Error(err)
	}
}
