package main

import (
	"image/color"
	"testing"

	"github.com/FloatTech/gg"
	"github.com/FloatTech/gg/2d/total"
)

func TestGR(*testing.T) {
	dc := gg.NewContext(400, 200)

	grad := gg.NewRadialGradient(100, 100, 10, 100, 120, 80)
	grad.AddColorStop(0, total.Green)
	grad.AddColorStop(1, total.Blue)

	dc.SetFillStyle(grad)
	dc.DrawRectangle(0, 0, 200, 200)
	dc.Fill()

	dc.SetColor(color.White)
	dc.DrawCircle(100, 100, 10)
	dc.Stroke()
	dc.DrawCircle(100, 120, 80)
	dc.Stroke()

	if err := dc.SavePNG(GetFileName() + ".png"); err != nil {
		panic(err)
	}
}
