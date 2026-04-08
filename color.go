package gg

import (
	"image"
	"image/color"
	"unsafe"
)

// Predefined colors.
//
// 预定义颜色。
var (
	White   = color.RGBA{255, 255, 255, 255}
	Black   = color.RGBA{0, 0, 0, 255}
	Red     = color.RGBA{255, 0, 0, 255}
	Green   = color.RGBA{0, 255, 0, 255}
	Blue    = color.RGBA{0, 0, 255, 255}
	Yellow  = color.RGBA{255, 255, 0, 255}
	Cyan    = color.RGBA{0, 255, 255, 255}
	Magenta = color.RGBA{255, 0, 255, 255}
	Grey    = color.RGBA{190, 190, 190, 255}
	Pink    = color.RGBA{255, 181, 197, 255}
	Orange  = color.RGBA{255, 165, 0, 255}
	Opaque  = color.RGBA{0, 0, 0, 0}
)

// TakeThemeColorsKMeans extracts the k dominant colors from an image using k-means.
//
// TakeThemeColorsKMeans 使用 k-means 算法从图像中提取 k 个主色。
func TakeThemeColorsKMeans(img image.Image, k uint16) ([]color.RGBA, error) {
	ki, err := newKMeansImage(img, k) // 初始化k个聚类中心
	if err != nil {
		return nil, err
	}
	defer ki.destroy()
	for {
		if err := ki.assign(); err != nil {
			return nil, err
		}
		ki.update()
		if ki.epilogue() {
			break
		}
	}
	return ki.result(), nil
}

// isArrayRGBAEqual compares two []color.RGBA is equal fastly.
//
// isArrayRGBAEqual 快速比较两个 []color.RGBA 是否相等。
func isArrayRGBAEqual(a, b []color.RGBA) bool {
	if len(a) != len(b) {
		return false
	}
	sz := len(a)
	if sz%2 == 0 { // can compare by uint64
		u64a := unsafe.Slice((*uint64)(unsafe.Pointer(unsafe.SliceData(a))), sz/2)
		u64b := unsafe.Slice((*uint64)(unsafe.Pointer(unsafe.SliceData(b))), sz/2)
		for i := range u64a {
			if u64a[i] != u64b[i] {
				return false
			}
		}
		return true
	}
	// compare by uint32
	u32a := unsafe.Slice((*uint32)(unsafe.Pointer(unsafe.SliceData(a))), sz)
	u32b := unsafe.Slice((*uint32)(unsafe.Pointer(unsafe.SliceData(b))), sz)
	for i := range u32a {
		if u32a[i] != u32b[i] {
			return false
		}
	}
	return true
}

// distanceRGBAsq calc between two color.RGBAs (RGB only, no sqrt to speedup)
//
// distanceRGBAsq 计算两个 color.RGBA 颜色之间的距离（仅 RGB，忽略开方以加速）
func distanceRGBAsq(a, b color.RGBA) float64 {
	return sq(float64(a.R)-float64(b.R)) +
		sq(float64(a.G)-float64(b.G)) +
		sq(float64(a.B)-float64(b.B))
}
