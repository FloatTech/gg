// Package fio provides file I/O utilities for loading and saving images.
//
// fio 包提供图像加载和保存的文件 I/O 工具。
package fio

import (
	"bufio"
	"image"
	"os"
)

// LoadImage loads an image from the specified file path.
//
// LoadImage 从指定的文件路径加载图像。
func LoadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	im, _, err := image.Decode(bufio.NewReader(file))
	return im, err
}
