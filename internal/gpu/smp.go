package gpu

import (
	"github.com/fumiama/gozel/gozel"
	"github.com/fumiama/gozel/ze"
)

// SamplerCreateNormalizedLinearClamp is the most commonly used sampler.
func SamplerCreateNormalizedLinearClamp() (ze.SamplerHandle, error) {
	return ctx.SamplerCreate(
		dev, gozel.ZE_SAMPLER_ADDRESS_MODE_CLAMP,
		gozel.ZE_SAMPLER_FILTER_MODE_LINEAR, 1,
	)
}
