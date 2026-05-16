package gpu

import (
	"github.com/fumiama/gozel/gozel"
	"github.com/fumiama/gozel/ze"
)

// SamplerCreateNormalizedClamp creates a sampler with normalized coordinates,
// clamp address mode and the specified filter mode.
func SamplerCreateNormalizedClamp(filtermode gozel.ZeSamplerFilterMode) (ze.SamplerHandle, error) {
	return g().ctx.SamplerCreate(
		g().dev, gozel.ZE_SAMPLER_ADDRESS_MODE_CLAMP,
		filtermode, 1,
	)
}
