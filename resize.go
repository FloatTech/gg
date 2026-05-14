package gg

import (
	"fmt"
	"image"
	"math"
	"os"

	"github.com/disintegration/imaging"
)

func cpuResize(img image.Image, dstW, dstH int, filter ResampleFilter) *image.NRGBA {
	var f imaging.ResampleFilter
	switch filter {
	case ResampleFilterNearestNeighbor:
		f = imaging.NearestNeighbor
	case ResampleFilterLinear:
		f = imaging.Linear
	default:
		f = imaging.Lanczos
	}
	return imaging.Resize(img, dstW, dstH, f)
}

// Resize resizes the image to the specified width and height using the specified
// resampling filter and returns the transformed image. If one of width or height
// is 0, the image aspect ratio is preserved.
//
// Resize 使用指定的重采样滤波器将图像缩放到指定的宽度和高度并返回变换后的图像。
// 如果 width 或 height 其中之一为 0，则保持图像宽高比。
func Resize(img image.Image, width, height int, filter ResampleFilter) *image.NRGBA {
	bounds := img.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()

	if width <= 0 && height <= 0 {
		width = srcW
		height = srcH
	} else if width <= 0 {
		width = int(math.Round(float64(height) * float64(srcW) / float64(srcH)))
	} else if height <= 0 {
		height = int(math.Round(float64(width) * float64(srcH) / float64(srcW)))
	}

	if width == srcW && height == srcH {
		return ImageToNRGBA(img)
	}

	if canUseResizeKernel {
		dst, err := gpuResize(img, srcW, srcH, width, height, filter)
		if err != nil {
			fmt.Fprintln(os.Stderr, "[gg.resize_ocl] gpuResize err:", err, "fallback to cpu")
			canUseResizeKernel = false
		} else {
			return dst
		}
	}

	// CPU fallback: simple nearest-neighbor or bilinear
	return cpuResize(img, width, height, filter)
}
