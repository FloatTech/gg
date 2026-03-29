package factory

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/FloatTech/gg"
	"github.com/disintegration/imaging"
)

// Factory 处理中图像
type Factory struct {
	im *image.NRGBA
}

func NewFactory(im *image.NRGBA) *Factory {
	return &Factory{im: im}
}

func (dst *Factory) W() int {
	return dst.im.Bounds().Size().X
}

func (dst *Factory) H() int {
	return dst.im.Bounds().Size().Y
}

func (dst *Factory) Image() *image.NRGBA {
	return dst.im
}

// Clone 克隆
func (dst *Factory) Clone() *Factory {
	var src Factory
	sz := dst.im.Bounds().Size()
	src.im = image.NewNRGBA(image.Rect(0, 0, sz.X, sz.Y))
	draw.Over.Draw(src.im, src.im.Bounds(), dst.im, dst.im.Bounds().Min)
	return &src
}

// Reshape 变形
func (dst *Factory) Reshape(w, h int) *Factory {
	dst = Size(dst.im, w, h)
	return dst
}

// FlipH 水平翻转
func (dst *Factory) FlipH() *Factory {
	return &Factory{
		im: imaging.FlipH(dst.im),
	}
}

// FlipV 垂直翻转
func (dst *Factory) FlipV() *Factory {
	return &Factory{
		im: imaging.FlipV(dst.im),
	}
}

// InsertUp 上部插入图片
func (dst *Factory) InsertUp(im image.Image, w, h, x, y int) *Factory {
	im1 := Size(im, w, h).im
	// 叠加图片
	draw.Over.Draw(dst.im, dst.im.Bounds(), im1, im1.Bounds().Min.Sub(image.Pt(x, y)))
	return dst
}

// InsertUpC 上部插入图片 x,y是中心点
func (dst *Factory) InsertUpC(im image.Image, w, h, x, y int) *Factory {
	im1 := Size(im, w, h)
	// 叠加图片
	draw.Over.Draw(dst.im, dst.im.Bounds(), im1.im, im1.im.Bounds().Min.Sub(image.Pt(x-im1.im.Bounds().Max.X/2, y-im1.im.Bounds().Max.Y/2)))
	return dst
}

// InsertBottom 底部插入图片
func (dst *Factory) InsertBottom(im image.Image, w, h, x, y int) *Factory {
	im1 := Size(im, w, h).im
	dc := dst.Clone()
	sz := dst.im.Bounds().Size()
	dst = NewFactoryBG(sz.X, sz.Y, color.NRGBA{0, 0, 0, 0})
	draw.Over.Draw(dst.im, dst.im.Bounds(), im1, im1.Bounds().Min.Sub(image.Pt(x, y)))
	draw.Over.Draw(dst.im, dst.im.Bounds(), dc.im, dc.im.Bounds().Min)
	return dst
}

// InsertBottomC 底部插入图片 x,y是中心点
func (dst *Factory) InsertBottomC(im image.Image, w, h, x, y int) *Factory {
	im1 := Size(im, w, h)
	dc := dst.Clone()
	sz := dst.im.Bounds().Size()
	dst = NewFactoryBG(sz.X, sz.Y, color.NRGBA{0, 0, 0, 0})
	draw.Over.Draw(dst.im, dst.im.Bounds(), im1.im, im1.im.Bounds().Min.Sub(image.Pt(x-im1.im.Bounds().Max.X/2, y-im1.im.Bounds().Max.Y/2)))
	draw.Over.Draw(dst.im, dst.im.Bounds(), dc.im, dc.im.Bounds().Min)
	return dst
}

// Circle 获取圆图
func (dst *Factory) Circle(r int) *Factory {
	sz := dst.im.Bounds().Size()
	if r == 0 {
		r = sz.Y / 2
	}
	dst = dst.Reshape(2*r, 2*r)
	b := dst.im.Bounds()
	for y1 := b.Min.Y; y1 < b.Max.Y; y1++ {
		for x1 := b.Min.X; x1 < b.Max.X; x1++ {
			if (x1-r)*(x1-r)+(y1-r)*(y1-r) > r*r {
				dst.im.Set(x1, y1, color.NRGBA{0, 0, 0, 0})
			}
		}
	}
	return dst
}

// Clip 剪取方图
func (dst *Factory) Clip(w, h, x, y int) *Factory {
	dst.im = dst.im.SubImage(image.Rect(x, y, x+w, y+h)).(*image.NRGBA)
	return dst
}

// ClipCircleFix 裁取圆图
func (dst *Factory) ClipCircleFix(x, y, r int) *Factory {
	dst = dst.Clip(2*r, 2*r, x-r, y-r)
	b := dst.im.Bounds()
	for y1 := b.Min.Y; y1 < b.Max.Y; y1++ {
		for x1 := b.Min.X; x1 < b.Max.X; x1++ {
			if (x1-x)*(x1-x)+(y1-y)*(y1-y) > r*r {
				dst.im.Set(x1, y1, color.NRGBA{0, 0, 0, 0})
			}
		}
	}
	return dst
}

// ClipCircle 扣取圆
func (dst *Factory) ClipCircle(x, y, r int) *Factory {
	//  dc := dst.Clip(x-r, y-r, 2*r, 2*r)
	b := dst.im.Bounds()
	for y1 := b.Min.Y; y1 < b.Max.Y; y1++ {
		for x1 := b.Min.X; x1 < b.Max.X; x1++ {
			if (x1-x)*(x1-x)+(y1-y)*(y1-y) <= r*r {
				dst.im.Set(x1, y1, color.NRGBA{0, 0, 0, 0})
			}
		}
	}
	return dst
}

// InsertText 插入文本
func (dst *Factory) InsertText(font string, size float64, col []int, x, y float64, txt string) *Factory {
	dc := gg.NewContextForImage(dst.im)
	// 字体, 大小, 颜色, 位置
	err := dc.LoadFontFace(font, size)
	if err != nil {
		return dst
	}
	dc.SetRGBA255(col[0], col[1], col[2], col[3])
	dc.DrawString(txt, x, y)
	ds := dc.Image()
	draw.Over.Draw(dst.im, dst.im.Bounds(), ds, ds.Bounds().Min)
	return dst
}

// InsertUpG gif 上部插入图片
func (dst *Factory) InsertUpG(im []*image.NRGBA, w, h, x, y int) []*image.NRGBA {
	if len(im) == 0 {
		return nil
	}
	ims := make([]*image.NRGBA, len(im))
	for i, v := range im {
		ims[i] = dst.Clone().InsertUp(v, w, h, x, y).im
	}
	return ims
}
