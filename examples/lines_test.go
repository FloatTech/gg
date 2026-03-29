package main

import (
	"math/rand"
	"testing"

	"github.com/FloatTech/gg"
)

func TestLines(*testing.T) {
	const W = 1024
	const H = 1024
	dc := gg.NewContext(W, H)
	dc.SetRGB(0, 0, 0)
	dc.Clear()
	for i := 0; i < 1000; i++ {
		x1 := rand.Float64() * W
		y1 := rand.Float64() * H
		x2 := rand.Float64() * W
		y2 := rand.Float64() * H
		r := rand.Float64()
		g := rand.Float64()
		b := rand.Float64()
		a := rand.Float64()*0.5 + 0.5
		w := rand.Float64()*4 + 1
		dc.SetRGBA(r, g, b, a)
		dc.SetLineWidth(w)
		dc.DrawLine(x1, y1, x2, y2)
		dc.Stroke()
	}
	if err := dc.SavePNG(GetFileName() + ".png"); err != nil {
		panic(err)
	}
}
