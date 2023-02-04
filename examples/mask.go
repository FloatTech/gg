package main

import (
	"log"

	"github.com/FloatTech/gg"
)

func main() {
	im, err := gg.LoadImage("examples/james-webb.png")
	if err != nil {
		log.Fatal(err)
	}

	dc := gg.NewContext(512, 512)
	dc.DrawRoundedRectangle(0, 0, 512, 512, 64)
	dc.Clip()
	dc.DrawImage(im, 0, 0)
	dc.SavePNG("out.png")
}
