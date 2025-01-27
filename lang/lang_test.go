package lang

import (
	"github.com/ivelsantos/cryptor/models"
	"testing"
)

var tests_2 = []struct {
	code string
	err  error
	exp  []string
}{
	{`if 45 < 46 and 2 == 2
		Buy()
	end`, nil, []string{"buy"}},
	{`if 45 < 46 and 2 != 2
		Buy()
	end`, nil, []string{}},
	{`if 45 < 46 or 2 != 2
		Buy()
	end`, nil, []string{"buy"}},
	{`if 45 < 46 or 2 != 2
		Buy()
	end`, nil, []string{"buy"}},
	{`if 45 < 46 or 2 != 2 or 1 > 3
		Buy()
	end`, nil, []string{"buy"}},
	{`if 45 > 46 and 2 != 2 or 1 < 3
		Buy()
	end`, nil, []string{"buy"}},
	{`if 45 < 46 or 2 == 2 and 1 > 3
		Buy()
	end`, nil, []string{"buy"}},
	{`if (45 < 46 or 2 == 2) and 1 > 3
		Buy()
	end`, nil, []string{}},
}

var tests = []struct {
	code string
	err  error
	exp  []string
}{
	{`if 45 < 46
		let b = 42
	end`, nil, []string{}},

	{`let a = 45`, nil, []string{}},

	{`if true == true
		if false != true
			let c = 5
		end
	end`, nil, []string{}},

	{`if 45 < 56
		Sell()
		let a = 33
	end`, nil, []string{"sell"}},

	{`if 45 < 56
		Sell()
		let a = 33
		Buy()
	end`, nil, []string{"sell", "buy"}},
	{`if 45 < 56
		Sell()
		let a = 33
		Buy()
	end
	Buy()
	Sell()`, nil, []string{"sell", "buy", "buy", "sell"}},

	{`if 2 < 1
		Sell()
		let a = 33
	end`, nil, []string{}},

	{`if 2 < 3
		Sell()
		let a = 33
		Buy()
		if 4 > 10
			Sell()
		end
		Buy()
	end
	Sell()`, nil, []string{"sell", "buy", "buy", "sell"}},

	{`if 2 < 3
		Sell()
		let a = 33
		Buy()
		if 4 < 10
			Sell()
		end
		Buy()
	end
	Sell()`, nil, []string{"sell", "buy", "sell", "buy", "sell"}},

	{`let a = 14
	let b = 15
	if a <= b
		Buy()
	end`, nil, []string{"buy"}},

	{`let a = 14
	let b = 15
	if a >= b
		Buy()
	end`, nil, []string{}},

	{`let a = @Price
	let b = 0
	if a > b
		Buy()
	end`, nil, []string{"buy"}},

	// {`let a = @Max(window_size = 14)
	// let b = 0
	// if a > b
	// 	Buy()
	// end`, nil, []string{"buy"}},
	// {`let a = @Max(window_size = 140, lag = 0)
	// let b = 0
	// if a > b
	// 	Buy()
	// end`, nil, []string{"buy"}},
	// {`let a = @Min(window_size = 7)
	// let b = 0
	// if a > b
	// 	Buy()
	// end`, nil, []string{"buy"}},
	// {`let a = @Min(window_size = 7, lag = 4)
	// let b = 0
	// if a > b
	// 	Buy()
	// end`, nil, []string{"buy"}},
	// {`let a = @Min(window_size = 7, lag = 0)
	// let b = 0
	// if a > b
	// 	Buy()
	// end`, nil, []string{"buy"}},
	// {`let a = @Max(window_size = 1, lag = 3)
	// let b = 0
	// if a > b
	// 	Buy()
	// end`, nil, []string{"buy"}},

	// {`let a = @Mean(window_size = 7)
	// let b = 0
	// if a > b
	// 	Buy()
	// end`, nil, []string{"buy"}},
	// {`let a = @Mean(window_size = 0)
	// let b = 0
	// if a > b
	// 	Buy()
	// end`, nil, []string{}},
	// {`let a = @Mean(window_size = 100)
	// let b = @Min(window_size = 100)
	// let c = @Max(window_size = 100)
	// if a > b
	// 	if a < c
	// 		Buy()
	// 	end
	// end`, nil, []string{"buy"}},

	// {`let a = @Median(window_size = 7, lag = 2)
	// let b = 0
	// if a > b
	// 	Buy()
	// end`, nil, []string{"buy"}},
	// {`let a = @Median(window_size = 0)
	// let b = 0
	// if a > b
	// 	Buy()
	// end`, nil, []string{}},
	// {`let a = @Median(window_size = 100)
	// let b = @Min(window_size = 100)
	// let c = @Max(window_size = 100)
	// if a > b
	// 	if a < c
	// 		Buy()
	// 	end
	// end`, nil, []string{"buy"}},

	// {`let a = @Range(window_size = 7)
	// let b = 0
	// if a > b
	// 	Buy()
	// end`, nil, []string{"buy"}},
	// {`let a = @Range(window_size = 14, lag = 3)
	// let b = 0
	// if a > b
	// 	Buy()
	// end`, nil, []string{"buy"}},
	// {`let a = @Range(window_size = 20)
	// let b = @Mean(window_size = 20)
	// let c = @Min(window_size = 20)
	// if a < b
	// 	if a < c
	// 		Buy()
	// 	end
	// end`, nil, []string{"buy"}},

	// {`let a = @Std(window_size = 25)
	// let b = 0
	// if a > b
	// 	Buy()
	// end`, nil, []string{"buy"}},
	// {`let a = @Std(window_size = 54, lag = 3)
	// let b = 0
	// if a > b
	// 	Buy()
	// end`, nil, []string{"buy"}},

	// {`let a = @Var(window_size = 25)
	// let b = 0
	// if a > b
	// 	Buy()
	// end`, nil, []string{"buy"}},
	// {`let a = @Var(window_size = 54, lag = 3)
	// let b = 0
	// if a > b
	// 	Buy()
	// end`, nil, []string{"buy"}},

	{`let a = @Mean(window_size = 2000, lag = 10)
	let b = 0
	if a > b
		Buy()
	end`, nil, []string{"buy"}},
}

var tests_3 = []struct {
	code string
	err  error
	exp  []string
}{
	{`let a = 0.15 + 0.15
	let b = 0.1 + 0.2
	if a == b
		Buy()
	end`, nil, []string{"buy"}},
	{`let a = "ok"
	if a == "ok"
		Buy()
	end`, nil, []string{"buy"}},
	{`let a = "ok"
	if a == "oka"
		Buy()
	end`, nil, []string{}},
	{`let a = "ok"
	if a == "o"
		Buy()
	end`, nil, []string{}},

	{`let a = "ok"
	//let a = "no"
	if a == "ok"
		Buy()
	end`, nil, []string{"buy"}},
	{`let a = "ok"
	if a == "ok"
		// Buy()
	end`, nil, []string{}},
	{`let a = 5
	let b = 5
	if a == b
		Buy()
	end`, nil, []string{"buy"}},
	{`let a = 5.1
	let b = 5.2
	if a == b
		Buy()
	end`, nil, []string{}},
	{`let a = 5.0001
	let b = 5.0002
	if a == b
		Buy()
	end`, nil, []string{}},
	{`let a = 5.0001
	let b = 5.0001
	if a == b
		Buy()
	end`, nil, []string{"buy"}},
	{`let a = 5.00000001
	let b = 5.00000002
	if a == b
		Buy()
	end`, nil, []string{}},
	{`let a = 5.00000001
	let b = 5.00000001
	if a == b
		Buy()
	end`, nil, []string{"buy"}},
	// {`let a = @Mean(ws = 14, lag = 20)
	// let b = @Mean(window_size = 14, lg = 20)
	// if a == b
	// 	Buy()
	// end`, nil, []string{"buy"}},
	// {`let a = @Mean(ws = 14)
	// let b = @Mean(window_size = 14, lg = 19)
	// if a == b
	// 	Buy()
	// end`, nil, []string{}},
}

var tests_4 = []struct {
	code string
	err  error
	exp  []string
}{
	{`a = 0.15 + 0.15
	b = 0.1 + 0.2
	if a == b
		Buy(percentage = 10)
	end`, nil, []string{"buy"}},
	{`a = 0.15 + 0.25
	b = 0.1 + 0.2
	if a == b
		Buy(percentage = 10)
	end`, nil, []string{}},
}

func TestExpressions(t *testing.T) {
	err := models.InitDB("../algor.db")
	if err != nil {
		t.Errorf("Error on database init: %v\n", err)
	}

	algos, err := models.GetAllAlgos()
	if err != nil {
		t.Errorf("Failed to get algos: %v", err)
	}

	algos[0].State = "lang_test"
	optAlgo := GlobalStore("Algo", algos[0])
	optTest := GlobalStore("Test", struct{}{})
	for _, test := range tests_4 {
		res, err := Parse("", []byte(test.code), optAlgo, optTest)
		if test.err != nil {
			if err == nil {
				t.Errorf("Wrong parsing in %s\nIt should not parse", test.code)
			}
		} else {
			if err != nil {
				t.Errorf("Parsing error: %v\n", err)
			}
		}
		if len(res.([]string)) == len(test.exp) {
			for i, exp := range test.exp {
				if exp != res.([]string)[i] {
					t.Errorf("Wrong result in: %s\nExpected %v, got %v\n", test.code, test.exp, res)
					break
				}
			}
		} else {
			t.Errorf("Wrong result in: %s\nExpected %v, got %v\n", test.code, test.exp, res)
		}
	}
}
