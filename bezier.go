package gg

import "math"

func QuadraticBezier(x0, y0, x1, y1, x2, y2 float64) []Point {
	l := math.Hypot(x1-x0, y1-y0) + math.Hypot(x2-x1, y2-y1)
	n := int(l + 0.5)
	if n < 4 {
		n = 4
	}
	result := make([]Point, n)
	quadratic(x0, y0, x1, y1, x2, y2, float64(n)-1, result)
	return result
}

func CubicBezier(x0, y0, x1, y1, x2, y2, x3, y3 float64) []Point {
	l := math.Hypot(x1-x0, y1-y0) + math.Hypot(x2-x1, y2-y1) + math.Hypot(x3-x2, y3-y2)
	n := int(l + 0.5)
	if n < 4 {
		n = 4
	}
	result := make([]Point, n)
	cubic(x0, y0, x1, y1, x2, y2, x3, y3, float64(n)-1, result)
	return result
}
