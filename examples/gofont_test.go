package main

import (
	"log"
	"testing"

	"github.com/FloatTech/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

func TestGoFont(t *testing.T) {
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}

	face := truetype.NewFace(font, &truetype.Options{Size: 48})

	dc := gg.NewContext(1024, 1024)
	dc.SetFontFace(face)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	dc.DrawStringAnchored("Hello, world!", 512, 512, 0.5, 0.5)
	dc.SavePNG(GetFileName() + ".png")
}
