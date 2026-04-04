package gg

import (
	"testing"

	"github.com/FloatTech/gg/fio"
)

func TestBrightness(t *testing.T) {
	im, err := fio.LoadPNG("examples/gopher.png")
	if err != nil {
		t.Fatal(err)
	}
	dc := NewContextForImage(im)
	dc.AdjustBrightness(-50)
	if err := saveImage(dc, "TestBrightness-50"); err != nil {
		t.Fatal(err)
	}
	checkHash(t, dc, "<gg.Context e495391cc1349c0db98431a3f4ef21d8>")
	dc.AdjustBrightness(70) // 70-50=20
	if err := saveImage(dc, "TestBrightness+20"); err != nil {
		t.Fatal(err)
	}
	checkHash(t, dc, "<gg.Context 2a1b5df094e1ffb6b8de120d922400bf>")
}

func TestContrast(t *testing.T) {
	im, err := fio.LoadPNG("examples/gopher.png")
	if err != nil {
		t.Fatal(err)
	}
	dc := NewContextForImage(im)
	dc.AdjustContrast(-50)
	if err := saveImage(dc, "TestContrast-50"); err != nil {
		t.Fatal(err)
	}
	checkHash(t, dc, "<gg.Context 71bfdcec14fa18485e85c80b7a7c3fa3>")
	dc.AdjustContrast(100)
	if err := saveImage(dc, "TestContrast+100"); err != nil {
		t.Fatal(err)
	}
	checkHash(t, dc, "<gg.Context 4854d92c7e15539580349c6876f39070>")
}
