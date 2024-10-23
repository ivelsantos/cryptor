package values

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

func GetPrice(symbol string) (float64, error) {
	type Price struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	}
	var price Price

	s := "https://api.binance.com/api/v3/ticker/price?symbol=" + symbol
	res, err := http.Get(s)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	err = json.Unmarshal(body, &price)
	if err != nil {
		return 0, err
	}

	p, ok := strconv.ParseFloat(price.Price, 64)
	if ok != nil {
		return 0, ok
	}

	return p, nil
}
