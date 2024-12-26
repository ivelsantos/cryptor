package backtesting_test

import (
	"testing"

	"github.com/ivelsantos/cryptor/backtesting"
	"github.com/ivelsantos/cryptor/models"
)

var code string = `
		let a = @Mean(window_size = 30)
		if a > 0
			Buy()
		end
	`

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

	// Injecting testingcode
	algos[0].Buycode = code

	err = backtesting.BackTesting(algos[0], 1)
	if err != nil {
		t.Error(err)
	}
}
