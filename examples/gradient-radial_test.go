package main

import (
	"image/color"
	"testing"

	"github.com/FloatTech/gg"
)

func TestGR(t *testing.T) {
	dc := gg.NewContext(400, 200)

	grad := gg.NewRadialGradient(100, 100, 10, 100, 120, 80)
	grad.AddColorStop(0, gg.Green)
	grad.AddColorStop(1, gg.Blue)

	dc.SetFillStyle(grad)
	dc.DrawRectangle(0, 0, 200, 200)
	dc.Fill()

	dc.SetColor(color.White)
	dc.DrawCircle(100, 100, 10)
	dc.Stroke()
	dc.DrawCircle(100, 120, 80)
	dc.Stroke()

	dc.SavePNG(GetFileName() + ".png")
}
