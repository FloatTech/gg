package total

import (
	"fmt"
	"image"
	"os"
)

func cpuIsDarkimg(im image.Image, scale float32) bool {
	img := Resize(
		im, int(float32(im.Bounds().Dx())*scale),
		int(float32(im.Bounds().Dy())*scale),
		ResampleFilterLinear,
	)

	bounds := img.Bounds()
	totalPixels := bounds.Dx() * bounds.Dy()
	if totalPixels == 0 {
		return true
	}

	const visibleThreshold = 15 // 亮度低于此值人眼几乎不可见
	visibleCount := 0

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		rowOffset := (y - bounds.Min.Y) * img.Stride
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			offset := rowOffset + (x-bounds.Min.X)*4
			r := uint32(img.Pix[offset])
			g := uint32(img.Pix[offset+1])
			b := uint32(img.Pix[offset+2])
			lum := (299*r + 587*g + 114*b) / 1000
			if lum > visibleThreshold {
				visibleCount++
			}
		}
	}

	// 若不到 5% 的像素有肉眼可见细节，则认为几乎全黑
	return visibleCount*100/totalPixels < 5
}

// IsDarkimg judges whether a given image can be visible to human eyes or,
// "not dark".
//
// IsDarkimg 判断图片是否为全黑或几乎全黑，以至于人眼不可辨识。
func IsDarkimg(im image.Image, scale float32) bool {
	if canUseDarkimgKernel {
		v, err := gpuIsDarkimg(im, scale)
		if err != nil {
			fmt.Fprintln(os.Stderr, "[gg.darkimg_ocl] gpuIsDarkimg err:", err, "fallback to cpu")
			canUseDarkimgKernel = false
		} else {
			return v
		}
	}

	return cpuIsDarkimg(im, scale)
}
