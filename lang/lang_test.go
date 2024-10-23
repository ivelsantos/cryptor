package lang

import (
	"github.com/ivelsantos/cryptor/services/crypt"
	"testing"
)

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
	let b = @Mean(14)
	if a < b
		Buy()
	end`, nil, []string{"buy"}},
	{`let a = @Price
	let b = @Mean(1000)
	if a < b
		Buy()
	end`, nil, []string{"buy"}},
}

func TestExpressions(t *testing.T) {
	err := crypt.InitCrypt()
	if err != nil {
		t.Errorf("Error initiating Crypt: %v", err)
	}

	for _, test := range tests {
		res, err := Parse("", []byte(test.code))
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
