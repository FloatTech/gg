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
	dc.Brightness(-50)
	if err := saveImage(dc, "TestBrightness-50"); err != nil {
		t.Fatal(err)
	}
	checkHash(t, dc, "<gg.Context 30e0b039be47de836a9fb4c8b9825ce8>")
	dc.Brightness(70) // 70-50=20
	if err := saveImage(dc, "TestBrightness+20"); err != nil {
		t.Fatal(err)
	}
	checkHash(t, dc, "<gg.Context d8aaf7a3d0aeb11dce3adc72ae9fb464>")
}

func TestContrast(t *testing.T) {
	im, err := fio.LoadPNG("examples/gopher.png")
	if err != nil {
		t.Fatal(err)
	}
	dc := NewContextForImage(im)
	dc.Contrast(-50)
	if err := saveImage(dc, "TestContrast-50"); err != nil {
		t.Fatal(err)
	}
	checkHash(t, dc, "<gg.Context e6f161e094f0d95603aad891c51c0b6b>")
	dc.Contrast(100)
	if err := saveImage(dc, "TestContrast+100"); err != nil {
		t.Fatal(err)
	}
	checkHash(t, dc, "<gg.Context e1608d9268930ac4f80bd5cb82af927a>")
}
