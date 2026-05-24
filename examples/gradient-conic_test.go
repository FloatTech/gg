package main

import (
	"image/color"
	"testing"

	"github.com/FloatTech/gg"
	"github.com/FloatTech/gg/2d/total"
)

func TestGC(*testing.T) {
	dc := gg.NewContext(400, 400)

	grad1 := gg.NewConicGradient(200, 200, 0)
	grad1.AddColorStop(0.0, color.Black)
	grad1.AddColorStop(0.5, color.RGBA{255, 215, 0, 255})
	grad1.AddColorStop(1.0, total.Red)

	grad2 := gg.NewConicGradient(200, 200, 90)
	grad2.AddColorStop(0.00, total.Red)
	grad2.AddColorStop(0.16, total.Yellow)
	grad2.AddColorStop(0.33, total.Green)
	grad2.AddColorStop(0.50, total.Cyan)
	grad2.AddColorStop(0.66, total.Blue)
	grad2.AddColorStop(0.83, total.Magenta)
	grad2.AddColorStop(1.00, total.Red)

	dc.SetStrokeStyle(grad1)
	dc.SetLineWidth(20)
	dc.DrawCircle(200, 200, 180)
	dc.Stroke()

	dc.SetFillStyle(grad2)
	dc.DrawCircle(200, 200, 150)
	dc.Fill()

	if err := dc.SavePNG("gradient-conic.png"); err != nil {
		panic(err)
	}
}
