package gg

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"unsafe"
)

// Predefined colors.
//
// 预定义颜色。
var (
	White   = color.RGBA{255, 255, 255, 255}
	Black   = color.RGBA{0, 0, 0, 255}
	Red     = color.RGBA{255, 0, 0, 255}
	Green   = color.RGBA{0, 255, 0, 255}
	Blue    = color.RGBA{0, 0, 255, 255}
	Yellow  = color.RGBA{255, 255, 0, 255}
	Cyan    = color.RGBA{0, 255, 255, 255}
	Magenta = color.RGBA{255, 0, 255, 255}
	Grey    = color.RGBA{190, 190, 190, 255}
	Pink    = color.RGBA{255, 181, 197, 255}
	Orange  = color.RGBA{255, 165, 0, 255}
	Opaque  = color.RGBA{0, 0, 0, 0}
)

// takeThemeColorsKMeans 实现基于k-means算法的图像取色算法
func takeThemeColorsKMeans(img image.Image, k int) []color.RGBA {
	rgbaimg := ImageToRGBA(img)
	pixels := unsafe.Slice(
		(*color.RGBA)(unsafe.Pointer(unsafe.SliceData(rgbaimg.Pix))),
		uintptr(len(rgbaimg.Pix))/unsafe.Sizeof(color.RGBA{}),
	)

	// 初始化k个聚类中心
	clusters := make([]color.RGBA, k)
	for i := range k {
		clusters[i] = pixels[rand.Intn(len(pixels))]
	}

	// 迭代聚类
	for {
		// 将每个像素点分配到最近的聚类中心
		clusterAssignments := make([]int, len(pixels))
		for i, pixel := range pixels {
			minDistance := math.MaxFloat64
			for j, cluster := range clusters {
				distance := distance(pixel, cluster)
				if distance < minDistance {
					minDistance = distance
					clusterAssignments[i] = j
				}
			}
		}

		// 计算每个聚类的新中心
		newClusters := make([]color.RGBA, k)
		for i := range k {
			var r, g, b uint32
			n := 0
			for j, cluster := range clusterAssignments {
				if cluster == i {
					pixel := pixels[j]
					r += uint32(pixel.R)
					g += uint32(pixel.G)
					b += uint32(pixel.B)
					n++
				}
			}
			if n != 0 {
				newClusters[i] = color.RGBA{uint8(r / uint32(n)), uint8(g / uint32(n)), uint8(b / uint32(n)), 255}
			}
		}

		// 如果聚类中心没有变化，则停止迭代
		if clustersEqual(clusters, newClusters) {
			break
		}
		clusters = newClusters
	}

	return clusters
}

// 计算两个颜色之间的距离
func distance(a, b color.RGBA) float64 {
	return math.Sqrt(sq(float64(a.R)-float64(b.R)) + sq(float64(a.G)-float64(b.G)) + sq(float64(a.B)-float64(b.B)))
}

// 计算平方
func sq(n float64) float64 {
	return n * n
}

// 比较两个聚类中心是否相等
func clustersEqual(a, b []color.RGBA) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
