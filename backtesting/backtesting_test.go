package backtesting_test

import (
	"testing"

	"github.com/ivelsantos/cryptor/backtesting"
	"github.com/ivelsantos/cryptor/models"
)

var code string = `
		let a = @Mean(window_size = 30)
		let b = a / @Price
		if b < 0.994
			Buy()
		end
		Stop_loss(0.008)
		Take_profit(0.005)
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

	err = backtesting.BackTesting(algos[0], 2)
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkGetData1(b *testing.B) {
	err := models.InitDB("../algor.db")
	if err != nil {
		b.Error(err)
	}

	algos, err := models.GetAllAlgos()
	if err != nil {
		b.Errorf("Failed to get algos: %v", err)
		return
	}

	// Injecting testingcode
	algos[0].Buycode = code

	for i := 0; i < b.N; i++ {
		err = backtesting.BackTesting(algos[0], 2)
		if err != nil {
			b.Error(err)
		}
	}
	// for i := 0; i < b.N; i++ {
	// 	getDataBench(1, b)
	// }
}
