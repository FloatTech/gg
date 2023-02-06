package gg

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math"
	"net/http"
	"os"
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
	im, _, err := image.Decode(bufio.NewReader(file))
	return im, err
}

// 加载指定路径的 JPG 图像
func LoadJPG(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return jpeg.Decode(bufio.NewReader(file))
}

// 保存 JPG 图像
func SaveJPG(path string, im image.Image, quality int) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return jpeg.Encode(file, im, &jpeg.Options{
		Quality: quality, // 质量百分比
	})
}

// 加载指定路径的 PNG 图像
func LoadPNG(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return png.Decode(bufio.NewReader(file))
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

// image.Image 转为 image.RGBA64
func ImageToRGBA64(src image.Image) *image.RGBA64 {
	bounds := src.Bounds()
	dst := image.NewRGBA64(bounds)
	draw.Draw(dst, bounds, src, bounds.Min, draw.Src)
	return dst
}

// image.Image 转为 image.NRGBA
func ImageToNRGBA(src image.Image) *image.NRGBA {
	bounds := src.Bounds()
	dst := image.NewNRGBA(bounds)
	draw.Draw(dst, bounds, src, bounds.Min, draw.Src)
	return dst
}

// image.Image 转为 *image.NRGBA64
func ImageToNRGBA64(src image.Image) *image.NRGBA64 {
	bounds := src.Bounds()
	dst := image.NewNRGBA64(bounds)
	draw.Draw(dst, bounds, src, bounds.Min, draw.Src)
	return dst
}

// 解析图片的宽高信息
func GetWH(path string) (int, int, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return 0, 0, err
	}
	return GetImgWH(b)
}

// 解析图片的宽高信息
func GetImgWH(imgBytes []byte) (int, int, error) {
	var width, height int
	t := http.DetectContentType(imgBytes) // 获取文件类型
	switch t {
	case "image/jpg", "image/jpeg":
		width, height = GetJPGWH(imgBytes)
	case "image/png":
		width, height = GetPNGWH(imgBytes)
	case "image/gif":
		width, height = GetGIFWH(imgBytes)
	default:
		return 0, 0, errors.New("错误的文件类型")
	}
	return width, height, nil
}

// 获取 JPG 图片的宽高
func GetJPGWH(img []byte) (int, int) {
	var index int
	iL := len(img)
	for i := 0; i < iL-1; i++ {
		if img[i] != 0xff {
			continue
		}
		if img[i+1] == 0xC0 || img[i+1] == 0xC1 || img[i+1] == 0xC2 {
			index = i
			break
		}
	}
	index += 5
	if index >= iL {
		return 0, 0
	}
	height := int(img[index])<<8 + int(img[index+1])
	width := int(img[index+2])<<8 + int(img[index+3])
	return width, height
}

// 获取 PNG 图片的宽高
func GetPNGWH(imgBytes []byte) (int, int) {
	pngHeader := "\x89PNG\r\n\x1a\n"
	if string(imgBytes[:len(pngHeader)]) != pngHeader {
		return 0, 0
	}
	index := 12
	if string(imgBytes[index:index+4]) != "IHDR" {
		return 0, 0
	}
	index += 4
	width := int(binary.BigEndian.Uint32(imgBytes[index : index+4]))
	height := int(binary.BigEndian.Uint32(imgBytes[index+4 : index+8]))
	return width, height
}

// 获取 GIF 图片的宽高
func GetGIFWH(imgBytes []byte) (int, int) {
	ver := string(imgBytes[:6])
	if ver != "GIF87a" && ver != "GIF89a" {
		return 0, 0
	}
	width := int(imgBytes[6]) + int(imgBytes[7])<<8
	height := int(imgBytes[8]) + int(imgBytes[9])<<8
	return width, height
}

// 解析十六进制颜色
func ParseHexColor(x string) (r, g, b, a int) {
	return parseHexColor(x)
}

// 解析十六进制颜色
func parseHexColor(x string) (r, g, b, a int) {
	x = strings.TrimPrefix(x, "#")
	a = 255
	switch len(x) {
	case 3:
		format := "%1x%1x%1x"
		fmt.Sscanf(x, format, &r, &g, &b)
		r |= r << 4
		g |= g << 4
		b |= b << 4
	case 6:
		format := "%02x%02x%02x"
		fmt.Sscanf(x, format, &r, &g, &b)
	case 8:
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
func LoadFontFace(path string, points float64) (face font.Face, err error) {
	fontBytes, err := os.ReadFile(path)
	if err != nil {
		return
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		return
	}
	face = truetype.NewFace(f, &truetype.Options{
		Size: points,
		// Hinting: font.HintingFull,
	})
	return
}

// ParseFontFace 是一个辅助函数，用于加载指定点大小的指定字体文件。
// 请注意，返回的 `font.Face` 对象不是线程安全的，不能跨 goroutine 并行使用。
// 您通常可以只使用 Context.LoadFontFace 函数而不是这个包级函数。
func ParseFontFace(b []byte, points float64) (face font.Face, err error) {
	f, err := truetype.Parse(b)
	if err != nil {
		return
	}
	face = truetype.NewFace(f, &truetype.Options{
		Size: points,
		// Hinting: font.HintingFull,
	})
	return
}
