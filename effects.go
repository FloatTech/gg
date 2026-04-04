package gg

import (
	"image"
	"math"

	"github.com/disintegration/imaging"
)

// AdjustBrightness 调整亮度 范围：±100%
func (dc *Context) AdjustBrightness(s float64) {
	if math.Abs(s) < 0.001 {
		return
	}
	dc.im = (*image.RGBA)(imaging.AdjustBrightness(dc.im, s))
}

// AdjustContrast 调整对比度 范围：±100%
func (dc *Context) AdjustContrast(s float64) {
	if math.Abs(s) < 0.001 {
		return
	}
	dc.im = (*image.RGBA)(imaging.AdjustContrast(dc.im, s))
}

// AdjustContrast 调整饱和度 范围：±100%
func (dc *Context) AdjustSaturation(s float64) {
	if math.Abs(s) < 0.001 {
		return
	}
	dc.im = (*image.RGBA)(imaging.AdjustSaturation(dc.im, s))
}

// Sharpen 锐化 范围：±100%
func (dc *Context) Sharpen(s float64) {
	if math.Abs(s) < 0.001 {
		return
	}
	dc.im = (*image.RGBA)(imaging.Sharpen(dc.im, s))
}

// Blur 模糊图像 正数
func (dc *Context) Blur(s float64) {
	if math.Abs(s) < 0.001 {
		return
	}
	dc.im = (*image.RGBA)(imaging.Blur(dc.im, s))
}
