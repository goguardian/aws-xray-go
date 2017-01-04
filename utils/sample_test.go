package utils

import (
	"testing"
	"time"
)

func TestNewSampler(t *testing.T) {
	const samplesCount = 1000

	tests := []struct {
		fixedRate        uint32
		fallbackRate     float64
		expectSampledMin int
		expectSampledMax int
	}{
		{
			fixedRate:        10,
			fallbackRate:     0,
			expectSampledMin: 20,
			expectSampledMax: 20,
		},
		{
			fixedRate:        0,
			fallbackRate:     0.5,
			expectSampledMin: 400,
			expectSampledMax: 600,
		},
		{
			fixedRate:        10,
			fallbackRate:     1,
			expectSampledMin: samplesCount,
			expectSampledMax: samplesCount,
		},
		{
			fixedRate:    1,
			fallbackRate: 0,
			// Test will span just more than 1 second, so expect 2
			expectSampledMin: 2,
			expectSampledMax: 2,
		},
	}

	for _, test := range tests {
		sampler := NewSampler(test.fixedRate, test.fallbackRate)

		start := time.Now()
		sampled := 0
		count := 0
		for count < samplesCount {
			isSampled := sampler.IsSampled()
			if isSampled {
				sampled++
			}
			count++
			time.Sleep(1 * time.Millisecond)
		}

		if sampled > test.expectSampledMax {
			t.Errorf("Sample count should not be greater than %d, got %d: %v",
				test.expectSampledMax, sampled, time.Since(start))
		}

		if sampled < test.expectSampledMin {
			t.Errorf("Sample count should not be less than %d, got %d",
				test.expectSampledMin, sampled)
		}
	}
}
