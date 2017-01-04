package utils

import (
	"math"
	"math/rand"
	"sync"
	"time"
)

// Sampler represents a sampler instance, which keeps track of the number of
// traces per second to be sampled and fallback rate for additional sampling.
// Additionally, a sampler instance determines whether a given trace should be
// sampled based on 'fixedTarget' and 'fallbackRate' settings.
type Sampler struct {
	fallbackRate   float64
	fixedTarget    uint32
	usedThisSecond uint32
	thisSecond     uint64
	rand           *rand.Rand

	sync.RWMutex
}

// NewSampler creates a new sampler with the specified configuration.
func NewSampler(fixedTarget uint32, fallbackRate float64) *Sampler {
	return &Sampler{
		fixedTarget:  fixedTarget,
		fallbackRate: fallbackRate,
		rand:         rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// IsSampled determines whether a given trace should be sampled.
func (s *Sampler) IsSampled() bool {
	now := uint64(math.Floor(float64(time.Now().UnixNano()) /
		float64(time.Second.Nanoseconds())))

	s.RLock()
	thisSecond := s.thisSecond
	usedThisSecond := s.usedThisSecond
	fixedTarget := s.fixedTarget
	fallbackRate := s.fallbackRate
	s.RUnlock()

	if now != thisSecond {
		s.Lock()
		s.usedThisSecond = 0
		s.thisSecond = now
		s.Unlock()
	}

	if usedThisSecond >= fixedTarget || fixedTarget == 0 {
		if fallbackRate == 0 {
			return false
		}

		return s.rand.Float64() < fallbackRate
	}

	s.Lock()
	s.usedThisSecond++
	s.Unlock()

	return true
}
