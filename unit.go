package gg

import (
	"fmt"
	"math"
	"strings"

	"golang.org/x/image/math/fixed"
)

// Radians converts degrees to radians.
//
// Radians 将角度转换为弧度。
func Radians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// Degrees converts radians to degrees.
//
// Degrees 将弧度转换为角度。
func Degrees(radians float64) float64 {
	return radians * 180 / math.Pi
}

// ParseHexColor parses a hex color string and returns r, g, b, a components.
// Supports 3, 6, and 8 digit hex strings with an optional leading '#'.
//
// ParseHexColor 解析十六进制颜色字符串并返回 r, g, b, a 分量。
// 支持 3、6 和 8 位十六进制字符串，前导 '#' 可选。
func ParseHexColor(x string) (r, g, b, a int) {
	x = strings.TrimPrefix(x, "#")
	a = 255
	switch len(x) {
	case 3:
		format := "%1x%1x%1x"
		_, _ = fmt.Sscanf(x, format, &r, &g, &b)
		r |= r << 4
		g |= g << 4
		b |= b << 4
	case 6:
		format := "%02x%02x%02x"
		_, _ = fmt.Sscanf(x, format, &r, &g, &b)
	case 8:
		format := "%02x%02x%02x%02x"
		_, _ = fmt.Sscanf(x, format, &r, &g, &b, &a)
	}
	return
}

func fixp(x, y float64) fixed.Point26_6 {
	return fixed.Point26_6{
		X: fix(x),
		Y: fix(y),
	}
}

func fix(x float64) fixed.Int26_6 {
	return fixed.Int26_6(math.Round(x * 64))
}

func unfix(x fixed.Int26_6) float64 {
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
