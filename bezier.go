package gg

import (
	"math"
)

func quadraticBezier(x0, y0, x1, y1, x2, y2, ds float64, p []Point) {
	if canUseBezierKernel {
		err := quadraticBezeirGPU(x0, y0, x1, y1, x2, y2, ds, p)
		if err == nil {
			return
		}
		canUseBezierKernel = false
	}
	quadraticBezierPlatform(x0, y0, x1, y1, x2, y2, ds, p)
}

func cubicBezier(x0, y0, x1, y1, x2, y2, x3, y3, ds float64, p []Point) {
	if canUseBezierKernel {
		err := cubicBezeirGPU(x0, y0, x1, y1, x2, y2, x3, y3, ds, p)
		if err == nil {
			return
		}
		canUseBezierKernel = false
	}
	cubicBezierPlatform(x0, y0, x1, y1, x2, y2, x3, y3, ds, p)
}

func quadraticBezierLen(x0, y0, x1, y1, x2, y2 float64) int {
	l := math.Hypot(x1-x0, y1-y0) + math.Hypot(x2-x1, y2-y1)
	n := max(int(l+0.5), 4)
	return n
}

func QuadraticBezier(x0, y0, x1, y1, x2, y2 float64) []Point {
	n := quadraticBezierLen(x0, y0, x1, y1, x2, y2)
	result := make([]Point, n)
	quadraticBezier(x0, y0, x1, y1, x2, y2, float64(n)-1, result)
	return result
}

func cubicBezierLen(x0, y0, x1, y1, x2, y2, x3, y3 float64) int {
	l := math.Hypot(x1-x0, y1-y0) + math.Hypot(x2-x1, y2-y1) + math.Hypot(x3-x2, y3-y2)
	n := max(int(l+0.5), 4)
	return n
}

func CubicBezier(x0, y0, x1, y1, x2, y2, x3, y3 float64) []Point {
	n := cubicBezierLen(x0, y0, x1, y1, x2, y2, x3, y3)
	result := make([]Point, n)
	cubicBezier(x0, y0, x1, y1, x2, y2, x3, y3, float64(n)-1, result)
	return result
}
