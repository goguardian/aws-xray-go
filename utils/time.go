package utils

import "time"

// CurrentTimeSecond returns a float64 representation of the current Unix time
// with sub-second decimal accuracy.
func CurrentTimeSecond() float64 {
	return float64(time.Now().UnixNano()) / float64(time.Second.Nanoseconds())
}
