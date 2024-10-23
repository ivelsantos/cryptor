package functions

import (
	"testing"
)

func TestMean(t *testing.T) {
	res, err := cryptMean(3)
	if err != nil {
		t.Errorf("Error wich cryptMean: %v", err)
	}
	_ = res
}
