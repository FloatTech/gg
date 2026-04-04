package gg

import (
	"image"
	"os"

	"github.com/fumiama/imgsz"
	"golang.org/x/image/draw"
)

// ImageToRGBA converts an image.Image to *image.RGBA.
//
// ImageToRGBA 将 image.Image 转换为 *image.RGBA。
func ImageToRGBA(src image.Image) *image.RGBA {
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

// ImageBoundsBelow returns resized image that newW < w and newH < h, while keeping the W/H ratio.
//
// ImageBoundsBelow 返回在保持宽高比的条件下，小于 w x h 的新 bound。
func ImageBoundsBelow(b image.Rectangle, w, h int) image.Rectangle {
	width := b.Dx()
	height := b.Dy()
	dstw, dsth := width, height
	if dstw > w {
		dstw = w
		ratio := float64(dstw) / float64(width)
		dsth *= int(float64(height) * ratio)
	}
	if dsth > h {
		dsth = h
		ratio := float64(dsth) / float64(height)
		dstw = int(float64(width) * ratio)
	}
	return image.Rect(0, 0, dstw, dsth)
}
