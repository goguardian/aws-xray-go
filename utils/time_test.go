package utils

import (
	"math"
	"testing"
	"time"
)

func TestCurrentTimeSecond(t *testing.T) {
	if float64(time.Now().Unix()) > CurrentTimeSecond() {
		t.Error("Calculated current time with non-zero decimal should be greater " +
			"than current time with zero decimal")
	}

	if float64(time.Now().Unix()) != math.Floor(CurrentTimeSecond()) {
		t.Errorf("Floored current time with non-zero decimal should equal the " +
			"current time with zero decimal")
	}
}
