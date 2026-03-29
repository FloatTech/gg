package fio

import (
	"bufio"
	"image"
	"image/png"
	"os"
)

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
