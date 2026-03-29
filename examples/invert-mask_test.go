package main

import (
	"testing"

	"github.com/FloatTech/gg"
)

func TestIM(*testing.T) {
	dc := gg.NewContext(1024, 1024)
	dc.DrawCircle(512, 512, 384)
	dc.Clip()
	dc.InvertMask()
	dc.DrawRectangle(0, 0, 1024, 1024)
	dc.SetRGB(0, 0, 0)
	dc.Fill()
	if err := dc.SavePNG(GetFileName() + ".png"); err != nil {
		panic(err)
	}
}
