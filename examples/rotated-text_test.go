package main

import (
	"testing"

	"github.com/FloatTech/gg"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

func TestRTXT(t *testing.T) {
	const S = 400
	dc := gg.NewContext(S, S)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	font, err := opentype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}
	face, err := opentype.NewFace(font, &opentype.FaceOptions{
		Size: 40,
	})
	if err != nil {
		panic(err)
	}
	dc.SetFontFace(face)
	text := "Hello, world!"
	w, h := dc.MeasureString(text)
	dc.Rotate(gg.Radians(10))
	dc.DrawRectangle(100, 180, w, h)
	dc.Stroke()
	dc.DrawStringAnchored(text, 100, 180, 0.0, 0.0)
	dc.SavePNG(GetFileName() + ".png")
}
