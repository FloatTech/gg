package main

import (
	"log"
	"testing"

	"github.com/FloatTech/gg"
	"github.com/FloatTech/gg/fio"
)

func TestMask(t *testing.T) {
	im, err := fio.LoadImage("james-webb.png")
	if err != nil {
		log.Fatal(err)
	}

	dc := gg.NewContext(512, 512)
	dc.DrawRoundedRectangle(0, 0, 512, 512, 64)
	dc.Clip()
	dc.DrawImage(im, 0, 0)
	dc.SavePNG(GetFileName() + ".png")
}
