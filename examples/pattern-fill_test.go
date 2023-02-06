package main

import (
	"testing"

	"github.com/FloatTech/gg"
)

func TestPF(t *testing.T) {
	im, err := gg.LoadPNG("james-webb.png")
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
	dc.SavePNG(GetFileName() + ".png")
}
