package main

import (
	"image"
	"image/color"
	"image/draw"
	"os"
	"testing"

	"github.com/Coloured-glaze/gg"
)

// go test -benchmem -bench .
// go test -benchmem -run=^$ -bench ^BenchmarkCircle2$
func BenchmarkCircle(b *testing.B) {
	f, _ := os.Open("../examples/gopher.png")
	//	gopherImg, _, _ := image.Decode(bufio.NewReader(f))
	gopherImg, _, _ := image.Decode(f)

	d := gopherImg.Bounds().Dx()
	b.ResetTimer() //重置时间
	for i := 0; i < b.N; i++ {
		c := circle{ //将一个cicle作为蒙层遮罩，圆心为图案中点，半径为边长的一半
			p: image.Point{
				X: d / 2,
				Y: d / 2,
			},
			r: d / 2,
		}
		circleImg := image.NewRGBA(image.Rect(0, 0, d, d))
		draw.DrawMask(circleImg, circleImg.Bounds(), gopherImg, image.Point{}, &c, image.Point{}, draw.Over)
	}
	//	gg.SavePNG("gopher_c.png", circleImg.SubImage(image.Rect(0, 0, d, d)))
	//	gg.SavePNG("gopher_c.png", circleImg)
	/*
		goos: windows
		goarch: amd64
		cpu: Intel(R) Core(TM) i3-10100F CPU @ 3.60GHz
		BenchmarkCircle-8   	     924	   1337658 ns/op	  430301 B/op	   70691 allocs/op
		PASS
		ok  	_/e_/1/github/gg/test	1.407s
	*/
}

func BenchmarkCircle2(b *testing.B) {
	var ft Factory
	one, _ := gg.LoadImage("../examples/gopher.png")
	ft.Im = gg.ImageToNRGBA(one)
	b.ResetTimer() //重置时间
	for i := 0; i < b.N; i++ {
		dx := ft.Im.Bounds().Dx()
		ft.ClipCircleFix(dx/2, dx/2, dx/2)
	}
	// gg.SavePNG("gopher_Circle.png", circleImg.SubImage(image.Rect(0, 0, dx, dx)))
	// gg.SavePNG("gopher_Circle.png", ft.Im)
	/*
		goos: windows
		goarch: amd64
		cpu: Intel(R) Core(TM) i3-10100F CPU @ 3.60GHz
		BenchmarkCircle2-8   	    6228	    170877 ns/op	   30532 B/op	    7618 allocs/op
		PASS
		ok  	_/e_/1/github/gg/test	1.114s
	*/
}

/////////////////////////////////////////////////////////////////////////

type circle struct { // 这里需要自己实现一个圆形遮罩，实现接口里的三个方法
	p image.Point // 圆心位置
	r int         // 半径
}

func (c *circle) ColorModel() color.Model {
	return color.AlphaModel
}
func (c *circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

// 对每个像素点进行色值设置，在半径以内的图案设成完全不透明
func (c *circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{A: 255}
	}
	return color.Alpha{}
}

//////////////////////////////////////////////////////////////////////////

type Factory struct {
	Im *image.NRGBA
}

// Clip 剪取方图
func (dst *Factory) Clip(w, h, x, y int) *Factory {
	dst.Im = dst.Im.SubImage(image.Rect(x, y, x+w, y+h)).(*image.NRGBA)
	return dst
}

// ClipCircleFix 裁取圆图 x,y中心点 r半径
func (dst *Factory) ClipCircleFix(x, y, r int) *Factory {
	dst = dst.Clip(2*r, 2*r, x-r, y-r)
	b := dst.Im.Bounds()
	for y1 := b.Min.Y; y1 < b.Max.Y; y1++ {
		for x1 := b.Min.X; x1 < b.Max.X; x1++ {

			if (x1-x)*(x1-x)+(y1-y)*(y1-y) > r*r {
				dst.Im.Set(x1, y1, color.NRGBA{0, 0, 0, 0})
			}
		}
	}
	return dst
}
