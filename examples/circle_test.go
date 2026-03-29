package main

import (
	"testing"

	"github.com/FloatTech/gg"
)

func TestCircle(*testing.T) {
	dc := gg.NewContext(1000, 1000)
	dc.DrawCircle(500, 500, 400)
	dc.SetRGB(0, 0, 0)
	dc.Fill()
	if err := dc.SavePNG(GetFileName() + ".png"); err != nil {
		panic(err)
	}
}
