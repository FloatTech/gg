package main

import (
	"math"
	"math/rand"
	"testing"

	"github.com/FloatTech/gg"
)

func Polygon(n int) []gg.Point {
	result := make([]gg.Point, n)
	for i := range n {
		a := float64(i)*2*math.Pi/float64(n) - math.Pi/2
		result[i] = gg.Point{X: math.Cos(a), Y: math.Sin(a)}
	}
	return result
}

func TestStars(*testing.T) {
	const W = 1200
	const H = 120
	const S = 100
	dc := gg.NewContext(W, H)
	dc.SetHexColor("#FFFFFF")
	dc.Clear()
	n := 5
	points := Polygon(n)
	for x := S / 2; x < W; x += S {
		dc.Push()
		s := rand.Float64()*S/4 + S/4
		dc.Translate(float64(x), H/2)
		dc.Rotate(rand.Float64() * 2 * math.Pi)
		dc.Scale(s, s)
		for i := range n + 1 {
			index := (i * 2) % n
			p := points[index]
			dc.LineTo(p.X, p.Y)
		}
		dc.SetLineWidth(10)
		dc.SetHexColor("#FFCC00")
		dc.StrokePreserve()
		dc.SetHexColor("#FFE43A")
		dc.Fill()
		dc.Pop()
	}
	if err := dc.SavePNG(GetFileName() + ".png"); err != nil {
		panic(err)
	}
}
