//go:build !amd64

package gg

func quadraticBezierPlatform(x0, y0, x1, y1, x2, y2, ds float64, p []Point) {
	quadraticBezierPure(x0, y0, x1, y1, x2, y2, ds, p)
}

func cubicBezierPlatform(x0, y0, x1, y1, x2, y2, x3, y3, ds float64, p []Point) {
	cubicBezierPure(x0, y0, x1, y1, x2, y2, x3, y3, ds, p)
}
