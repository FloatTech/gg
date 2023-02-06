package main

import (
	"math"
	"testing"

	"github.com/FloatTech/gg"
)

func Polygon4(n int, x, y, r float64) []gg.Point {
	result := make([]gg.Point, n)
	for i := 0; i < n; i++ {
		a := float64(i)*2*math.Pi/float64(n) - math.Pi/2
		result[i] = gg.Point{X: x + r*math.Cos(a), Y: y + r*math.Sin(a)}
	}
	return result
}

func TestStar(t *testing.T) {
	n := 5
	points := Polygon4(n, 512, 512, 400)
	dc := gg.NewContext(1024, 1024)
	dc.SetHexColor("fff")
	dc.Clear()
	for i := 0; i < n+1; i++ {
		index := (i * 2) % n
		p := points[index]
		dc.LineTo(p.X, p.Y)
	}
	dc.SetRGBA(0, 0.5, 0, 1)
	dc.SetFillRule(gg.FillRuleEvenOdd)
	dc.FillPreserve()
	dc.SetRGBA(0, 1, 0, 0.5)
	dc.SetLineWidth(16)
	dc.Stroke()
	dc.SavePNG(GetFileName() + ".png")
}
