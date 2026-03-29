package main

import (
	"testing"

	"github.com/FloatTech/gg"
	"github.com/FloatTech/gg/fio"
)

func TestTiling(*testing.T) {
	const NX = 4
	const NY = 3
	im, err := fio.LoadPNG("gopher.png")
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
	if err := dc.SavePNG(GetFileName() + ".png"); err != nil {
		panic(err)
	}
}
