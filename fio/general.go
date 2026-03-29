package fio

import (
	"bufio"
	"image"
	"os"
)

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
