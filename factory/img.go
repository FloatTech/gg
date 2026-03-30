// Package factory provides image manipulation utilities built on top of gg.
//
// factory 包提供基于 gg 的图像操作工具库。
package factory

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"io"
	"net/http"
	"strings"

	"github.com/FloatTech/gg"
	"github.com/FloatTech/gg/fio"
	"github.com/disintegration/imaging"
)

// Load 加载图片
func Load(path string) (img image.Image, err error) {
	if strings.HasPrefix(path, "http") {
		var res *http.Response
		res, err = http.Get(path)
		if err != nil {
			return
		}
		img, _, err = image.Decode(res.Body)
		_ = res.Body.Close()
		return
	}
	return fio.LoadImage(path)
}

// Parse 解析图片数据流
func Parse(r io.Reader) (img image.Image, err error) {
	img, _, err = image.Decode(r)
	return
}

// NewFactoryBG creates a new Factory with a solid background color.
//
// NewFactoryBG 创建一个具有纯色背景的新 Factory。
func NewFactoryBG(w, h int, fillColor color.Color) *Factory {
	c := color.NRGBAModel.Convert(fillColor).(color.NRGBA)
	if (c == color.NRGBA{0, 0, 0, 0}) {
		return NewFactory(image.NewNRGBA(image.Rect(0, 0, w, h)))
	}
	return NewFactory(&image.NRGBA{
		Pix:    bytes.Repeat([]byte{c.R, c.G, c.B, c.A}, w*h),
		Stride: 4 * w,
		Rect:   image.Rect(0, 0, w, h),
	})
}

// LoadFirstFrame 载入图片第一帧作底图
func LoadFirstFrame(path string, w, h int) (*Factory, error) {
	im, err := Load(path)
	if err != nil {
		return nil, err
	}
	return Size(im, w, h), nil
}

// ParseFirstFrame 解析图片第一帧作底图
func ParseFirstFrame(r io.Reader, w, h int) (*Factory, error) {
	im, err := Parse(r)
	if err != nil {
		return nil, err
	}
	return Size(im, w, h), nil
}

// LoadAllFrames 加载每一帧图片
func LoadAllFrames(path string, w, h int) ([]*Factory, error) {
	var res *http.Response
	var err error
	var im *gif.GIF
	if strings.HasPrefix(path, "http") {
		res, err = http.Get(path)
		if err != nil {
			return nil, err
		}
		im, err = gif.DecodeAll(res.Body)
		_ = res.Body.Close()
		if err != nil {
			return nil, err
		}
	} else {
		im, err = fio.LoadGIF(path)
	}
	img, err := Load(path)
	if err != nil {
		return nil, err
	}
	im0 := Size(img, w, h)
	ims := make([]*Factory, len(im.Image))
	for i, v := range im.Image {
		ims[i] = im0.InsertUp(Size(v, w, h).im, 0, 0, 0, 0).Clone()
	}
	return ims, nil
}

// LoadAllTrueFrames 加载每一帧显示出的图片
func LoadAllTrueFrames(path string, w, h int) ([]*Factory, error) {
	var res *http.Response
	var err error
	var im *gif.GIF
	if strings.HasPrefix(path, "http") {
		res, err = http.Get(path)
		if err != nil {
			return nil, err
		}
		im, err = gif.DecodeAll(res.Body)
		_ = res.Body.Close()
		if err != nil {
			return nil, err
		}
	} else {
		im, err = fio.LoadGIF(path)
	}
	imgWidth, imgHeight := getGifDimensions(im)
	overpaintImage := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	img, err := Load(path)
	if err != nil {
		return nil, err
	}
	im0 := Size(img, w, h)
	ims := make([]*Factory, len(im.Image))
	for i, v := range im.Image {
		draw.Draw(overpaintImage, overpaintImage.Bounds(), v, image.Point{}, draw.Over)
		ims[i] = im0.InsertUp(Size(overpaintImage, w, h).im, 0, 0, 0, 0).Clone()
	}
	return ims, nil
}

// Size 变形
func Size(im image.Image, w, h int) *Factory {
	sz := im.Bounds().Size()
	// 修改尺寸
	switch {
	case w > 0 && h > 0:
		return NewFactory(imaging.Resize(im, w, h, imaging.Lanczos))
	case w == 0 && h > 0:
		return NewFactory(imaging.Resize(im, h*sz.X/sz.Y, h, imaging.Lanczos))
	case h == 0 && w > 0:
		return NewFactory(imaging.Resize(im, w, w*sz.Y/sz.X, imaging.Lanczos))
	default:
		nim := image.NewNRGBA(image.Rect(0, 0, sz.X, sz.Y))
		draw.Over.Draw(nim, nim.Bounds(), im, im.Bounds().Min)
		return NewFactory(nim)
	}
}

// Rotate 旋转
func Rotate(img image.Image, angle float64, w, h int) *Factory {
	return NewFactory(imaging.Rotate(Size(img, w, h).im, angle, color.NRGBA{0, 0, 0, 0}))
}

// MergeW 横向合并图片
func MergeW(im []*image.NRGBA) *Factory {
	dc := make([]*Factory, len(im))
	h := im[0].Bounds().Size().Y
	w := 0
	for i, value := range im {
		dc[i] = Size(value, 0, h)
		w += dc[i].W()
	}
	ds := NewFactoryBG(w, h, color.NRGBA{0, 0, 0, 0})
	x := 0
	for _, value := range dc {
		ds = ds.InsertUp(value.im, value.W(), h, x, 0)
		x += value.W()
	}
	return ds
}

// MergeH 纵向合并图片
func MergeH(im []*image.NRGBA) *Factory {
	dc := make([]*Factory, len(im))
	w := im[0].Bounds().Size().X
	h := 0
	for i, value := range im {
		dc[i] = Size(value, 0, w)
		h += dc[i].H()
	}
	ds := NewFactoryBG(w, h, color.NRGBA{0, 0, 0, 0})
	y := 0
	for _, value := range dc {
		ds = ds.InsertUp(value.im, w, value.H(), 0, y)
		y += value.H()
	}
	return ds
}

// Text 文本框 字体, 大小, 颜色 , 背景色, 文本
func Text(font string, size float64, col []int, col1 []int, txt string) *Factory {
	var dst Factory
	dc := gg.NewContext(10, 10)
	dc.SetRGBA255(0, 0, 0, 0)
	dc.Clear()
	dc.SetRGBA255(col[0], col[1], col[2], col[3])
	err := dc.LoadFontFace(font, size+size/2)
	if err != nil {
		return &dst
	}
	w, h := dc.MeasureString(txt)
	w -= size * 2
	dc1 := gg.NewContext(int(w), int(h))
	dc1.SetRGBA255(col1[0], col1[1], col1[2], col1[3])
	dc1.Clear()
	dc1.SetRGBA255(col[0], col[1], col[2], col[3])
	err = dc1.LoadFontFace(font, size)
	if err != nil {
		return &dst
	}
	dc1.DrawStringAnchored(txt, w/2, h/2, 0.5, 0.5)
	dst.im = image.NewNRGBA(image.Rect(0, 0, int(w), int(h)))
	draw.Over.Draw(dst.im, dst.im.Bounds(), dc1.Image(), dc1.Image().Bounds().Min)
	return &dst
}

// Limit 限制图片在 xmax*ymax 之内
func Limit(img image.Image, xmax, ymax int) image.Image {
	// 避免图片过大, 最大 xmax*ymax
	x := img.Bounds().Size().X
	y := img.Bounds().Size().Y
	hasChanged := false
	if x > xmax {
		y = y * xmax / x
		x = xmax
		hasChanged = true
	}
	if y > ymax {
		x = x * ymax / y
		y = ymax
		hasChanged = true
	}
	if hasChanged {
		img = Size(img, x, y).im
	}
	return img
}

// 获取gif真实大小
func getGifDimensions(gif *gif.GIF) (x, y int) {
	var lowestX int
	var lowestY int
	var highestX int
	var highestY int

	for _, img := range gif.Image {
		if img.Rect.Min.X < lowestX {
			lowestX = img.Rect.Min.X
		}
		if img.Rect.Min.Y < lowestY {
			lowestY = img.Rect.Min.Y
		}
		if img.Rect.Max.X > highestX {
			highestX = img.Rect.Max.X
		}
		if img.Rect.Max.Y > highestY {
			highestY = img.Rect.Max.Y
		}
	}

	return highestX - lowestX, highestY - lowestY
}
