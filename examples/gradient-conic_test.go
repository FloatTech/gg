package main

import (
	"image/color"
	"testing"

	"github.com/FloatTech/gg"
)

func TestGC(t *testing.T) {
	dc := gg.NewContext(400, 400)

	grad1 := gg.NewConicGradient(200, 200, 0)
	grad1.AddColorStop(0.0, color.Black)
	grad1.AddColorStop(0.5, color.RGBA{255, 215, 0, 255})
	grad1.AddColorStop(1.0, gg.Red)

	grad2 := gg.NewConicGradient(200, 200, 90)
	grad2.AddColorStop(0.00, gg.Red)
	grad2.AddColorStop(0.16, gg.Yellow)
	grad2.AddColorStop(0.33, gg.Green)
	grad2.AddColorStop(0.50, gg.Cyan)
	grad2.AddColorStop(0.66, gg.Blue)
	grad2.AddColorStop(0.83, gg.Magenta)
	grad2.AddColorStop(1.00, gg.Red)

	dc.SetStrokeStyle(grad1)
	dc.SetLineWidth(20)
	dc.DrawCircle(200, 200, 180)
	dc.Stroke()

	dc.SetFillStyle(grad2)
	dc.DrawCircle(200, 200, 150)
	dc.Fill()

	dc.SavePNG("gradient-conic.png")
}
