package point

import (
	"math"

	"golang.org/x/image/math/fixed"
)

// FixedPoint convert (x, y) from float64 to fixed.Point26_6
func FixedPoint(x, y float64) fixed.Point26_6 {
	return fixed.Point26_6{
		X: Fixed(x),
		Y: Fixed(y),
	}
}

// Fixed convert float64 to fixed.Int26_6
func Fixed(x float64) fixed.Int26_6 {
	return fixed.Int26_6(math.Round(x * 64))
}

// DeFixed convert fixed.Int26_6 to float64
func DeFixed(x fixed.Int26_6) float64 {
	const shift, mask = 6, 1<<6 - 1
	if x >= 0 {
		return float64(x>>shift) + float64(x&mask)/64
	}
	x = -x
	if x >= 0 {
		return -(float64(x>>shift) + float64(x&mask)/64)
	}
	return 0
}
