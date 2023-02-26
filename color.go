package gg

import (
	"image"
	"image/color"
	"math"
	"math/rand"
)

// takecolor 实现基于k-means算法的图像取色算法
func takecolor(img image.Image, k int) []color.RGBA {
	// 获取图像的宽高
	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	// 获取图像的像素点
	var pixels []color.RGBA
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			r, g, b, _ := img.At(x, y).RGBA()
			pixels = append(pixels, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 255})
		}
	}

	// 初始化k个聚类中心
	clusters := make([]color.RGBA, k)
	for i := 0; i < k; i++ {
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
		for i := 0; i < k; i++ {
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
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
