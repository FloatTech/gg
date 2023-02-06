package main

import (
	"testing"

	"github.com/FloatTech/gg"
)

func TestClip(t *testing.T) {
	dc := gg.NewContext(1000, 1000)
	dc.DrawCircle(350, 500, 300)
	dc.Clip()
	dc.DrawCircle(650, 500, 300)
	dc.Clip()
	dc.DrawRectangle(0, 0, 1000, 1000)
	dc.SetRGB(0, 0, 0)
	dc.Fill()
	dc.SavePNG(GetFileName() + ".png")
}
