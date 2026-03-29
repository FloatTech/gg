package main

import (
	"testing"

	"github.com/FloatTech/gg"
)

func TestText(*testing.T) {
	const S = 1024
	dc := gg.NewContext(S, S)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	// 字体路径
	if err := dc.LoadFontFace("/System/Library/Fonts/Supplemental/Arial.ttf", 96); err != nil {
		panic(err)
	}
	dc.DrawStringAnchored("Hello, world!", S/2, S/2, 0.5, 0.5)
	if err := dc.SavePNG(GetFileName() + ".png"); err != nil {
		panic(err)
	}
}
