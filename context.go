// Package gg provides a simple API for rendering 2D graphics in pure Go.
// 包 gg 提供了一个简单的API，用于在纯Go中渲染二维图形。
package gg

import (
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"strings"

	"github.com/golang/freetype/raster"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/f64"
)

type LineCap int

const (
	LineCapRound LineCap = iota
	LineCapButt
	LineCapSquare
)

type LineJoin int

const (
	LineJoinRound LineJoin = iota
	LineJoinBevel
)

type FillRule int

const (
	FillRuleWinding FillRule = iota
	FillRuleEvenOdd
)

type Align int

const (
	AlignLeft Align = iota
	AlignCenter
	AlignRight
)

var (
	defaultFillStyle   = NewSolidPattern(color.White)
	defaultStrokeStyle = NewSolidPattern(color.Black)
)

type Context struct {
	width         int
	height        int
	rasterizer    *raster.Rasterizer
	im            *image.RGBA
	mask          *image.Alpha
	color         color.Color
	fillPattern   Pattern
	strokePattern Pattern
	strokePath    raster.Path
	fillPath      raster.Path
	start         Point
	current       Point
	hasCurrent    bool
	dashes        []float64
	dashOffset    float64
	lineWidth     float64
	lineCap       LineCap
	lineJoin      LineJoin
	fillRule      FillRule
	fontFace      font.Face
	fontHeight    float64
	matrix        Matrix
	stack         []*Context
	scaleStyle    ScaleStyle
}

// NewContext creates a new image.RGBA with the specified width and height
// and prepares a context for rendering onto that image.
//
// 创建一个具有指定宽度和高度的新 image.RGBA
// 并准备渲染到该图像上的上下文。
func NewContext(width, height int) *Context {
	return NewContextForRGBA(image.NewRGBA(image.Rect(0, 0, width, height)))
}

// NewContextForImage copies the specified image into a new image.RGBA
// and prepares a context for rendering onto that image.
//
// 将指定图像复制到一个新的 image.RGBA
// 并准备渲染到该图像上的上下文。
func NewContextForImage(im image.Image) *Context {
	return NewContextForRGBA(imageToRGBA(im))
}

// NewContextForRGBA prepares a context for rendering onto the specified image.
// No copy is made.
//
// 准备渲染到指定图像的上下文。没有复制。
func NewContextForRGBA(im *image.RGBA) *Context {
	w := im.Bounds().Size().X
	h := im.Bounds().Size().Y
	return &Context{
		width:         w,
		height:        h,
		rasterizer:    raster.NewRasterizer(w, h),
		im:            im,
		color:         color.Transparent,
		fillPattern:   defaultFillStyle,
		strokePattern: defaultStrokeStyle,
		lineWidth:     1,
		fillRule:      FillRuleWinding,
		fontFace:      basicfont.Face7x13,
		fontHeight:    13,
		matrix:        Identity(),
	}
}

// ScaleStyle determines the way image pixels are interpolated when scaled.
// See
//   https://pkg.go.dev/golang.org/x/image/draw
// for the corresponding interpolators.
//
// 确定缩放时图像像素的插值方式。 请看
// https://pkg.go.dev/golang.org/x/image/draw
// 对应的插值器。
type ScaleStyle int

const (
	// BiLinear is the tent kernel. It is slow, but usually gives high quality
	// results.
	//
	// BiLinear 是帐篷内核。 它很慢，但通常会产生高质量的结果。
	BiLinear ScaleStyle = iota

	// ApproxBiLinear is a mixture of the nearest neighbor and bi-linear
	// interpolators. It is fast, but usually gives medium quality results.
	//
	// It implements bi-linear interpolation when upscaling and a bi-linear
	// blend of the 4 nearest neighbor pixels when downscaling. This yields
	// nicer quality than nearest neighbor interpolation when upscaling, but
	// the time taken is independent of the number of source pixels, unlike the
	// bi-linear interpolator. When downscaling a large image, the performance
	// difference can be significant.
	//
	// ApproxBiLinear 是最近邻和双线性插值器的混合。 它速度很快，但通常会给出中等质量的结果。
	// 它在放大时实现双线性插值，在缩小时实现 4 个最近邻像素的双线性混合。
	// 这在放大时产生比最近邻插值更好的质量，但与双线性插值器不同，所花费的时间与源像素的数量无关。
	// 缩小大图像时，性能差异可能很大。
	ApproxBiLinear

	// NearestNeighbor is the nearest neighbor interpolator. It is very fast,
	// but usually gives very low quality results. When scaling up, the result
	// will look 'blocky'.
	//
	// NearestNeighbor 是最近邻插值器。 它非常快，但通常会给出非常低质量的结果。
	// 放大时，结果将看起来“块状”。
	NearestNeighbor

	// CatmullRom is the Catmull-Rom kernel. It is very slow, but usually gives
	// very high quality results.
	//
	// It is an instance of the more general cubic BC-spline kernel with parameters
	// B=0 and C=0.5. See Mitchell and Netravali, "Reconstruction Filters in
	// Computer Graphics", Computer Graphics, Vol. 22, No. 4, pp. 221-228.
	//
	// CatmullRom 是 Catmull-Rom 内核。 它很慢，但通常会给出非常高质量的结果。
	// 它是参数 B=0 和 C=0.5 的更一般的三次 BC 样条核的实例。
	// 参见 Mitchell 和 Netravali，“计算机图形学中的重构过滤器”，
	// 计算机图形学，卷.22，第 4 期，第 221-228 页。
	CatmullRom
)

// 变压器
func (s ScaleStyle) transformer() draw.Interpolator {
	switch s {
	case BiLinear:
		return draw.BiLinear
	case ApproxBiLinear:
		return draw.ApproxBiLinear
	case NearestNeighbor:
		return draw.NearestNeighbor
	case CatmullRom:
		return draw.CatmullRom
	}
	return draw.BiLinear // BiLinear by default. 默认情况下为双线性。
}

// 设置比缩放样式
func (dc *Context) SetScaleStyle(s ScaleStyle) {
	dc.scaleStyle = s
}
// 设置缩放双线性
func (dc *Context) SetScaleBiLinear() {
	dc.SetScaleStyle(BiLinear)
}
// 设置缩放近似双线性
func (dc *Context) SetScaleApproxBiLinear() {
	dc.SetScaleStyle(ApproxBiLinear)
}
// 设置最近邻的缩放
func (dc *Context) SetScaleNearestNeighbor() {
	dc.SetScaleStyle(NearestNeighbor)
}
// 设置缩放为 CatmullRom
func (dc *Context) SetScaleCatmullRom() {
	dc.SetScaleStyle(CatmullRom)
}

// GetCurrentPoint will return the current point and if there is a current point.
// The point will have been transformed by the context's transformation matrix.
//
//GetCurrentPoint将返回当前点，如果存在当前点。
//该点将通过上下文的变换矩阵进行变换。
func (dc *Context) GetCurrentPoint() (Point, bool) {
	if dc.hasCurrent {
		return dc.current, true
	}
	return Point{}, false
}

// Image returns the image that has been drawn by this context.
//
// 返回在此上下文中绘制的图像。
func (dc *Context) Image() image.Image {
	return dc.im
}

// Width returns the width of the image in pixels.
//
// 返回图像的宽度（以像素为单位）
func (dc *Context) Width() int {
	return dc.width
}

// Height returns the height of the image in pixels.
//
// 返回图像的高度（以像素为单位）
func (dc *Context) Height() int {
	return dc.height
}

// SavePNG encodes the image as a PNG and writes it to disk.
//
// 将图像编码为 PNG 并将其写入磁盘。
func (dc *Context) SavePNG(path string) error {
	return SavePNG(path, dc.im)
}

// SaveJPG encodes the image as a JPG and writes it to disk.
//
// 将图像编码为 JPG 并将其写入磁盘。
func (dc *Context) SaveJPG(path string, quality int) error {
	return SaveJPG(path, dc.im, quality)
}

// EncodePNG encodes the image as a PNG and writes it to the provided io.Writer.
//
// 将图像编码为 PNG 并将其写入提供的 io.Writer
func (dc *Context) EncodePNG(w io.Writer) error {
	return png.Encode(w, dc.im)
}

// EncodeJPG encodes the image as a JPG and writes it to the provided io.Writer
// in JPEG 4:2:0 baseline format with the given options.
// Default parameters are used if a nil *jpeg.Options is passed.
//
// 将图像编码为JPG，并将其写入提供的 io.Writer
// 在 JPEG 4:2:0 基线格式中，使用给定选项。
// 如果为nil，则使用默认参数 *jpeg.Options 进行传递
func (dc *Context) EncodeJPG(w io.Writer, o *jpeg.Options) error {
	return jpeg.Encode(w, dc.im, o)
}

// SetDash sets the current dash pattern to use. Call with zero arguments to
// disable dashes. The values specify the lengths of each dash, with
// alternating on and off lengths.
//
// SetDash设置要使用的当前破折号图案。使用零参数调用
// 禁用破折号。这些值指定每个破折号的长度，包括：
// 交替开启和关闭长度。
func (dc *Context) SetDash(dashes ...float64) {
	dc.dashes = dashes
}

// SetDashOffset sets the initial offset into the dash pattern to use when
// stroking dashed paths.
//
// 将初始偏移量设置为虚线模式，以在描边虚线路径时使用。
func (dc *Context) SetDashOffset(offset float64) {
	dc.dashOffset = offset
}

// 设置线宽
func (dc *Context) SetLineWidth(lineWidth float64) {
	dc.lineWidth = lineWidth
}
// 设置线帽
func (dc *Context) SetLineCap(lineCap LineCap) {
	dc.lineCap = lineCap
}
// 设置线帽圆
func (dc *Context) SetLineCapRound() {
	dc.lineCap = LineCapRound
}
// 设置线帽对齐
func (dc *Context) SetLineCapButt() {
	dc.lineCap = LineCapButt
}
// 设置线帽正方形
func (dc *Context) SetLineCapSquare() {
	dc.lineCap = LineCapSquare
}
// 设置线帽连接
func (dc *Context) SetLineJoin(lineJoin LineJoin) {
	dc.lineJoin = lineJoin
}
// 设置线帽连接圆
func (dc *Context) SetLineJoinRound() {
	dc.lineJoin = LineJoinRound
}
/// 设置线帽连接斜面
func (dc *Context) SetLineJoinBevel() {
	dc.lineJoin = LineJoinBevel
}
// 设置填充规则
func (dc *Context) SetFillRule(fillRule FillRule) {
	dc.fillRule = fillRule
}
// 设置填充规则绕组
func (dc *Context) SetFillRuleWinding() {
	dc.fillRule = FillRuleWinding
}
// 设置填充规则偶数奇数
func (dc *Context) SetFillRuleEvenOdd() {
	dc.fillRule = FillRuleEvenOdd
}

// Color Setters
// 颜色设定器

// 设置填充和描边颜色
func (dc *Context) setFillAndStrokeColor(c color.Color) {
	dc.color = c
	dc.fillPattern = NewSolidPattern(c)
	dc.strokePattern = NewSolidPattern(c)
}

// SetFillStyle sets current fill style
// 设置当前填充样式
func (dc *Context) SetFillStyle(pattern Pattern) {
	// if pattern is SolidPattern, also change dc.color(for dc.Clear, dc.drawString)
	if fillStyle, ok := pattern.(*solidPattern); ok {
		dc.color = fillStyle.color
	}
	dc.fillPattern = pattern
}

// SetStrokeStyle sets current stroke style
// 设置当前笔划样式
func (dc *Context) SetStrokeStyle(pattern Pattern) {
	dc.strokePattern = pattern
}

// SetColor sets the current color(for both fill and stroke).
// 设置当前颜色（用于填充和笔划）
func (dc *Context) SetColor(c color.Color) {
	dc.setFillAndStrokeColor(c)
}

// SetHexColor sets the current color using a hex string. The leading pound
// sign (#) is optional. Both 3- and 6-digit variations are supported. 8 digits
// may be provided to set the alpha value as well.
//
// 使用十六进制字符串设置当前颜色。前置的
// 符号（#）是可选的。支持3位数和6位数的变体。8位数字
// 也可以提供设置 alpha 值。
func (dc *Context) SetHexColor(x string) {
	r, g, b, a := parseHexColor(x)
	dc.SetRGBA255(r, g, b, a)
}

// SetRGBA255 sets the current color. r, g, b, a values should be between 0 and
// 255, inclusive.
//
// 设置当前颜色。r、g、b、a 值应介于 0 和 255，包括在内。
func (dc *Context) SetRGBA255(r, g, b, a int) {
	dc.color = color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	dc.setFillAndStrokeColor(dc.color)
}

// SetRGB255 sets the current color. r, g, b values should be between 0 and 255,
// inclusive. Alpha will be set to 255 (fully opaque).
//
// 设置当前颜色。 r、g、b 值应介于 0 和 255 之间，包括。
// Alpha 将设置为 255（完全不透明）。
func (dc *Context) SetRGB255(r, g, b int) {
	dc.SetRGBA255(r, g, b, 255)
}

// SetRGBA sets the current color. r, g, b, a values should be between 0 and 1,
// inclusive.
//
// SetRGBA 设置当前颜色。 r、g、b、a 值应介于 0 和 1 之间，
// 包括。
func (dc *Context) SetRGBA(r, g, b, a float64) {
	dc.color = color.NRGBA{
		uint8(r * 255),
		uint8(g * 255),
		uint8(b * 255),
		uint8(a * 255),
	}
	dc.setFillAndStrokeColor(dc.color)
}

// SetRGB sets the current color. r, g, b values should be between 0 and 1,
// inclusive. Alpha will be set to 1 (fully opaque).
//
// SetRGB 设置当前颜色。 r、g、b 值应介于 0 和 1 之间，
// 包括。 Alpha 将设置为 1（完全不透明）。
func (dc *Context) SetRGB(r, g, b float64) {
	dc.SetRGBA(r, g, b, 1)
}

// Path Manipulation
// 路径操作

// MoveTo starts a new subpath within the current path starting at the
// specified point.
//
// 在当前路径中从指定点开始新的子路径。
func (dc *Context) MoveTo(x, y float64) {
	if dc.hasCurrent {
		dc.fillPath.Add1(dc.start.Fixed())
	}
	x, y = dc.TransformPoint(x, y)
	p := Point{x, y}
	dc.strokePath.Start(p.Fixed())
	dc.fillPath.Start(p.Fixed())
	dc.start = p
	dc.current = p
	dc.hasCurrent = true
}

// LineTo adds a line segment to the current path starting at the current
// point. If there is no current point, it is equivalent to MoveTo(x, y)
//
// 从当前点开始向当前路径添加一条线段。 
// 如果没有当前点，则等价于 MoveTo(x, y)
func (dc *Context) LineTo(x, y float64) {
	if !dc.hasCurrent {
		dc.MoveTo(x, y)
	} else {
		x, y = dc.TransformPoint(x, y)
		p := Point{x, y}
		dc.strokePath.Add1(p.Fixed())
		dc.fillPath.Add1(p.Fixed())
		dc.current = p
	}
}

// QuadraticTo adds a quadratic bezier curve to the current path starting at
// the current point. If there is no current point, it first performs
// MoveTo(x1, y1)
//
// 将二次贝塞尔曲线添加到从当前点开始的当前路径。 
// 如果没有当前点，则首先执行 MoveTo(x1, y1)
func (dc *Context) QuadraticTo(x1, y1, x2, y2 float64) {
	if !dc.hasCurrent {
		dc.MoveTo(x1, y1)
	}
	x1, y1 = dc.TransformPoint(x1, y1)
	x2, y2 = dc.TransformPoint(x2, y2)
	p1 := Point{x1, y1}
	p2 := Point{x2, y2}
	dc.strokePath.Add2(p1.Fixed(), p2.Fixed())
	dc.fillPath.Add2(p1.Fixed(), p2.Fixed())
	dc.current = p2
}

// CubicTo adds a cubic bezier curve to the current path starting at the
// current point. If there is no current point, it first performs
// MoveTo(x1, y1). Because freetype/raster does not support cubic beziers,
// this is emulated with many small line segments.
//
// 向当前路径添加一条三次贝塞尔曲线，从
// 当前点。 如果没有当前点，则首先执行
// MoveTo(x1, y1) 因为 freetype/raster 不支持三次贝塞尔曲线，
// 这是用许多小线段模拟的。
func (dc *Context) CubicTo(x1, y1, x2, y2, x3, y3 float64) {
	if !dc.hasCurrent {
		dc.MoveTo(x1, y1)
	}
	x0, y0 := dc.current.X, dc.current.Y
	x1, y1 = dc.TransformPoint(x1, y1)
	x2, y2 = dc.TransformPoint(x2, y2)
	x3, y3 = dc.TransformPoint(x3, y3)
	points := CubicBezier(x0, y0, x1, y1, x2, y2, x3, y3)
	previous := dc.current.Fixed()
	for _, p := range points[1:] {
		f := p.Fixed()
		if f == previous {
			// TODO: this fixes some rendering issues but not all
			continue
		}
		previous = f
		dc.strokePath.Add1(f)
		dc.fillPath.Add1(f)
		dc.current = p
	}
}

// ClosePath adds a line segment from the current point to the beginning
// of the current subpath. If there is no current point, this is a no-op.
//
// 添加从当前点到当前子路径开头的线段。 
// 如果没有当前点，这是一个空操作。
func (dc *Context) ClosePath() {
	if dc.hasCurrent {
		dc.strokePath.Add1(dc.start.Fixed())
		dc.fillPath.Add1(dc.start.Fixed())
		dc.current = dc.start
	}
}

// ClearPath clears the current path. There is no current point after this
// operation.
//
// 清除当前路径。 此操作后没有当前点。
func (dc *Context) ClearPath() {
	dc.strokePath.Clear()
	dc.fillPath.Clear()
	dc.hasCurrent = false
}

// NewSubPath starts a new subpath within the current path. There is no current
// point after this operation.
//
// 在当前路径中开始一个新的子路径。 此操作后没有当前点。
func (dc *Context) NewSubPath() {
	if dc.hasCurrent {
		dc.fillPath.Add1(dc.start.Fixed())
	}
	dc.hasCurrent = false
}

// Path Drawing
// 路径绘制

// 压缩
func (dc *Context) capper() raster.Capper {
	switch dc.lineCap {
	case LineCapButt:
		return raster.ButtCapper
	case LineCapRound:
		return raster.RoundCapper
	case LineCapSquare:
		return raster.SquareCapper
	}
	return nil
}
// 木匠
func (dc *Context) joiner() raster.Joiner {
	switch dc.lineJoin {
	case LineJoinBevel:
		return raster.BevelJoiner
	case LineJoinRound:
		return raster.RoundJoiner
	}
	return nil
}
// 打击
func (dc *Context) stroke(painter raster.Painter) {
	path := dc.strokePath
	if len(dc.dashes) > 0 {
		path = dashed(path, dc.dashes, dc.dashOffset)
	} else {
		// TODO: this is a temporary workaround to remove tiny segments
		// that result in rendering issues
		// TODO:这是一个临时解决方案，用于删除微小的片段
		// 这会导致渲染问题
		path = rasterPath(flattenPath(path))
	}
	r := dc.rasterizer
	r.UseNonZeroWinding = true
	r.Clear()
	r.AddStroke(path, fix(dc.lineWidth), dc.capper(), dc.joiner())
	r.Rasterize(painter)
}
// 填充
func (dc *Context) fill(painter raster.Painter) {
	path := dc.fillPath
	if dc.hasCurrent {
		path = make(raster.Path, len(dc.fillPath))
		copy(path, dc.fillPath)
		path.Add1(dc.start.Fixed())
	}
	r := dc.rasterizer
	r.UseNonZeroWinding = dc.fillRule == FillRuleWinding
	r.Clear()
	r.AddPath(path)
	r.Rasterize(painter)
}

// StrokePreserve strokes the current path with the current color, line width,
// line cap, line join and dash settings. The path is preserved after this
// operation.
//
// StrokePreserve 使用当前颜色、线宽、线帽、线连接和虚线设置描边当前路径。
// 此操作后将保留路径。
func (dc *Context) StrokePreserve() {
	var painter raster.Painter
	if dc.mask == nil {
		if pattern, ok := dc.strokePattern.(*solidPattern); ok {
			// with a nil mask and a solid color pattern, we can be more efficient
			// TODO: refactor so we don't have to do this type assertion stuff?
			p := raster.NewRGBAPainter(dc.im)
			p.SetColor(pattern.color)
			painter = p
		}
	}
	if painter == nil {
		painter = newPatternPainter(dc.im, dc.mask, dc.strokePattern)
	}
	dc.stroke(painter)
}

// Stroke strokes the current path with the current color, line width,
// line cap, line join and dash settings. The path is cleared after this
// operation.
//
// 使用当前颜色、线宽、线帽、线连接和虚线设置描边当前路径。 此操作后路径被清除。
func (dc *Context) Stroke() {
	dc.StrokePreserve()
	dc.ClearPath()
}

// FillPreserve fills the current path with the current color. Open subpaths
// are implicity closed. The path is preserved after this operation.
//
// 用当前颜色填充当前路径。 打开的子路径是隐式关闭的。 此操作后将保留路径。
func (dc *Context) FillPreserve() {
	var painter raster.Painter
	if dc.mask == nil {
		if pattern, ok := dc.fillPattern.(*solidPattern); ok {
			// with a nil mask and a solid color pattern, we can be more efficient
			// TODO: refactor so we don't have to do this type assertion stuff?
			// 使用 nil 掩码和纯色图案，我们可以更高效
			// TODO: 重构所以我们不必做这种类型断言的东西？
			p := raster.NewRGBAPainter(dc.im)
			p.SetColor(pattern.color)
			painter = p
		}
	}
	if painter == nil {
		painter = newPatternPainter(dc.im, dc.mask, dc.fillPattern)
	}
	dc.fill(painter)
}

// Fill fills the current path with the current color. Open subpaths
// are implicity closed. The path is cleared after this operation.
//
// 用当前颜色填充当前路径。 
// 打开的子路径是隐式关闭的。 此操作后路径被清除。
func (dc *Context) Fill() {
	dc.FillPreserve()
	dc.ClearPath()
}

// ClipPreserve updates the clipping region by intersecting the current
// clipping region with the current path as it would be filled by dc.Fill().
// The path is preserved after this operation.
//
// 通过将当前剪辑区域与当前路径相交来更新剪辑区域，因为它将由 dc.Fill() 填充。
// 此操作后将保留路径。
func (dc *Context) ClipPreserve() {
	clip := image.NewAlpha(image.Rect(0, 0, dc.width, dc.height))
	painter := raster.NewAlphaOverPainter(clip)
	dc.fill(painter)
	if dc.mask == nil {
		dc.mask = clip
	} else {
		mask := image.NewAlpha(image.Rect(0, 0, dc.width, dc.height))
		draw.DrawMask(mask, mask.Bounds(), clip, image.Point{}, dc.mask, image.Point{}, draw.Over)
		dc.mask = mask
	}
}

// SetMask allows you to directly set the *image.Alpha to be used as a clipping
// mask. It must be the same size as the context, else an error is returned
// and the mask is unchanged.
//
// 允许您直接设置 *image.Alpha 用作剪贴蒙版。
// 它必须与上下文大小相同，否则返回错误并且掩码不变。
func (dc *Context) SetMask(mask *image.Alpha) error {
	if mask.Bounds().Size() != dc.im.Bounds().Size() {
		return errors.New("mask size must match context size")
	}
	dc.mask = mask
	return nil
}

// AsMask returns an *image.Alpha representing the alpha channel of this
// context. This can be useful for advanced clipping operations where you first
// render the mask geometry and then use it as a mask.
//
// 返回一个 *image.Alpha 表示此上下文的 Alpha 通道。
// 这对于您首先渲染蒙版几何图形然后将其用作蒙版的高级裁剪操作很有用。
func (dc *Context) AsMask() *image.Alpha {
	mask := image.NewAlpha(dc.im.Bounds())
	draw.Draw(mask, dc.im.Bounds(), dc.im, image.Point{}, draw.Src)
	return mask
}

// InvertMask inverts the alpha values in the current clipping mask such that
// a fully transparent region becomes fully opaque and vice versa.
//
// 反转当前剪贴蒙版中的 alpha 值，使完全透明的区域变得完全不透明，反之亦然。
func (dc *Context) InvertMask() {
	if dc.mask == nil {
		dc.mask = image.NewAlpha(dc.im.Bounds())
	} else {
		for i, a := range dc.mask.Pix {
			dc.mask.Pix[i] = 255 - a
		}
	}
}

// Clip updates the clipping region by intersecting the current
// clipping region with the current path as it would be filled by dc.Fill().
// The path is cleared after this operation.
//
// 通过与当前的相交来更新剪辑区域
// 使用当前路径剪切区域，因为它将由 dc.Fill() 填充。
// 此操作后路径被清除。
func (dc *Context) Clip() {
	dc.ClipPreserve()
	dc.ClearPath()
}

// ResetClip clears the clipping region.
//
// 清除剪裁区域。
func (dc *Context) ResetClip() {
	dc.mask = nil
}

// Convenient Drawing Functions
// 方便的绘图功能

// Clear fills the entire image with the current color.
// 用当前颜色填充整个图像。
func (dc *Context) Clear() {
	src := image.NewUniform(dc.color)
	draw.Draw(dc.im, dc.im.Bounds(), src, image.Point{}, draw.Src)
}

// SetPixel sets the color of the specified pixel using the current color.
//
// 使用当前颜色设置指定像素的颜色。
func (dc *Context) SetPixel(x, y int) {
	dc.im.Set(x, y, dc.color)
}

// DrawPoint is like DrawCircle but ensures that a circle of the specified
// size is drawn regardless of the current transformation matrix. The position
// is still transformed, but not the shape of the point.
//
// DrawPoint 与 DrawCircle 类似，但确保绘制指定大小的圆，而不管当前的变换矩阵如何。
// 位置仍然被改变，但不是点的形状。
func (dc *Context) DrawPoint(x, y, r float64) {
	dc.Push()
	tx, ty := dc.TransformPoint(x, y)
	dc.Identity()
	dc.DrawCircle(tx, ty, r)
	dc.Pop()
}

// 绘制一条线
func (dc *Context) DrawLine(x1, y1, x2, y2 float64) {
	dc.MoveTo(x1, y1)
	dc.LineTo(x2, y2)
}

// 绘制矩形
func (dc *Context) DrawRectangle(x, y, w, h float64) {
	dc.NewSubPath()
	dc.MoveTo(x, y)
	dc.LineTo(x+w, y)
	dc.LineTo(x+w, y+h)
	dc.LineTo(x, y+h)
	dc.ClosePath()
}

// 绘制一个圆角矩形
func (dc *Context) DrawRoundedRectangle(x, y, w, h, r float64) {
	x0, x1, x2, x3 := x, x+r, x+w-r, x+w
	y0, y1, y2, y3 := y, y+r, y+h-r, y+h
	dc.NewSubPath()
	dc.MoveTo(x1, y0)
	dc.LineTo(x2, y0)
	dc.DrawArc(x2, y1, r, Radians(270), Radians(360))
	dc.LineTo(x3, y2)
	dc.DrawArc(x2, y2, r, Radians(0), Radians(90))
	dc.LineTo(x1, y3)
	dc.DrawArc(x1, y2, r, Radians(90), Radians(180))
	dc.LineTo(x0, y1)
	dc.DrawArc(x1, y1, r, Radians(180), Radians(270))
	dc.ClosePath()
}

// 绘制椭圆弧
func (dc *Context) DrawEllipticalArc(x, y, rx, ry, angle1, angle2 float64) {
	const n = 16
	for i := 0; i < n; i++ {
		p1 := float64(i+0) / n
		p2 := float64(i+1) / n
		a1 := angle1 + (angle2-angle1)*p1
		a2 := angle1 + (angle2-angle1)*p2
		x0 := x + rx*math.Cos(a1)
		y0 := y + ry*math.Sin(a1)
		x1 := x + rx*math.Cos((a1+a2)/2)
		y1 := y + ry*math.Sin((a1+a2)/2)
		x2 := x + rx*math.Cos(a2)
		y2 := y + ry*math.Sin(a2)
		cx := 2*x1 - x0/2 - x2/2
		cy := 2*y1 - y0/2 - y2/2
		if i == 0 {
			if dc.hasCurrent {
				dc.LineTo(x0, y0)
			} else {
				dc.MoveTo(x0, y0)
			}
		}
		dc.QuadraticTo(cx, cy, x2, y2)
	}
}

// 绘制椭圆
func (dc *Context) DrawEllipse(x, y, rx, ry float64) {
	dc.NewSubPath()
	dc.DrawEllipticalArc(x, y, rx, ry, 0, 2*math.Pi)
	dc.ClosePath()
}

// 绘制弧线
func (dc *Context) DrawArc(x, y, r, angle1, angle2 float64) {
	dc.DrawEllipticalArc(x, y, r, r, angle1, angle2)
}

// 绘制圆圈
func (dc *Context) DrawCircle(x, y, r float64) {
	dc.NewSubPath()
	dc.DrawEllipticalArc(x, y, r, r, 0, 2*math.Pi)
	dc.ClosePath()
}

// 绘制正多边形
func (dc *Context) DrawRegularPolygon(n int, x, y, r, rotation float64) {
	angle := 2 * math.Pi / float64(n)
	rotation -= math.Pi / 2
	if n%2 == 0 {
		rotation += angle / 2
	}
	dc.NewSubPath()
	for i := 0; i < n; i++ {
		a := rotation + angle*float64(i)
		dc.LineTo(x+r*math.Cos(a), y+r*math.Sin(a))
	}
	dc.ClosePath()
}

// DrawImage draws the specified image at the specified point.
//
// 在指定点绘制指定图像。
func (dc *Context) DrawImage(im image.Image, x, y int) {
	dc.DrawImageAnchored(im, x, y, 0, 0)
}

// DrawImageAnchored draws the specified image at the specified anchor point.
// The anchor point is x - w * ax, y - h * ay, where w, h is the size of the
// image. Use ax=0.5, ay=0.5 to center the image at the specified point.
//
// 在指定锚点处绘制指定图像。
// 锚点是x - w * ax, y - h * ay，其中w, h是图像的大小。
// 使用 ax=0.5, ay=0.5 使图像在指定点居中。
func (dc *Context) DrawImageAnchored(im image.Image, x, y int, ax, ay float64) {
	s := im.Bounds().Size()
	x -= int(ax * float64(s.X))
	y -= int(ay * float64(s.Y))
	transformer := dc.scaleStyle.transformer()
	fx, fy := float64(x), float64(y)
	m := dc.matrix.Translate(fx, fy)
	s2d := f64.Aff3{m.XX, m.XY, m.X0, m.YX, m.YY, m.Y0}
	if dc.mask == nil {
		transformer.Transform(dc.im, s2d, im, im.Bounds(), draw.Over, nil)
	} else {
		transformer.Transform(dc.im, s2d, im, im.Bounds(), draw.Over, &draw.Options{
			DstMask:  dc.mask,
			DstMaskP: image.Point{},
		})
	}
}

// Text Functions
// 文本函数

// 设置字体面
func (dc *Context) SetFontFace(fontFace font.Face) {
	dc.fontFace = fontFace
	dc.fontHeight = float64(fontFace.Metrics().Height) / 64
}

// Load the font from the specified path
//
// 加载指定路径的字体
func (dc *Context) LoadFontFace(path string, points float64) error {
	face, err := LoadFontFace(path, points)
	if err == nil {
		dc.fontFace = face
		dc.fontHeight = points * 72 / 96
	}
	return err
}

// 返回字体高度
func (dc *Context) FontHeight() float64 {
	return dc.fontHeight
}
// 绘制文本
func (dc *Context) drawString(im *image.RGBA, s string, x, y float64) {
	d := &font.Drawer{
		Dst:  im,
		Src:  image.NewUniform(dc.color),
		Face: dc.fontFace,
		Dot:  fixp(x, y),
	}
	// based on Drawer.DrawString() in golang.org/x/image/font/font.go
	prevC := rune(-1)
	for _, c := range s {
		if prevC >= 0 {
			d.Dot.X += d.Face.Kern(prevC, c)
		}
		dr, mask, maskp, advance, ok := d.Face.Glyph(d.Dot, c)
		if !ok {
			// TODO: is falling back on the U+FFFD glyph the responsibility of
			// the Drawer or the Face?
			// TODO: set prevC = '\ufffd'?
			continue
		}
		sr := dr.Sub(dr.Min)
		transformer := draw.BiLinear
		fx, fy := float64(dr.Min.X), float64(dr.Min.Y)
		m := dc.matrix.Translate(fx, fy)
		s2d := f64.Aff3{m.XX, m.XY, m.X0, m.YX, m.YY, m.Y0}
		transformer.Transform(d.Dst, s2d, d.Src, sr, draw.Over, &draw.Options{
			SrcMask:  mask,
			SrcMaskP: maskp,
		})
		d.Dot.X += advance
		prevC = c
	}
}

// DrawString draws the specified text at the specified point.
//
// 在指定点绘制指定文本。
func (dc *Context) DrawString(s string, x, y float64) {
	dc.DrawStringAnchored(s, x, y, 0, 0)
}

// DrawStringAnchored draws the specified text at the specified anchor point.
// The anchor point is x - w * ax, y - h * ay, where w, h is the size of the
// text. Use ax=0.5, ay=0.5 to center the text at the specified point.
//
// 在指定锚点处绘制指定文本。
// 锚点是 x - w * ax, y - h * ay，其中 w, h 是文本的大小。
// 使用 ax=0.5, ay=0.5 使文本在指定点居中。
func (dc *Context) DrawStringAnchored(s string, x, y, ax, ay float64) {
	w, h := dc.MeasureString(s)
	x -= ax * w
	y += ay * h
	if dc.mask == nil {
		dc.drawString(dc.im, s, x, y)
	} else {
		im := image.NewRGBA(image.Rect(0, 0, dc.width, dc.height))
		dc.drawString(im, s, x, y)
		draw.DrawMask(dc.im, dc.im.Bounds(), im, image.Point{}, dc.mask, image.Point{}, draw.Over)
	}
}

// DrawStringWrapped word-wraps the specified string to the given max width
// and then draws it at the specified anchor point using the given line
// spacing and text alignment.
//
// 将指定的字符串换行到给定的最大宽度，
// 然后使用给定的行距和文本对齐在指定的锚点处绘制它。
func (dc *Context) DrawStringWrapped(s string, x, y, ax, ay, width, lineSpacing float64, align Align) {
	lines := dc.WordWrap(s, width)

	// sync h formula with MeasureMultilineString
	// 同步 h 公式与 度量多行字符串
	h := float64(len(lines)) * dc.fontHeight * lineSpacing
	h -= (lineSpacing - 1) * dc.fontHeight

	x -= ax * width
	y -= ay * h
	switch align {
	case AlignLeft:
		ax = 0
	case AlignCenter:
		ax = 0.5
		x += width / 2
	case AlignRight:
		ax = 1
		x += width
	}
	ay = 1
	for _, line := range lines {
		dc.DrawStringAnchored(line, x, y, ax, ay)
		y += dc.fontHeight * lineSpacing
	}
}

// 具有度量多行字符串的公式
func (dc *Context) MeasureMultilineString(s string, lineSpacing float64) (width, height float64) {
	lines := strings.Split(s, "\n")

	// sync h formula with DrawStringWrapped
	height = float64(len(lines)) * dc.fontHeight * lineSpacing
	height -= (lineSpacing - 1) * dc.fontHeight

	d := &font.Drawer{
		Face: dc.fontFace,
	}

	// max width from lines
	for _, line := range lines {
		adv := d.MeasureString(line)
		currentWidth := float64(adv >> 6) // from gg.Context.MeasureString
		if currentWidth > width {
			width = currentWidth
		}
	}

	return width, height
}

// MeasureString returns the rendered width and height of the specified text
// given the current font face.
//
// 返回给定当前字体的指定文本的渲染宽度和高度。
func (dc *Context) MeasureString(s string) (w, h float64) {
	d := &font.Drawer{
		Face: dc.fontFace,
	}
	a := d.MeasureString(s)
	return float64(a >> 6), dc.fontHeight
}

// WordWrap wraps the specified string to the given max width and current
// font face.
//
// 将指定的字符串换行到给定的最大宽度和当前字体。
func (dc *Context) WordWrap(s string, w float64) []string {
	return wordWrap(dc, s, w)
}

// Transformation Matrix Operations
//矩阵运算的变换

// Identity resets the current transformation matrix to the identity matrix.
// This results in no translating, scaling, rotating, or shearing.
//
// 将当前变换矩阵重置为单位矩阵。
// 这不会导致平移、缩放、旋转或剪切。
func (dc *Context) Identity() {
	dc.matrix = Identity()
}

// Translate updates the current matrix with a translation.
//
// 使用平移更新当前矩阵。
func (dc *Context) Translate(x, y float64) {
	dc.matrix = dc.matrix.Translate(x, y)
}

// Scale updates the current matrix with a scaling factor.
// Scaling occurs about the origin.
//
// 使用缩放因子更新当前矩阵。缩放发生在原点附近。
func (dc *Context) Scale(x, y float64) {
	dc.matrix = dc.matrix.Scale(x, y)
}

// ScaleAbout updates the current matrix with a scaling factor.
// Scaling occurs about the specified point.
//
// 使用缩放因子更新当前矩阵。在指定点附近发生缩放。
func (dc *Context) ScaleAbout(sx, sy, x, y float64) {
	dc.Translate(x, y)
	dc.Scale(sx, sy)
	dc.Translate(-x, -y)
}

// Rotate updates the current matrix with a anticlockwise rotation.
// Rotation occurs about the origin. Angle is specified in radians.
//
// 逆时针旋转更新当前矩阵。围绕原点发生旋转。 角度以弧度指定。
func (dc *Context) Rotate(angle float64) {
	dc.matrix = dc.matrix.Rotate(angle)
}

// RotateAbout updates the current matrix with a anticlockwise rotation.
// Rotation occurs about the specified point. Angle is specified in radians.
//
// 逆时针旋转更新当前矩阵。围绕指定点进行旋转。 角度以弧度指定。
func (dc *Context) RotateAbout(angle, x, y float64) {
	dc.Translate(x, y)
	dc.Rotate(angle)
	dc.Translate(-x, -y)
}

// Shear updates the current matrix with a shearing angle.
// Shearing occurs about the origin.
//
// 用剪切角更新当前矩阵。剪切发生在原点附近。
func (dc *Context) Shear(x, y float64) {
	dc.matrix = dc.matrix.Shear(x, y)
}

// ShearAbout updates the current matrix with a shearing angle.
// Shearing occurs about the specified point.
//
// 用剪切角更新当前矩阵。剪切发生在指定点附近。
func (dc *Context) ShearAbout(sx, sy, x, y float64) {
	dc.Translate(x, y)
	dc.Shear(sx, sy)
	dc.Translate(-x, -y)
}

// TransformPoint multiplies the specified point by the current matrix,
// returning a transformed position.
//
// 将指定点乘以当前矩阵，返回转换后的位置。
func (dc *Context) TransformPoint(x, y float64) (tx, ty float64) {
	return dc.matrix.TransformPoint(x, y)
}

// InvertY flips the Y axis so that Y grows from bottom to top and Y=0 is at
// the bottom of the image.
//
// 反转Y轴，使Y从下到上增长，Y=0 位于图像的底部。
func (dc *Context) InvertY() {
	dc.Translate(0, float64(dc.height))
	dc.Scale(1, -1)
}

// Stack
// 堆栈

// Push saves the current state of the context for later retrieval. These
// can be nested.
//
// 保存上下文的当前状态以供以后检索。 这些可以嵌套。
func (dc *Context) Push() {
	x := *dc
	dc.stack = append(dc.stack, &x)
}

// Pop restores the last saved context state from the stack.
//
// 从堆栈中恢复上次保存的上下文状态。
func (dc *Context) Pop() {
	before := *dc
	s := dc.stack
	x, s := s[len(s)-1], s[:len(s)-1]
	*dc = *x
	dc.mask = before.mask
	dc.strokePath = before.strokePath
	dc.fillPath = before.fillPath
	dc.start = before.start
	dc.current = before.current
	dc.hasCurrent = before.hasCurrent
}
