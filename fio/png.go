package fio

import (
	"bufio"
	"image"
	"image/png"
	"os"
)

// LoadPNG loads a PNG image from the specified file path.
//
// LoadPNG 从指定的文件路径加载 PNG 图像。
func LoadPNG(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return png.Decode(bufio.NewReader(file))
}

// SavePNG encodes an image as PNG and saves it to the specified file path.
//
// SavePNG 将图像编码为 PNG 并保存到指定的文件路径。
func SavePNG(path string, im image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, im)
}
