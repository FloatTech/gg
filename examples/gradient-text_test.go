package main

import (
	"testing"

	"github.com/FloatTech/gg"
)

const (
	W = 1024
	H = 512
)

func TestGT(*testing.T) {
	dc := gg.NewContext(W, H)

	// draw text
	dc.SetRGB(0, 0, 0)
	if err := dc.LoadFontFace("/System/Library/Fonts/Supplemental/Impact.ttf", 128); err != nil {
		panic(err)
	}
	dc.DrawStringAnchored("Gradient Text", W/2, H/2, 0.5, 0.5)

	// get the context as an alpha mask
	mask := dc.AsMask()

	// clear the context
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// set a gradient
	g := gg.NewLinearGradient(0, 0, W, H)
	g.AddColorStop(0, gg.Red)
	g.AddColorStop(1, gg.Blue)
	dc.SetFillStyle(g)

	// using the mask, fill the context with the gradient
	if err := dc.SetMask(mask); err != nil {
		panic(err)
	}
	dc.DrawRectangle(0, 0, W, H)
	dc.Fill()

	if err := dc.SavePNG(GetFileName() + ".png"); err != nil {
		panic(err)
	}
}
