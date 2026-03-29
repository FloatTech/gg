package fio

import (
	"bufio"
	"image"
	"image/jpeg"
	"os"
)

// LoadJPG loads a JPG image from the specified file path.
//
// LoadJPG 从指定的文件路径加载 JPG 图像。
func LoadJPG(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return jpeg.Decode(bufio.NewReader(file))
}

// SaveJPG encodes an image as JPG and saves it to the specified file path.
//
// SaveJPG 将图像编码为 JPG 并保存到指定的文件路径。
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
