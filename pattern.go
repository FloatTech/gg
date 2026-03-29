package gg

import (
	"image"
	"image/color"

	"github.com/golang/freetype/raster"
)

// RepeatOp defines how a surface pattern repeats.
//
// RepeatOp 定义表面图案的重复方式。
type RepeatOp int

// Pattern repeat modes.
//
// 图案重复模式。
const (
	RepeatBoth RepeatOp = iota // Repeat in both directions. 在两个方向上重复。
	RepeatX                    // Repeat horizontally only. 仅水平重复。
	RepeatY                    // Repeat vertically only. 仅垂直重复。
	RepeatNone                 // No repeat. 不重复。
)

// Pattern defines an interface for generating colors at given coordinates.
//
// Pattern 定义了在给定坐标生成颜色的接口。
type Pattern interface {
	ColorAt(x, y int) color.Color
}

// Solid Pattern
type solidPattern struct {
	color color.Color
}

func (p *solidPattern) ColorAt(_, _ int) color.Color {
	return p.color
}

// NewSolidPattern creates a pattern that always returns the given color.
//
// NewSolidPattern 创建一个始终返回指定颜色的图案。
func NewSolidPattern(color color.Color) Pattern {
	return &solidPattern{color: color}
}

// Surface Pattern
type surfacePattern struct {
	im image.Image
	op RepeatOp
}

func (p *surfacePattern) ColorAt(x, y int) color.Color {
	b := p.im.Bounds()
	switch p.op {
	case RepeatX:
		if y >= b.Dy() {
			return color.Transparent
		}
	case RepeatY:
		if x >= b.Dx() {
			return color.Transparent
		}
	case RepeatNone:
		if x >= b.Dx() || y >= b.Dy() {
			return color.Transparent
		}
	}
	x = x%b.Dx() + b.Min.X
	y = y%b.Dy() + b.Min.Y
	return p.im.At(x, y)
}

// NewSurfacePattern creates a pattern from an image with the given repeat mode.
//
// NewSurfacePattern 使用指定重复模式从图像创建图案。
func NewSurfacePattern(im image.Image, op RepeatOp) Pattern {
	return &surfacePattern{im: im, op: op}
}

type patternPainter struct {
	im   *image.RGBA
	mask *image.Alpha
	p    Pattern
}

// Paint satisfies the Painter interface.
func (r *patternPainter) Paint(ss []raster.Span, _ bool) {
	b := r.im.Bounds()
	for _, s := range ss {
		if s.Y < b.Min.Y {
			continue
		}
		if s.Y >= b.Max.Y {
			return
		}
		if s.X0 < b.Min.X {
			s.X0 = b.Min.X
		}
		if s.X1 > b.Max.X {
			s.X1 = b.Max.X
		}
		if s.X0 >= s.X1 {
			continue
		}
		const m = 1<<16 - 1
		y := s.Y - r.im.Rect.Min.Y
		x0 := s.X0 - r.im.Rect.Min.X
		// RGBAPainter.Paint() in $GOPATH/src/github.com/golang/freetype/raster/paint.go
		i0 := (s.Y-r.im.Rect.Min.Y)*r.im.Stride + (s.X0-r.im.Rect.Min.X)*4
		i1 := i0 + (s.X1-s.X0)*4
		for i, x := i0, x0; i < i1; i, x = i+4, x+1 {
			ma := s.Alpha
			if r.mask != nil {
				ma = ma * uint32(r.mask.AlphaAt(x, y).A) / 255
				if ma == 0 {
					continue
				}
			}
			c := r.p.ColorAt(x, y)
			cr, cg, cb, ca := c.RGBA()
			dr := uint32(r.im.Pix[i+0])
			dg := uint32(r.im.Pix[i+1])
			db := uint32(r.im.Pix[i+2])
			da := uint32(r.im.Pix[i+3])
			a := (m - (ca * ma / m)) * 0x101
			r.im.Pix[i+0] = uint8((dr*a + cr*ma) / m >> 8)
			r.im.Pix[i+1] = uint8((dg*a + cg*ma) / m >> 8)
			r.im.Pix[i+2] = uint8((db*a + cb*ma) / m >> 8)
			r.im.Pix[i+3] = uint8((da*a + ca*ma) / m >> 8)
		}
	}
}

func newPatternPainter(im *image.RGBA, mask *image.Alpha, p Pattern) *patternPainter {
	return &patternPainter{im, mask, p}
}
