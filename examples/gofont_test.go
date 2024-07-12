package main

import (
	"log"
	"testing"

	"github.com/FloatTech/gg"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

func TestGoFont(t *testing.T) {
	font, err := opentype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}

	face, err := opentype.NewFace(font, &opentype.FaceOptions{Size: 48})
	if err != nil {
		log.Fatal(err)
	}

	dc := gg.NewContext(1024, 1024)
	dc.SetFontFace(face)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	dc.DrawStringAnchored("Hello, world!", 512, 512, 0.5, 0.5)
	dc.SavePNG(GetFileName() + ".png")
}
