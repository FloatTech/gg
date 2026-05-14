package main

import (
	"testing"

	"github.com/FloatTech/gg"
	"github.com/FloatTech/gg/fio"
)

func TestResize(*testing.T) {
	img, err := fio.LoadPNG("gopher.png")
	if err != nil {
		panic(err)
	}
	small := gg.Resize(img, img.Bounds().Dx()/4, img.Bounds().Dy()/4, gg.ResampleFilterLinear)
	if err := fio.SavePNG(GetFileName()+".png", small); err != nil {
		panic(err)
	}
}
