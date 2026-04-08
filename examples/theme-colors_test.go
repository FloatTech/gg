package main

import (
	"image/color"
	"math"
	"testing"

	"github.com/FloatTech/gg"
	"github.com/FloatTech/gg/fio"
)

func TestThemeColorsGopher(t *testing.T) {
	im, err := fio.LoadPNG("gopher.png")
	if err != nil {
		t.Fatal(err)
	}

	expected := []color.RGBA{
		{105, 213, 226, 255}, {249, 240, 227, 255}, {4, 6, 6, 255},
	}

	const (
		k           = 3
		maxAttempts = 30
		tolerance   = 600.0 // 允许的颜色距离平方
	)

	var result []color.RGBA
	found := false
	for range maxAttempts {
		var err error
		result, err = gg.TakeThemeColorsKMeans(im, k)
		if err != nil {
			t.Fatal(err)
		}
		if len(result) != k {
			t.Fatalf("expected %d colors, got %d", k, len(result))
		}
		if allColorsMatch(expected, result, tolerance) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("theme colors did not match expected in %d attempts\nexpected: %v\ngot:      %v", maxAttempts, expected, result)
	}

	// 绘制提取到的主题色方块
	const (
		blockW = 120
		blockH = 120
		pad    = 20
	)
	w := pad + k*(blockW+pad)
	h := pad + blockH + pad
	dc := gg.NewContext(w, h)

	// 白色背景
	dc.SetColor(color.White)
	dc.Clear()

	for i, c := range result {
		x := float64(pad + i*(blockW+pad))
		y := float64(pad)
		dc.SetColor(c)
		dc.DrawRectangle(x, y, blockW, blockH)
		dc.Fill()
	}

	if err := dc.SavePNG(GetFileName() + ".png"); err != nil {
		t.Fatal(err)
	}
}

// allColorsMatch 检查 expected 中每个颜色都能在 result 中找到匹配（无序）
func allColorsMatch(expected, result []color.RGBA, tolerance float64) bool {
	if len(expected) != len(result) {
		return false
	}
	used := make([]bool, len(result))
	for _, e := range expected {
		matched := false
		for j, r := range result {
			if !used[j] && colorDistSq(e, r) <= tolerance {
				used[j] = true
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	return true
}

func colorDistSq(a, b color.RGBA) float64 {
	dr := float64(a.R) - float64(b.R)
	dg := float64(a.G) - float64(b.G)
	db := float64(a.B) - float64(b.B)
	return math.Abs(dr*dr) + math.Abs(dg*dg) + math.Abs(db*db)
}
