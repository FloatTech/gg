package main

import (
	"image/color"
	"testing"

	"github.com/FloatTech/gg"
)

func TestGL(t *testing.T) {
	dc := gg.NewContext(500, 400)

	grad := gg.NewLinearGradient(20, 320, 400, 20)
	grad.AddColorStop(0, gg.Green)
	grad.AddColorStop(1, gg.Blue)
	grad.AddColorStop(0.5, gg.Red)

	dc.SetColor(color.White)
	dc.DrawRectangle(20, 20, 400-20, 300)
	dc.Stroke()

	dc.SetStrokeStyle(grad)
	dc.SetLineWidth(4)
	dc.MoveTo(10, 10)
	dc.LineTo(410, 10)
	dc.LineTo(410, 100)
	dc.LineTo(10, 100)
	dc.ClosePath()
	dc.Stroke()

	dc.SetFillStyle(grad)
	dc.MoveTo(10, 120)
	dc.LineTo(410, 120)
	dc.LineTo(410, 300)
	dc.LineTo(10, 300)
	dc.ClosePath()
	dc.Fill()

	dc.SavePNG(GetFileName() + ".png")
}
