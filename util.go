package gg

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/freetype/truetype"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// 弧度
func Radians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// 角度
func Degrees(radians float64) float64 {
	return radians * 180 / math.Pi
}

// 加载指定路径的图像
func LoadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	im, _, err := image.Decode(file)
	return im, err
}

// 加载指定路径的 PNG 图像
func LoadPNG(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return png.Decode(file)
}

// 保存 PNG 图像
func SavePNG(path string, im image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, im)
}

// 加载指定路径的 JPG 图像
func LoadJPG(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return jpeg.Decode(file)
}

// 保存 JPG 图像
func SaveJPG(path string, im image.Image, quality int) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	var opt jpeg.Options
	opt.Quality = quality // 质量数值
	return jpeg.Encode(file, im, &opt)
}

// image.Image 转为 image.RGBA
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

// 解析十六进制颜色
func ParseHexColor(x string) (r, g, b, a int) {
	return parseHexColor(x)
}

// 解析十六进制颜色
func parseHexColor(x string) (r, g, b, a int) {
	x = strings.TrimPrefix(x, "#")
	a = 255
	if len(x) == 3 {
		format := "%1x%1x%1x"
		fmt.Sscanf(x, format, &r, &g, &b)
		r |= r << 4
		g |= g << 4
		b |= b << 4
	}
	if len(x) == 6 {
		format := "%02x%02x%02x"
		fmt.Sscanf(x, format, &r, &g, &b)
	}
	if len(x) == 8 {
		format := "%02x%02x%02x%02x"
		fmt.Sscanf(x, format, &r, &g, &b, &a)
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
func LoadFontFace(path string, points float64) (font.Face, error) {
	fontBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}
	face := truetype.NewFace(f, &truetype.Options{
		Size: points,
		// Hinting: font.HintingFull,
	})
	return face, nil
}

//  解析图片的宽高信息
func GetWH(path string) (float64, float64, error) {
	b, err := ioutil.ReadFile(path)
	if err !=nil {
	return 0, 0, err
	}
	return GetImgWH(b, path)
}

//  解析图片的宽高信息
func GetImgWH(imgBytes []byte, path string) (float64, float64, error) {
	var (
		imgConf image.Config
		err     error
	)
	switch strings.ToLower(filepath.Ext(path)) {
	case ".jpg", ".jpeg":
		imgConf, err = jpeg.DecodeConfig(bytes.NewReader(imgBytes))
	case ".png":
		imgConf, err = png.DecodeConfig(bytes.NewReader(imgBytes))
	case ".gif":
		imgConf, err = gif.DecodeConfig(bytes.NewReader(imgBytes))
	default:
		return 0, 0, errors.New("错误的文件类型")
	}
	if err != nil {
		return 0, 0, err
	}
	return float64(imgConf.Width), float64(imgConf.Height), nil
}

