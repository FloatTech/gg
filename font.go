package gg

import (
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

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
	f, err := opentype.ParseCollection(fontBytes)
	if err != nil {
		return
	}
	fnf, err := f.Font(0)
	if err != nil {
		return
	}
	face, err = opentype.NewFace(fnf, &opentype.FaceOptions{
		Size: points,
		DPI:  72,
		// Hinting: font.HintingFull,
	})
	return
}

// ParseFontFace 是一个辅助函数，用于加载指定点大小的指定字体文件。
// 请注意，返回的 `font.Face` 对象不是线程安全的，不能跨 goroutine 并行使用。
// 您通常可以只使用 Context.LoadFontFace 函数而不是这个包级函数。
func ParseFontFace(b []byte, points float64) (face font.Face, err error) {
	f, err := opentype.ParseCollection(b)
	if err != nil {
		return
	}
	fnf, err := f.Font(0)
	if err != nil {
		return
	}
	face, err = opentype.NewFace(fnf, &opentype.FaceOptions{
		Size: points,
		DPI:  72,
		// Hinting: font.HintingFull,
	})
	return
}
