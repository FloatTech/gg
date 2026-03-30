package fio

import (
	"image/gif"
	"os"
)

// LoadGIF loads a GIF image from the specified file path.
//
// LoadGIF 从指定的文件路径加载 GIF 图像。
func LoadGIF(path string) (*gif.GIF, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return gif.DecodeAll(file)
}

// SaveGIF encodes a GIF and saves it to the specified file path.
//
// SaveGIF 保存 gif 到 path
func SaveGIF(path string, g *gif.GIF) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return gif.EncodeAll(f, g)
}
