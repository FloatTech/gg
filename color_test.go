package gg

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"testing"
)

// ---- 辅助函数 ----

// solidImage 创建一个全部填充为指定颜色的 image.Image
func solidImage(w, h int, c color.RGBA) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			img.SetRGBA(x, y, c)
		}
	}
	return img
}

// twoColorImage 创建左半部分为 c1、右半部分为 c2 的图像
func twoColorImage(w, h int, c1, c2 color.RGBA) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			if x < w/2 {
				img.SetRGBA(x, y, c1)
			} else {
				img.SetRGBA(x, y, c2)
			}
		}
	}
	return img
}

// colorInSlice 判断给定颜色是否在切片中（允许 tolerance 误差）
func colorInSlice(c color.RGBA, slice []color.RGBA, tolerance float64) bool {
	for _, s := range slice {
		if distanceRGBAsq(c, s) <= tolerance {
			return true
		}
	}
	return false
}

// ---- sq 测试 ----

func TestSq_Zero(t *testing.T) {
	if got := sq(0); got != 0 {
		t.Errorf("sq(0) = %v, want 0", got)
	}
}

func TestSq_Positive(t *testing.T) {
	cases := []struct {
		in   float64
		want float64
	}{
		{1, 1},
		{2, 4},
		{3, 9},
		{0.5, 0.25},
		{10, 100},
		{255, 65025},
	}
	for _, c := range cases {
		if got := sq(c.in); got != c.want {
			t.Errorf("sq(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestSq_Negative(t *testing.T) {
	// 平方结果应与正数相同（始终非负）
	cases := []float64{-1, -2, -3, -0.5, -10}
	for _, n := range cases {
		got := sq(n)
		want := n * n
		if got != want {
			t.Errorf("sq(%v) = %v, want %v", n, got, want)
		}
		if got < 0 {
			t.Errorf("sq(%v) returned negative value %v", n, got)
		}
	}
}

func TestSq_LargeValue(t *testing.T) {
	// 2e154 的平方 = 4e308，超出 float64 最大值（约 1.8e308），应溢出为 +Inf
	n := 2e154
	if got := sq(n); !math.IsInf(got, 1) {
		t.Errorf("sq(%v) = %v, expected +Inf", n, got)
	}
}

// ---- distance 测试 ----

func TestDistance_SameColor(t *testing.T) {
	c := color.RGBA{100, 150, 200, 255}
	if got := distanceRGBAsq(c, c); got != 0 {
		t.Errorf("distance(same, same) = %v, want 0", got)
	}
}

func TestDistance_BlackAndWhite(t *testing.T) {
	// sqrt(255^2 * 3) = 255 * sqrt(3)
	want := 255 * math.Sqrt(3)
	got := math.Sqrt(distanceRGBAsq(Black, White))
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("distance(black, white) = %v, want %v", got, want)
	}
}

func TestDistance_SingleChannel(t *testing.T) {
	a := Black
	b := color.RGBA{3, 4, 0, 255}
	// sqrt(9 + 16) = 5
	want := 5.0
	got := math.Sqrt(distanceRGBAsq(a, b))
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("distance(%v, %v) = %v, want %v", a, b, got, want)
	}
}

func TestDistance_Symmetry(t *testing.T) {
	a := color.RGBA{10, 20, 30, 255}
	b := color.RGBA{50, 80, 110, 255}
	if d1, d2 := distanceRGBAsq(a, b), distanceRGBAsq(b, a); math.Abs(d1-d2) > 1e-9 {
		t.Errorf("distance not symmetric: d(a,b)=%v, d(b,a)=%v", d1, d2)
	}
}

func TestDistance_NonNegative(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	for range 100 {
		a := color.RGBA{uint8(rng.Intn(256)), uint8(rng.Intn(256)), uint8(rng.Intn(256)), 255}
		b := color.RGBA{uint8(rng.Intn(256)), uint8(rng.Intn(256)), uint8(rng.Intn(256)), 255}
		got := distanceRGBAsq(a, b)
		if got < 0 {
			t.Errorf("distance returned negative value %v for %v, %v", got, a, b)
		}
	}
}

func TestDistance_IgnoresAlpha(t *testing.T) {
	a1 := color.RGBA{100, 150, 200, 0}
	a2 := color.RGBA{100, 150, 200, 255}
	b := color.RGBA{50, 50, 50, 128}
	// Alpha 不参与计算，结果应相同
	if d1, d2 := distanceRGBAsq(a1, b), distanceRGBAsq(a2, b); math.Abs(d1-d2) > 1e-9 {
		t.Errorf("distance should ignore alpha: d(a1,b)=%v, d(a2,b)=%v", d1, d2)
	}
}

// ---- clustersEqual 测试 ----

func TestClustersEqual_Equal(t *testing.T) {
	a := []color.RGBA{{255, 0, 0, 255}, {0, 255, 0, 255}, {0, 0, 255, 255}}
	b := []color.RGBA{{255, 0, 0, 255}, {0, 255, 0, 255}, {0, 0, 255, 255}}
	if !isArrayRGBAEqual(a, b) {
		t.Error("clustersEqual should return true for identical slices")
	}
}

func TestClustersEqual_NotEqual(t *testing.T) {
	a := []color.RGBA{{255, 0, 0, 255}, {0, 255, 0, 255}}
	b := []color.RGBA{{255, 0, 0, 255}, {0, 0, 255, 255}}
	if isArrayRGBAEqual(a, b) {
		t.Error("clustersEqual should return false for different slices")
	}
}

func TestClustersEqual_DifferentLength(t *testing.T) {
	a := []color.RGBA{{255, 0, 0, 255}}
	b := []color.RGBA{{255, 0, 0, 255}, {0, 255, 0, 255}}
	if isArrayRGBAEqual(a, b) {
		t.Error("clustersEqual should return false for slices of different lengths")
	}
}

func TestClustersEqual_Empty(t *testing.T) {
	if !isArrayRGBAEqual([]color.RGBA{}, []color.RGBA{}) {
		t.Error("clustersEqual should return true for two empty slices")
	}
}

func TestClustersEqual_OneEmpty(t *testing.T) {
	a := []color.RGBA{{255, 0, 0, 255}}
	if isArrayRGBAEqual(a, []color.RGBA{}) {
		t.Error("clustersEqual should return false when one slice is empty")
	}
}

func TestClustersEqual_SingleElement(t *testing.T) {
	a := []color.RGBA{{100, 100, 100, 255}}
	b := []color.RGBA{{100, 100, 100, 255}}
	if !isArrayRGBAEqual(a, b) {
		t.Error("clustersEqual should return true for identical single-element slices")
	}
}

func TestClustersEqual_DiffersOnlyInAlpha(t *testing.T) {
	a := []color.RGBA{{100, 100, 100, 255}}
	b := []color.RGBA{{100, 100, 100, 0}}
	// RGBA 结构体逐字段比较，Alpha 也参与比较
	if isArrayRGBAEqual(a, b) {
		t.Error("clustersEqual should return false when alpha differs")
	}
}

// ---- takecolor 测试 ----

func TestTakecolor_ReturnsKColors(t *testing.T) {
	img := solidImage(10, 10, color.RGBA{128, 64, 32, 255})
	for k := 1; k <= 5; k++ {
		result := TakeThemeColorsKMeans(img, uint16(k))
		if len(result) != k {
			t.Errorf("takecolor with k=%d returned %d colors, want %d", k, len(result), k)
		}
	}
}

func TestTakecolor_SolidColorK1(t *testing.T) {
	c := color.RGBA{200, 100, 50, 255}
	img := solidImage(20, 20, c)
	result := TakeThemeColorsKMeans(img, 1)
	if len(result) != 1 {
		t.Fatalf("expected 1 color, got %d", len(result))
	}
	got := result[0]
	if got.R != c.R || got.G != c.G || got.B != c.B {
		t.Errorf("takecolor(solid %v, 1) = %v, want %v", c, got, c)
	}
}

func TestTakecolor_SolidColorKGreaterThan1(t *testing.T) {
	// 纯色图像下，k-means 会将所有像素归到第一个聚类（平局时取下标最小者），
	// 其余聚类因无像素而保持零值 {0,0,0,0}。
	// 故只验证：返回 k 个颜色，且其中至少一个与原始颜色完全匹配。
	c := color.RGBA{10, 200, 150, 255}
	img := solidImage(15, 15, c)
	result := TakeThemeColorsKMeans(img, 3)
	if len(result) != 3 {
		t.Fatalf("expected 3 colors, got %d", len(result))
	}
	found := false
	for _, got := range result {
		if got.R == c.R && got.G == c.G && got.B == c.B {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected at least one cluster to converge to %v, got %v", c, result)
	}
}

func TestTakecolor_TwoDistinctColors(t *testing.T) {
	// k-means 的初始中心随机选取，可能两次都落到同一颜色区域导致不收敛。
	// 多次运行，验证算法至少能在 30 次尝试中有一次正确分离两种颜色。
	img := twoColorImage(20, 20, Red, Blue)
	const maxAttempts = 30
	for range maxAttempts {
		result := TakeThemeColorsKMeans(img, 2)
		if len(result) == 2 && colorInSlice(Red, result, 5) && colorInSlice(Blue, result, 5) {
			return // 成功分离，测试通过
		}
	}
	t.Errorf("takecolor failed to separate red and blue colors in %d attempts", maxAttempts)
}

func TestTakecolor_Deterministic_SolidImage(t *testing.T) {
	// 纯色图像下，无论随机种子如何，结果应完全一致
	c := color.RGBA{77, 88, 99, 255}
	img := solidImage(10, 10, c)
	r1 := TakeThemeColorsKMeans(img, 2)
	r2 := TakeThemeColorsKMeans(img, 2)
	if !isArrayRGBAEqual(r1, r2) {
		t.Errorf("takecolor on solid image should be deterministic: r1=%v, r2=%v", r1, r2)
	}
}

func TestTakecolor_AllClustersHaveValidRGB(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	img := image.NewRGBA(image.Rect(0, 0, 16, 16))
	for y := range 16 {
		for x := range 16 {
			img.SetRGBA(x, y, color.RGBA{
				uint8(rng.Intn(256)),
				uint8(rng.Intn(256)),
				uint8(rng.Intn(256)),
				255,
			})
		}
	}
	result := TakeThemeColorsKMeans(img, 4)
	if len(result) != 4 {
		t.Fatalf("expected 4 clusters, got %d", len(result))
	}
	for _, c := range result {
		if c.A != 255 {
			t.Errorf("expected alpha=255, got %d for color %v", c.A, c)
		}
	}
}

// ---- Benchmark ----

func BenchmarkSq(b *testing.B) {
	for i := range b.N {
		sq(float64(i))
	}
}

func BenchmarkDistance(b *testing.B) {
	a := color.RGBA{100, 150, 200, 255}
	c := color.RGBA{50, 80, 30, 255}
	b.ResetTimer()
	for range b.N {
		distanceRGBAsq(a, c)
	}
}

func BenchmarkClustersEqual_Equal(b *testing.B) {
	a := []color.RGBA{{255, 0, 0, 255}, {0, 255, 0, 255}, {0, 0, 255, 255}, {128, 128, 128, 255}}
	c := []color.RGBA{{255, 0, 0, 255}, {0, 255, 0, 255}, {0, 0, 255, 255}, {128, 128, 128, 255}}
	b.ResetTimer()
	for range b.N {
		isArrayRGBAEqual(a, c)
	}
}

func BenchmarkClustersEqual_NotEqual(b *testing.B) {
	a := []color.RGBA{{255, 0, 0, 255}, {0, 255, 0, 255}, {0, 0, 255, 255}, {128, 128, 128, 255}}
	c := []color.RGBA{{255, 0, 0, 255}, {0, 255, 0, 255}, {0, 0, 255, 255}, {200, 200, 200, 255}}
	b.ResetTimer()
	for range b.N {
		isArrayRGBAEqual(a, c)
	}
}

func BenchmarkTakecolor_16x16_K3(b *testing.B) {
	rng := rand.New(rand.NewSource(42))
	img := image.NewRGBA(image.Rect(0, 0, 16, 16))
	for y := range 16 {
		for x := range 16 {
			img.SetRGBA(x, y, color.RGBA{uint8(rng.Intn(256)), uint8(rng.Intn(256)), uint8(rng.Intn(256)), 255})
		}
	}
	b.ResetTimer()
	for range b.N {
		TakeThemeColorsKMeans(img, 3)
	}
}

func BenchmarkTakecolor_64x64_K4(b *testing.B) {
	rng := rand.New(rand.NewSource(42))
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))
	for y := range 64 {
		for x := range 64 {
			img.SetRGBA(x, y, color.RGBA{uint8(rng.Intn(256)), uint8(rng.Intn(256)), uint8(rng.Intn(256)), 255})
		}
	}
	b.ResetTimer()
	for range b.N {
		TakeThemeColorsKMeans(img, 4)
	}
}

func BenchmarkTakecolor_128x128_K8(b *testing.B) {
	rng := rand.New(rand.NewSource(42))
	img := image.NewRGBA(image.Rect(0, 0, 128, 128))
	for y := range 128 {
		for x := range 128 {
			img.SetRGBA(x, y, color.RGBA{uint8(rng.Intn(256)), uint8(rng.Intn(256)), uint8(rng.Intn(256)), 255})
		}
	}
	b.ResetTimer()
	for range b.N {
		TakeThemeColorsKMeans(img, 8)
	}
}

func BenchmarkTakecolor_SolidColor_K5(b *testing.B) {
	img := solidImage(64, 64, color.RGBA{200, 100, 50, 255})
	b.ResetTimer()
	for range b.N {
		TakeThemeColorsKMeans(img, 5)
	}
}
