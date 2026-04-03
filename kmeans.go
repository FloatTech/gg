package gg

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"unsafe"

	"github.com/disintegration/imaging"
	"github.com/fumiama/gozel/ze"
)

type kmeansImage struct {
	pixels                []color.RGBA
	rs, gs, bs, ns        []uint32
	clusters, newClusters []color.RGBA
	clusterAssignments    []uint16

	canUseGPU                bool
	isRemainingAssign        bool
	gcx, gcy                 uint32
	bounds                   image.Rectangle
	krn1st, krnrem           ze.KernelHandle
	inputImgHost             unsafe.Pointer
	inputImgDevice           unsafe.Pointer
	inputImgHandle           ze.ImageHandle
	smp                      ze.SamplerHandle
	clustersHost             unsafe.Pointer // clustersHost is the underlying buf of clusters
	clustersDevice           unsafe.Pointer
	clustersImgHandle        ze.ImageHandle
	clusterAssignmentsHost   unsafe.Pointer // clusterAssignmentsHost is the underlying buf of clusterAssignments
	clusterAssignmentsDevice unsafe.Pointer
	sampleResultHost         unsafe.Pointer // sampleResultHost is the underlying buf of pixels
	sampleResultDevice       unsafe.Pointer
	sampleResult             ze.ImageHandle
}

func newKMeansImage(img image.Image, k uint16) kmeansImage {
	rgbaimg := ImageToRGBA(img)
	pixels := unsafe.Slice(
		(*color.RGBA)(unsafe.Pointer(unsafe.SliceData(rgbaimg.Pix))),
		uintptr(len(rgbaimg.Pix))/unsafe.Sizeof(color.RGBA{}),
	)
	clusters := make([]color.RGBA, k)
	for i := range k {
		clusters[i] = pixels[rand.Intn(len(pixels))]
	}
	ki := kmeansImage{
		pixels: pixels,
		rs:     make([]uint32, k), gs: make([]uint32, k), bs: make([]uint32, k), ns: make([]uint32, k),
		clusters:           clusters,
		clusterAssignments: make([]uint16, len(pixels)),

		bounds: img.Bounds(),
	}
	if canUseKmeansKernel {
		ki.gpuInit()
	}
	if !ki.canUseGPU {
		width := img.Bounds().Dx()
		height := img.Bounds().Dy()
		dstw, dsth := width, height
		ratio := 0.
		if dstw > 512 {
			dstw = 512
			ratio = float64(dstw) / float64(width)
			dsth *= int(float64(height) * ratio)
		}
		if dsth > 512 {
			dsth = 512
			ratio = float64(dsth) / float64(height)
			dstw = int(float64(width) * ratio)
		}
		ki.bounds = image.Rect(0, 0, dstw, dsth)
		rgbaimg = (*image.RGBA)(imaging.Resize(img, dstw, dsth, imaging.Lanczos))
		pixels = unsafe.Slice(
			(*color.RGBA)(unsafe.Pointer(unsafe.SliceData(rgbaimg.Pix))),
			uintptr(len(rgbaimg.Pix))/unsafe.Sizeof(color.RGBA{}),
		)
		ki.pixels = pixels
		ki.clusterAssignments = ki.clusterAssignments[:len(pixels)]
	}
	return ki
}

// assign 将每个像素点分配到最近的聚类中心
func (ki *kmeansImage) assign() {
	if ki.canUseGPU {
		err := ki.gpuAssign()
		if err == nil {
			return
		}
		ki.gpuDestroy(true)
	}

	for i, pixel := range ki.pixels {
		minDistance := math.MaxFloat64
		assign := uint16(math.MaxUint16)
		for j, cluster := range ki.clusters {
			distance := distanceRGBAsq(pixel, cluster)
			if distance < minDistance {
				minDistance = distance
				assign = uint16(j)
			}
		}
		ki.clusterAssignments[i] = assign
	}
}

// update 计算每个聚类的新中心
func (ki *kmeansImage) update() {
	for i, pixelCluster := range ki.clusterAssignments {
		if pixelCluster == uint16(math.MaxUint16) {
			continue
		}
		pixel := ki.pixels[i]
		ki.rs[pixelCluster] += uint32(pixel.R)
		ki.gs[pixelCluster] += uint32(pixel.G)
		ki.bs[pixelCluster] += uint32(pixel.B)
		ki.ns[pixelCluster]++
	}
}

func (ki *kmeansImage) epilogue() bool {
	if ki.newClusters == nil {
		ki.newClusters = make([]color.RGBA, len(ki.clusters))
	}
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
	if ki.canUseGPU {
		c := make([]color.RGBA, len(ki.clusters))
		copy(c, ki.clusters)
		return c
	}
	return ki.clusters
}

func (ki *kmeansImage) destroy() {
	if ki.canUseGPU {
		ki.gpuDestroy(false)
	}
}
