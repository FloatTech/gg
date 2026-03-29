package main

import (
	"testing"

	"github.com/FloatTech/gg"
)

func TestUnicode(*testing.T) {
	const S = 4096 * 2
	const T = 16 * 2
	const F = 28
	dc := gg.NewContext(S, S)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	if err := dc.LoadFontFace(fontPath("Impact"), F); err != nil {
		panic(err)
	}
	for r := range 256 {
		for c := range 256 {
			i := r*256 + c
			x := float64(c*T) + T/2
			y := float64(r*T) + T/2
			dc.DrawStringAnchored(string(rune(i)), x, y, 0.5, 0.5)
		}
	}
	if err := dc.SavePNG(GetFileName() + ".png"); err != nil {
		panic(err)
	}
}
