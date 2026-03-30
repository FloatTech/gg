package gg

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"unsafe"
)

type kmeansImage struct {
	pixels                []color.RGBA
	rs, gs, bs, ns        []uint64
	clusters, newClusters []color.RGBA
	clusterAssignments    []int
}

func newKMeansImage(img image.Image, k int) kmeansImage {
	rgbaimg := ImageToRGBA(img)
	pixels := unsafe.Slice(
		(*color.RGBA)(unsafe.Pointer(unsafe.SliceData(rgbaimg.Pix))),
		uintptr(len(rgbaimg.Pix))/unsafe.Sizeof(color.RGBA{}),
	)
	clusters := make([]color.RGBA, k)
	for i := range k {
		clusters[i] = pixels[rand.Intn(len(pixels))]
	}
	return kmeansImage{
		pixels: pixels,
		rs:     make([]uint64, k), gs: make([]uint64, k), bs: make([]uint64, k), ns: make([]uint64, k),
		clusters: clusters, newClusters: make([]color.RGBA, k),
		clusterAssignments: make([]int, len(pixels)),
	}
}

// assign 将每个像素点分配到最近的聚类中心
func (ki *kmeansImage) assign() {
	for i, pixel := range ki.pixels {
		minDistance := math.MaxFloat64
		for j, cluster := range ki.clusters {
			distance := distanceRGBAsq(pixel, cluster)
			if distance < minDistance {
				minDistance = distance
				ki.clusterAssignments[i] = j
			}
		}
	}
}

// update 计算每个聚类的新中心
func (ki *kmeansImage) update() {
	for i, pixelCluster := range ki.clusterAssignments {
		for currentCluster := range ki.rs {
			if pixelCluster == currentCluster {
				pixel := ki.pixels[i]
				ki.rs[currentCluster] += uint64(pixel.R)
				ki.gs[currentCluster] += uint64(pixel.G)
				ki.bs[currentCluster] += uint64(pixel.B)
				ki.ns[currentCluster]++
			}
		}
	}
}

func (ki *kmeansImage) epilogue() bool {
	for i, n := range ki.ns {
		if n == 0 {
			ki.newClusters[i] = ki.clusters[i]
		} else {
			ki.newClusters[i] = color.RGBA{
				uint8(ki.rs[i] / n),
				uint8(ki.gs[i] / n),
				uint8(ki.bs[i] / n),
				255,
			}
		}
	}

	if isArrayRGBAEqual(ki.clusters, ki.newClusters) {
		return true
	}

	clear(ki.rs)
	clear(ki.gs)
	clear(ki.bs)
	clear(ki.ns)

	copy(ki.clusters, ki.newClusters)

	return false
}

func (ki *kmeansImage) result() []color.RGBA {
	return ki.clusters
}
