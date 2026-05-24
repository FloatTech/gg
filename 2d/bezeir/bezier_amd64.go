package bezeir

import "github.com/FloatTech/gg/2d/point"

func quadraticBezierASM(x0, y0, x1, y1, x2, y2, ds float64, p []point.Point)

func cubicBezierASM(x0, y0, x1, y1, x2, y2, x3, y3, ds float64, p []point.Point)

func quadraticBezierPlatform(x0, y0, x1, y1, x2, y2, ds float64, p []point.Point) {
	quadraticBezierASM(x0, y0, x1, y1, x2, y2, ds, p)
}

func cubicBezierPlatform(x0, y0, x1, y1, x2, y2, x3, y3, ds float64, p []point.Point) {
	cubicBezierASM(x0, y0, x1, y1, x2, y2, x3, y3, ds, p)
}
