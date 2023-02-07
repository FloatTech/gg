package gg

import "testing"

func TestBrightness(t *testing.T) {
	im, err := LoadPNG("examples/gopher.png")
	if err != nil {
		t.Fatal(err)
	}
	dc := NewContextForImage(im)
	dc.Brightness(-50)
	saveImage(dc, "TestBrightness")
	checkHash(t, dc, "<gg.Context 30e0b039be47de836a9fb4c8b9825ce8>")
}

func TestContrast(t *testing.T) {
	im, err := LoadPNG("examples/gopher.png")
	if err != nil {
		t.Fatal(err)
	}
	dc := NewContextForImage(im)
	dc.Contrast(-50)
	saveImage(dc, "TestContrast")
	checkHash(t, dc, "<gg.Context e6f161e094f0d95603aad891c51c0b6b>")
}
