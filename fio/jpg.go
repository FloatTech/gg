package fio

import (
	"bufio"
	"image"
	"image/jpeg"
	"os"
)

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
