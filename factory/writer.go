package factory

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
)

// SaveGIF2Path 保存 gif 到 path
func SaveGIF2Path(path string, g *gif.GIF) error {
	f, err := os.Create(path) // 创建文件
	if err == nil {
		_ = gif.EncodeAll(f, g) // 写入
		_ = f.Close()           // 关闭文件
	}
	return err
}

// SavePNG2Path 保存 png 到 path
func SavePNG2Path(path string, im image.Image) error {
	f, err := os.Create(path) // 创建文件
	if err == nil {
		err = png.Encode(f, im) // 写入
		_ = f.Close()
	}
	return err
}

// ToBase64 img 内容转为base64
func ToBase64(img image.Image) (base64Bytes []byte, err error) {
	buffer := bytes.NewBuffer(make([]byte, 0, 1024*1024*4)) // 4MB
	encoder := base64.NewEncoder(base64.StdEncoding, buffer)
	var opt jpeg.Options
	opt.Quality = 70
	if err := jpeg.Encode(encoder, img, &opt); err != nil {
		return nil, err
	}
	_ = encoder.Close()
	base64Bytes = buffer.Bytes()
	return
}

// ToBytes img 内容转为 []byte
func ToBytes(img image.Image) (data []byte, err error) {
	buffer := bytes.NewBuffer(make([]byte, 0, 1024*1024*4)) // 4MB
	err = jpeg.Encode(buffer, img, &jpeg.Options{Quality: 70})
	data = buffer.Bytes()
	return
}

// WriteTo img 内容写入 Writer
func WriteTo(img image.Image, f io.Writer) (n int64, err error) {
	data, err := ToBytes(img)
	if err != nil {
		return
	}
	c, err := f.Write(data)
	return int64(c), err
}

// GIF2Base64 gif 内容转为 base64
func GIF2Base64(gifImage *gif.GIF) (string, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 1024*1024*4)) // 4MB
	err := gif.EncodeAll(buf, gifImage)
	if err != nil {
		return "", err
	}
	encodedGIF := base64.StdEncoding.EncodeToString(buf.Bytes())
	return "base64://" + encodedGIF, nil
}
