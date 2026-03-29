package main

import (
	"math"
	"testing"

	"github.com/FloatTech/gg"
	"github.com/FloatTech/gg/fio"
)

func TestConcat(*testing.T) {
	im1, err := fio.LoadPNG("james-webb.png")
	if err != nil {
		panic(err)
	}

	im2, err := fio.LoadPNG("gopher.png")
	if err != nil {
		panic(err)
	}

	s1 := im1.Bounds().Size()
	s2 := im2.Bounds().Size()

	width := int(math.Max(float64(s1.X), float64(s2.X)))
	height := s1.Y + s2.Y

	dc := gg.NewContext(width, height)
	dc.DrawImage(im1, 0, 0)
	dc.DrawImage(im2, 0, s1.Y)
	if err := dc.SavePNG(GetFileName() + ".png"); err != nil {
		panic(err)
	}
}
