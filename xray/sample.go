package xray

import (
	"github.com/goguardian/aws-xray-go/segment"
	"github.com/goguardian/aws-xray-go/utils"
)

// SetSampler updates the sampler used for segment sampling.
func SetSampler(fixedTarget uint32, fallbackRate float64) {
	segment.SetSampler(utils.NewSampler(fixedTarget, fallbackRate))
}
