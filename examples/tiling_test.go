package main

import (
	"testing"

	"github.com/FloatTech/gg"
)

func TestTiling(t *testing.T) {
	const NX = 4
	const NY = 3
	im, err := gg.LoadPNG("gopher.png")
	if err != nil {
		panic(err)
	}
	w := im.Bounds().Size().X
	h := im.Bounds().Size().Y
	dc := gg.NewContext(w*NX, h*NY)
	for y := 0; y < NY; y++ {
		for x := 0; x < NX; x++ {
			dc.DrawImage(im, x*w, y*h)
		}
	}
	dc.SavePNG(GetFileName() + ".png")
}
