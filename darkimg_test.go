package gg

import (
	"image"
	"image/color"
	"testing"
)

// darkimgTestImage 创建指定颜色的 RGBA 图像
func darkimgTestImage(w, h int, c color.RGBA) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			img.SetRGBA(x, y, c)
		}
	}
	return img
}

// darkimgGradientImage 创建从黑到白的水平渐变图像
func darkimgGradientImage(w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			v := uint8(x * 255 / w)
			img.SetRGBA(x, y, color.RGBA{v, v, v, 255})
		}
	}
	return img
}

// darkimgPartialBrightImage 创建部分亮像素的图像
// brightPct 指定亮像素的百分比 (0-100)
func darkimgPartialBrightImage(w, h int, brightPct int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	total := w * h
	brightCount := total * brightPct / 100
	idx := 0
	for y := range h {
		for x := range w {
			if idx < brightCount {
				img.SetRGBA(x, y, color.RGBA{200, 200, 200, 255})
			} else {
				img.SetRGBA(x, y, color.RGBA{0, 0, 0, 255})
			}
			idx++
		}
	}
	return img
}

// ---- cpuIsDarkimg 单元测试 ----

func TestCpuIsDarkimg_AllBlack(t *testing.T) {
	img := darkimgTestImage(100, 100, color.RGBA{0, 0, 0, 255})
	if !cpuIsDarkimg(img, 1.0) {
		t.Error("全黑图像应判定为 dark")
	}
}

func TestCpuIsDarkimg_AllWhite(t *testing.T) {
	img := darkimgTestImage(100, 100, color.RGBA{255, 255, 255, 255})
	if cpuIsDarkimg(img, 1.0) {
		t.Error("全白图像不应判定为 dark")
	}
}

func TestCpuIsDarkimg_AllTransparentBlack(t *testing.T) {
	// 全透明黑色 — 像素值仍为 0, 应判定为 dark
	img := darkimgTestImage(100, 100, color.RGBA{0, 0, 0, 0})
	if !cpuIsDarkimg(img, 1.0) {
		t.Error("全透明黑色图像应判定为 dark（基于像素值而非可见性）")
	}
}

func TestCpuIsDarkimg_NearBlackBelowThreshold(t *testing.T) {
	// 亮度 = (299*14+587*14+114*14)/1000 = 14, 低于 visibleThreshold=15
	img := darkimgTestImage(100, 100, color.RGBA{14, 14, 14, 255})
	if !cpuIsDarkimg(img, 1.0) {
		t.Error("亮度恰好低于阈值的图像应判定为 dark")
	}
}

func TestCpuIsDarkimg_NearBlackAboveThreshold(t *testing.T) {
	// 亮度 = (299*16+587*16+114*16)/1000 = 16, 高于 visibleThreshold=15
	img := darkimgTestImage(100, 100, color.RGBA{16, 16, 16, 255})
	if cpuIsDarkimg(img, 1.0) {
		t.Error("全部像素亮度高于阈值的图像不应判定为 dark")
	}
}

func TestCpuIsDarkimg_ExactlyAtThreshold(t *testing.T) {
	// 亮度 = 15, 等于 visibleThreshold, 不满足 > 条件
	img := darkimgTestImage(100, 100, color.RGBA{15, 15, 15, 255})
	if !cpuIsDarkimg(img, 1.0) {
		t.Error("亮度恰好等于阈值的图像应判定为 dark（条件为严格大于）")
	}
}

func TestCpuIsDarkimg_Exactly5PercentBright(t *testing.T) {
	// 恰好 5% 亮像素，不满足 < 5 条件，不应判定为 dark
	img := darkimgPartialBrightImage(100, 100, 5)
	if cpuIsDarkimg(img, 1.0) {
		t.Error("恰好 5%% 亮像素不应判定为 dark（条件为严格小于 5%%）")
	}
}

func TestCpuIsDarkimg_4PercentBright(t *testing.T) {
	// 4% 亮像素，应判定为 dark
	img := darkimgPartialBrightImage(100, 100, 4)
	if !cpuIsDarkimg(img, 1.0) {
		t.Error("仅 4%% 亮像素的图像应判定为 dark")
	}
}

func TestCpuIsDarkimg_6PercentBright(t *testing.T) {
	// 6% 亮像素，不应判定为 dark
	img := darkimgPartialBrightImage(100, 100, 6)
	if cpuIsDarkimg(img, 1.0) {
		t.Error("6%% 亮像素的图像不应判定为 dark")
	}
}

func TestCpuIsDarkimg_GradientImage(t *testing.T) {
	// 渐变图像约有一半亮像素，不应判定为 dark
	img := darkimgGradientImage(256, 100)
	if cpuIsDarkimg(img, 1.0) {
		t.Error("渐变图像不应判定为 dark")
	}
}

func TestCpuIsDarkimg_ZeroSizeImage(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 0, 0))
	// 零大小图像经过 Resize 后仍为 0 像素，应返回 true
	if !cpuIsDarkimg(img, 1.0) {
		t.Error("零大小图像应判定为 dark")
	}
}

func TestCpuIsDarkimg_SinglePixelBright(t *testing.T) {
	img := darkimgTestImage(1, 1, color.RGBA{255, 255, 255, 255})
	if cpuIsDarkimg(img, 1.0) {
		t.Error("单个白色像素不应判定为 dark")
	}
}

func TestCpuIsDarkimg_SinglePixelDark(t *testing.T) {
	img := darkimgTestImage(1, 1, color.RGBA{0, 0, 0, 255})
	if !cpuIsDarkimg(img, 1.0) {
		t.Error("单个黑色像素应判定为 dark")
	}
}

func TestCpuIsDarkimg_Scale(t *testing.T) {
	// 大图 + 小 scale, 验证缩放后行为一致
	img := darkimgTestImage(1000, 1000, color.RGBA{255, 255, 255, 255})
	if cpuIsDarkimg(img, 0.1) {
		t.Error("全白大图缩放到 10%% 仍不应判定为 dark")
	}
}

func TestCpuIsDarkimg_RedChannel(t *testing.T) {
	// 纯红色：亮度 = 299*255/1000 ≈ 76，高于阈值
	img := darkimgTestImage(100, 100, color.RGBA{255, 0, 0, 255})
	if cpuIsDarkimg(img, 1.0) {
		t.Error("纯红色图像不应判定为 dark")
	}
}

func TestCpuIsDarkimg_GreenChannel(t *testing.T) {
	// 纯绿色：亮度 = 587*255/1000 ≈ 149，高于阈值
	img := darkimgTestImage(100, 100, color.RGBA{0, 255, 0, 255})
	if cpuIsDarkimg(img, 1.0) {
		t.Error("纯绿色图像不应判定为 dark")
	}
}

func TestCpuIsDarkimg_BlueChannel(t *testing.T) {
	// 纯蓝色：亮度 = 114*255/1000 ≈ 29，高于阈值
	img := darkimgTestImage(100, 100, color.RGBA{0, 0, 255, 255})
	if cpuIsDarkimg(img, 1.0) {
		t.Error("纯蓝色图像不应判定为 dark")
	}
}

func TestCpuIsDarkimg_VeryDarkBlue(t *testing.T) {
	// 暗蓝色: 亮度 = 114*10/1000 = 1, 远低于阈值
	img := darkimgTestImage(100, 100, color.RGBA{0, 0, 10, 255})
	if !cpuIsDarkimg(img, 1.0) {
		t.Error("极暗蓝色图像应判定为 dark")
	}
}

func TestCpuIsDarkimg_LargeImage(t *testing.T) {
	// 大图性能/正确性测试
	img := darkimgTestImage(500, 500, color.RGBA{0, 0, 0, 255})
	if !cpuIsDarkimg(img, 0.5) {
		t.Error("全黑大图缩放后应判定为 dark")
	}
}

func TestCpuIsDarkimg_NonSquareImage(t *testing.T) {
	img := darkimgTestImage(200, 50, color.RGBA{128, 128, 128, 255})
	if cpuIsDarkimg(img, 1.0) {
		t.Error("中等亮度非方形图像不应判定为 dark")
	}
}

// ---- IsDarkimg 整体测试（GPU/CPU 一致性）----

func TestIsDarkimg_AllBlack(t *testing.T) {
	img := darkimgTestImage(100, 100, color.RGBA{0, 0, 0, 255})
	if !canUseDarkimgKernel {
		t.Skip("no gpu available")
		defer func() {
			if !canUseDarkimgKernel {
				panic("unexpected")
			}
		}()
		return
	}
	result := IsDarkimg(img, 1.0)
	expected := cpuIsDarkimg(img, 1.0)
	if result != expected {
		t.Errorf("IsDarkimg 全黑图像结果 %v 与 cpuIsDarkimg %v 不一致", result, expected)
	}
	if !result {
		t.Error("全黑图像应判定为 dark")
	}
}

func TestIsDarkimg_AllWhite(t *testing.T) {
	img := darkimgTestImage(100, 100, color.RGBA{255, 255, 255, 255})
	if !canUseDarkimgKernel {
		t.Skip("no gpu available")
		defer func() {
			if !canUseDarkimgKernel {
				panic("unexpected")
			}
		}()
		return
	}
	result := IsDarkimg(img, 1.0)
	expected := cpuIsDarkimg(img, 1.0)
	if result != expected {
		t.Errorf("IsDarkimg 全白图像结果 %v 与 cpuIsDarkimg %v 不一致", result, expected)
	}
	if result {
		t.Error("全白图像不应判定为 dark")
	}
}

func TestIsDarkimg_BelowThreshold(t *testing.T) {
	img := darkimgTestImage(100, 100, color.RGBA{14, 14, 14, 255})
	if !canUseDarkimgKernel {
		t.Skip("no gpu available")
		defer func() {
			if !canUseDarkimgKernel {
				panic("unexpected")
			}
		}()
		return
	}
	result := IsDarkimg(img, 1.0)
	expected := cpuIsDarkimg(img, 1.0)
	if result != expected {
		t.Errorf("IsDarkimg 低于阈值图像结果 %v 与 cpuIsDarkimg %v 不一致", result, expected)
	}
}

func TestIsDarkimg_AboveThreshold(t *testing.T) {
	img := darkimgTestImage(100, 100, color.RGBA{17, 17, 17, 255})
	if !canUseDarkimgKernel {
		t.Skip("no gpu available")
		defer func() {
			if !canUseDarkimgKernel {
				panic("unexpected")
			}
		}()
		return
	}
	result := IsDarkimg(img, 1.0)
	expected := cpuIsDarkimg(img, 1.0)
	if result != expected {
		t.Errorf("IsDarkimg 高于阈值图像结果 %v 与 cpuIsDarkimg %v 不一致", result, expected)
	}
}

func TestIsDarkimg_Gradient(t *testing.T) {
	img := darkimgGradientImage(256, 100)
	if !canUseDarkimgKernel {
		t.Skip("no gpu available")
		defer func() {
			if !canUseDarkimgKernel {
				panic("unexpected")
			}
		}()
		return
	}
	result := IsDarkimg(img, 1.0)
	expected := cpuIsDarkimg(img, 1.0)
	if result != expected {
		t.Errorf("IsDarkimg 渐变图像结果 %v 与 cpuIsDarkimg %v 不一致", result, expected)
	}
}

func TestIsDarkimg_3PercentBright(t *testing.T) {
	img := darkimgPartialBrightImage(100, 100, 3)
	if !canUseDarkimgKernel {
		t.Skip("no gpu available")
		defer func() {
			if !canUseDarkimgKernel {
				panic("unexpected")
			}
		}()
		return
	}
	result := IsDarkimg(img, 1.0)
	expected := cpuIsDarkimg(img, 1.0)
	if result != expected {
		t.Errorf("IsDarkimg 4%% 亮像素结果 %v 与 cpuIsDarkimg %v 不一致", result, expected)
	}
}

func TestIsDarkimg_6PercentBright(t *testing.T) {
	img := darkimgPartialBrightImage(100, 100, 6)
	if !canUseDarkimgKernel {
		t.Skip("no gpu available")
		defer func() {
			if !canUseDarkimgKernel {
				panic("unexpected")
			}
		}()
		return
	}
	result := IsDarkimg(img, 1.0)
	expected := cpuIsDarkimg(img, 1.0)
	if result != expected {
		t.Errorf("IsDarkimg 6%% 亮像素结果 %v 与 cpuIsDarkimg %v 不一致", result, expected)
	}
}

func TestIsDarkimg_ScaleDown(t *testing.T) {
	img := darkimgTestImage(500, 500, color.RGBA{100, 100, 100, 255})
	if !canUseDarkimgKernel {
		t.Skip("no gpu available")
		defer func() {
			if !canUseDarkimgKernel {
				panic("unexpected")
			}
		}()
		return
	}
	result := IsDarkimg(img, 0.2)
	expected := cpuIsDarkimg(img, 0.2)
	if result != expected {
		t.Errorf("IsDarkimg 缩放结果 %v 与 cpuIsDarkimg %v 不一致", result, expected)
	}
}

func TestIsDarkimg_Red(t *testing.T) {
	img := darkimgTestImage(100, 100, color.RGBA{255, 0, 0, 255})
	if !canUseDarkimgKernel {
		t.Skip("no gpu available")
		defer func() {
			if !canUseDarkimgKernel {
				panic("unexpected")
			}
		}()
		return
	}
	result := IsDarkimg(img, 1.0)
	expected := cpuIsDarkimg(img, 1.0)
	if result != expected {
		t.Errorf("IsDarkimg 红色图像结果 %v 与 cpuIsDarkimg %v 不一致", result, expected)
	}
}

func TestIsDarkimg_Green(t *testing.T) {
	img := darkimgTestImage(100, 100, color.RGBA{0, 255, 0, 255})
	if !canUseDarkimgKernel {
		t.Skip("no gpu available")
		defer func() {
			if !canUseDarkimgKernel {
				panic("unexpected")
			}
		}()
		return
	}
	result := IsDarkimg(img, 1.0)
	expected := cpuIsDarkimg(img, 1.0)
	if result != expected {
		t.Errorf("IsDarkimg 绿色图像结果 %v 与 cpuIsDarkimg %v 不一致", result, expected)
	}
}

func TestIsDarkimg_Blue(t *testing.T) {
	img := darkimgTestImage(100, 100, color.RGBA{0, 0, 255, 255})
	if !canUseDarkimgKernel {
		t.Skip("no gpu available")
		defer func() {
			if !canUseDarkimgKernel {
				panic("unexpected")
			}
		}()
		return
	}
	result := IsDarkimg(img, 1.0)
	expected := cpuIsDarkimg(img, 1.0)
	if result != expected {
		t.Errorf("IsDarkimg 蓝色图像结果 %v 与 cpuIsDarkimg %v 不一致", result, expected)
	}
}

func TestIsDarkimg_VeryDarkBlue(t *testing.T) {
	img := darkimgTestImage(100, 100, color.RGBA{0, 0, 10, 255})
	if !canUseDarkimgKernel {
		t.Skip("no gpu available")
		defer func() {
			if !canUseDarkimgKernel {
				panic("unexpected")
			}
		}()
		return
	}
	result := IsDarkimg(img, 1.0)
	expected := cpuIsDarkimg(img, 1.0)
	if result != expected {
		t.Errorf("IsDarkimg 极暗蓝色结果 %v 与 cpuIsDarkimg %v 不一致", result, expected)
	}
}

func TestIsDarkimg_NonSquare(t *testing.T) {
	img := darkimgTestImage(300, 50, color.RGBA{200, 200, 200, 255})
	if !canUseDarkimgKernel {
		t.Skip("no gpu available")
		defer func() {
			if !canUseDarkimgKernel {
				panic("unexpected")
			}
		}()
		return
	}
	result := IsDarkimg(img, 0.5)
	expected := cpuIsDarkimg(img, 0.5)
	if result != expected {
		t.Errorf("IsDarkimg 非方形图像结果 %v 与 cpuIsDarkimg %v 不一致", result, expected)
	}
}

func TestIsDarkimg_MixedScales(t *testing.T) {
	if !canUseDarkimgKernel {
		t.Skip("no gpu available")
		defer func() {
			if !canUseDarkimgKernel {
				panic("unexpected")
			}
		}()
		return
	}
	img := darkimgPartialBrightImage(200, 200, 10)
	scales := []float32{0.1, 0.25, 0.5, 0.75, 1.0}
	for _, s := range scales {
		result := IsDarkimg(img, s)
		expected := cpuIsDarkimg(img, s)
		if result != expected {
			t.Errorf("IsDarkimg scale=%.2f 结果 %v 与 cpuIsDarkimg %v 不一致", s, result, expected)
		}
	}
}

// ---- 基准测试 ----

func BenchmarkCpuIsDarkimg_100x100(b *testing.B) {
	img := darkimgTestImage(100, 100, color.RGBA{10, 10, 10, 255})
	b.ResetTimer()
	for range b.N {
		cpuIsDarkimg(img, 1.0)
	}
}

func BenchmarkCpuIsDarkimg_500x500(b *testing.B) {
	img := darkimgTestImage(500, 500, color.RGBA{10, 10, 10, 255})
	b.ResetTimer()
	for range b.N {
		cpuIsDarkimg(img, 1.0)
	}
}

func BenchmarkCpuIsDarkimg_4096x4096(b *testing.B) {
	img := darkimgTestImage(4096, 4096, color.RGBA{10, 10, 10, 255})
	b.ResetTimer()
	canUseResizeKernel = false
	for range b.N {
		cpuIsDarkimg(img, 0.1)
	}
}

func BenchmarkIsDarkimg_100x100(b *testing.B) {
	img := darkimgTestImage(100, 100, color.RGBA{10, 10, 10, 255})
	b.ResetTimer()
	for range b.N {
		IsDarkimg(img, 1.0)
	}
}

func BenchmarkIsDarkimg_500x500(b *testing.B) {
	img := darkimgTestImage(500, 500, color.RGBA{10, 10, 10, 255})
	b.ResetTimer()
	for range b.N {
		IsDarkimg(img, 1.0)
	}
}

func BenchmarkIsDarkimg_4096x4096(b *testing.B) {
	img := darkimgTestImage(4096, 4096, color.RGBA{10, 10, 10, 255})
	b.ResetTimer()
	for range b.N {
		IsDarkimg(img, 0.1)
	}
}
