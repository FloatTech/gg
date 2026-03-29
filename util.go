package gg

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"os"
	"strings"

	"github.com/fumiama/imgsz"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// Radians converts degrees to radians.
//
// Radians 将角度转换为弧度。
func Radians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// Degrees converts radians to degrees.
//
// Degrees 将弧度转换为角度。
func Degrees(radians float64) float64 {
	return radians * 180 / math.Pi
}

// ImageToRGBA converts an image.Image to *image.RGBA.
//
// ImageToRGBA 将 image.Image 转换为 *image.RGBA。
func ImageToRGBA(src image.Image) *image.RGBA {
	return imageToRGBA(src)
}

// image.Image 转为 image.RGBA
func imageToRGBA(src image.Image) *image.RGBA {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, src, bounds.Min, draw.Src)
	return dst
}

// ImageToRGBA64 converts an image.Image to *image.RGBA64.
//
// ImageToRGBA64 将 image.Image 转换为 *image.RGBA64。
func ImageToRGBA64(src image.Image) *image.RGBA64 {
	bounds := src.Bounds()
	dst := image.NewRGBA64(bounds)
	draw.Draw(dst, bounds, src, bounds.Min, draw.Src)
	return dst
}

// ImageToNRGBA converts an image.Image to *image.NRGBA.
//
// ImageToNRGBA 将 image.Image 转换为 *image.NRGBA。
func ImageToNRGBA(src image.Image) *image.NRGBA {
	bounds := src.Bounds()
	dst := image.NewNRGBA(bounds)
	draw.Draw(dst, bounds, src, bounds.Min, draw.Src)
	return dst
}

// ImageToNRGBA64 converts an image.Image to *image.NRGBA64.
//
// ImageToNRGBA64 将 image.Image 转换为 *image.NRGBA64。
func ImageToNRGBA64(src image.Image) *image.NRGBA64 {
	bounds := src.Bounds()
	dst := image.NewNRGBA64(bounds)
	draw.Draw(dst, bounds, src, bounds.Min, draw.Src)
	return dst
}

// GetImageWxH returns the width and height of the image at the given path.
//
// GetImageWxH 返回指定路径图片的宽度和高度。
func GetImageWxH(path string) (int, int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()
	sz, _, err := imgsz.DecodeSize(f)
	return sz.Width, sz.Height, err
}

// ParseHexColor parses a hex color string and returns r, g, b, a components.
// Supports 3, 6, and 8 digit hex strings with an optional leading '#'.
//
// ParseHexColor 解析十六进制颜色字符串并返回 r, g, b, a 分量。
// 支持 3、6 和 8 位十六进制字符串，前导 '#' 可选。
func ParseHexColor(x string) (r, g, b, a int) {
	x = strings.TrimPrefix(x, "#")
	a = 255
	switch len(x) {
	case 3:
		format := "%1x%1x%1x"
		_, _ = fmt.Sscanf(x, format, &r, &g, &b)
		r |= r << 4
		g |= g << 4
		b |= b << 4
	case 6:
		format := "%02x%02x%02x"
		_, _ = fmt.Sscanf(x, format, &r, &g, &b)
	case 8:
		format := "%02x%02x%02x%02x"
		_, _ = fmt.Sscanf(x, format, &r, &g, &b, &a)
	}
	return
}

func fixp(x, y float64) fixed.Point26_6 {
	return fixed.Point26_6{
		X: fix(x),
		Y: fix(y),
	}
}

func fix(x float64) fixed.Int26_6 {
	return fixed.Int26_6(math.Round(x * 64))
}

func unfix(x fixed.Int26_6) float64 {
	const shift, mask = 6, 1<<6 - 1
	if x >= 0 {
		return float64(x>>shift) + float64(x&mask)/64
	}
	x = -x
	if x >= 0 {
		return -(float64(x>>shift) + float64(x&mask)/64)
	}
	return 0
}

// LoadFontFace is a helper function to load the specified font file with
// the specified point size. Note that the returned `font.Face` objects
// are not thread safe and cannot be used in parallel across goroutines.
// You can usually just use the Context.LoadFontFace function instead of
// this package-level function.
//
// LoadFontFace 是一个辅助函数，用于加载指定点大小的指定字体文件。
// 请注意，返回的 `font.Face` 对象不是线程安全的，不能跨 goroutine 并行使用。
// 您通常可以只使用 Context.LoadFontFace 函数而不是这个包级函数。
func LoadFontFace(path string, points float64) (face font.Face, err error) {
	fontBytes, err := os.ReadFile(path)
	if err != nil {
		return
	}
	f, err := opentype.ParseCollection(fontBytes)
	if err != nil {
		return
	}
	fnf, err := f.Font(0)
	if err != nil {
		return
	}
	face, err = opentype.NewFace(fnf, &opentype.FaceOptions{
		Size: points,
		DPI:  72,
		// Hinting: font.HintingFull,
	})
	return
}

// ParseFontFace 是一个辅助函数，用于加载指定点大小的指定字体文件。
// 请注意，返回的 `font.Face` 对象不是线程安全的，不能跨 goroutine 并行使用。
// 您通常可以只使用 Context.LoadFontFace 函数而不是这个包级函数。
func ParseFontFace(b []byte, points float64) (face font.Face, err error) {
	f, err := opentype.ParseCollection(b)
	if err != nil {
		return
	}
	fnf, err := f.Font(0)
	if err != nil {
		return
	}
	face, err = opentype.NewFace(fnf, &opentype.FaceOptions{
		Size: points,
		DPI:  72,
		// Hinting: font.HintingFull,
	})
	return
}

// TakeColor extracts the k dominant colors from an image using k-means.
//
// TakeColor 使用 k-means 算法从图像中提取 k 个主色。
func TakeColor(img image.Image, k int) []color.RGBA {
	return takecolor(img, k)
}
