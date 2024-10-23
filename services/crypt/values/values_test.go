package values

import (
	"fmt"
	"testing"
)

func TestPrice(t *testing.T) {
	res, err := GetPrice("BTCBRL")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%f\n", res)
}
