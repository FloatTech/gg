package gg

import (
	"math"

	"golang.org/x/image/math/fixed"
)

// Point represents a 2D point with X and Y coordinates.
//
// Point 表示一个具有 X 和 Y 坐标的二维点。
type Point struct {
	X, Y float64
}

// Fixed converts the point to a fixed-point representation.
//
// Fixed 将点转换为定点数表示。
func (a Point) Fixed() fixed.Point26_6 {
	return fixp(a.X, a.Y)
}

// Distance returns the Euclidean distance between two points.
//
// Distance 返回两点之间的欧几里得距离。
func (a Point) Distance(b Point) float64 {
	return math.Hypot(a.X-b.X, a.Y-b.Y)
}

// Interpolate linearly interpolates between two points by parameter t.
//
// Interpolate 按参数 t 在两点之间进行线性插值。
func (a Point) Interpolate(b Point, t float64) Point {
	x := a.X + (b.X-a.X)*t
	y := a.Y + (b.Y-a.Y)*t
	return Point{x, y}
}
