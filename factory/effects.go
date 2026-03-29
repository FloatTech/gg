package factory

import (
	"image/color"

	"github.com/disintegration/imaging"
)

// float64转uint8
func floatUint8(a float64) uint8 {
	b := int64(a + 0.5)
	if b > 255 {
		return 255
	}
	if b > 0 {
		return uint8(b)
	}
	return 0
}

// AdjustBrightness 亮度(-100, 100)
func (dst *Factory) AdjustBrightness(s float64) *Factory {
	return &Factory{
		im: imaging.AdjustBrightness(dst.im, s),
	}
}

// AdjustContrast 对比度(-100, 100)
func (dst *Factory) AdjustContrast(a float64) *Factory {
	return &Factory{
		im: imaging.AdjustContrast(dst.im, a),
	}
}

// AdjustSaturation 饱和度(-100, 100)
func (dst *Factory) AdjustSaturation(a float64) *Factory {
	return &Factory{
		im: imaging.AdjustSaturation(dst.im, a),
	}
}

// Sharpen 锐化
func (dst *Factory) Sharpen(a float64) *Factory {
	return &Factory{
		im: imaging.Sharpen(dst.im, a),
	}
}

// Blur 模糊图像 正数
func (dst *Factory) Blur(a float64) *Factory {
	return &Factory{
		im: imaging.Blur(dst.im, a),
	}
}

// Grayscale 灰度
func (dst *Factory) Grayscale() *Factory {
	b := dst.im.Bounds()
	for y1 := b.Min.Y; y1 <= b.Max.Y; y1++ {
		for x1 := b.Min.X; x1 <= b.Max.X; x1++ {
			a := dst.im.At(x1, y1)
			c := color.NRGBAModel.Convert(a).(color.NRGBA)
			f := 0.299*float64(c.R) + 0.587*float64(c.G) + 0.114*float64(c.B)
			c.R = floatUint8(f)
			c.G = floatUint8(f)
			c.B = floatUint8(f)
			dst.im.Set(x1, y1, c)
		}
	}
	return dst
}

// Invert 反色
func (dst *Factory) Invert() *Factory {
	b := dst.im.Bounds()
	for y1 := b.Min.Y; y1 <= b.Max.Y; y1++ {
		for x1 := b.Min.X; x1 <= b.Max.X; x1++ {
			a := dst.im.At(x1, y1)
			c := color.NRGBAModel.Convert(a).(color.NRGBA)
			c.R = 255 - c.R
			c.G = 255 - c.G
			c.B = 255 - c.B
			dst.im.Set(x1, y1, c)
		}
	}
	return dst
}

// Relief 浮雕
func (dst *Factory) Relief() *Factory {
	return &Factory{
		im: imaging.Convolve3x3(
			dst.im,
			[9]float64{
				-1, -1, 0,
				-1, 1, 1,
				0, 1, 1,
			},
			nil,
		),
	}
}
