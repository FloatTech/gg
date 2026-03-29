package main

import (
	"testing"

	"github.com/FloatTech/gg"
	"github.com/FloatTech/gg/fio"
)

func TestPF(*testing.T) {
	im, err := fio.LoadPNG("james-webb.png")
	if err != nil {
		panic(err)
	}
	pattern := gg.NewSurfacePattern(im, gg.RepeatBoth)
	dc := gg.NewContext(600, 600)
	dc.MoveTo(20, 20)
	dc.LineTo(590, 20)
	dc.LineTo(590, 590)
	dc.LineTo(20, 590)
	dc.ClosePath()
	dc.SetFillStyle(pattern)
	dc.Fill()
	if err := dc.SavePNG(GetFileName() + ".png"); err != nil {
		panic(err)
	}
}
