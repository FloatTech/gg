package total

import (
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"unsafe"

	"github.com/fumiama/gozel/gozel"
	"github.com/fumiama/gozel/ze"

	"github.com/FloatTech/gg/gpu"
)

//go:generate ocloc compile -file build/resize_ocl.cl -spv_only -options "-cl-mad-enable -cl-fast-relaxed-math -cl-finite-math-only -cl-single-precision-constant" -internal_options "-O3" -output build/resize_ocl
//go:generate llvm-spirv -to-text build/resize_ocl_.spv -o build/resize_ocl.spt

//go:embed build/resize_ocl_.spv
var resizespv []byte

// ResampleFilter is a GPU-supported resampling filter for image resizing.
//
// Intel GPU 硬件仅支持以下两种采样滤波模式：
//   - Nearest: 最近邻采样（无抗锯齿）
//   - Linear: 双线性插值采样
type ResampleFilter = gozel.ZeSamplerFilterMode

var (
	// ResampleFilterNearestNeighbor is a nearest-neighbor filter (no anti-aliasing).
	//
	// ResampleFilterNearestNeighbor 最近邻采样滤波器（无抗锯齿）。
	ResampleFilterNearestNeighbor ResampleFilter = gozel.ZE_SAMPLER_FILTER_MODE_NEAREST

	// ResampleFilterLinear is a bilinear interpolation filter.
	//
	// ResampleFilterLinear 双线性插值采样滤波器。
	ResampleFilterLinear ResampleFilter = gozel.ZE_SAMPLER_FILTER_MODE_LINEAR
)

var (
	canUseResizeKernel = false
	resizeModel        ze.ModuleHandle
)

func init() {
	if !gpu.IsAvailable() {
		return
	}

	var err error
	resizeModel, err = gpu.ModuleCreateAndCheckKernels(resizespv, "scale")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[gg.resize_ocl] gpu init err:", err)
		return
	}

	canUseResizeKernel = true
}

func gpuResize(img image.Image, srcW, srcH, dstW, dstH int, filter ResampleFilter) (*image.NRGBA, error) {
	rgbaimg := ImageToRGBA(img)
	pixels := rgbaimg.Pix

	krn, err := resizeModel.KernelCreate("scale")
	if err != nil {
		return nil, err
	}
	defer krn.Destroy()

	// Allocate input image buffer (host + device)
	srcSize := uintptr(srcW*srcH) * unsafe.Sizeof(color.RGBA{})
	inputImgHost, inputImgDevice, err := gpu.MemAllocHostDevicePair(srcSize, unsafe.Sizeof(color.RGBA{}))
	if err != nil {
		return nil, err
	}
	defer gpu.MemFree(inputImgHost)
	defer gpu.MemFree(inputImgDevice)

	// Copy source pixels to host buffer
	himg := unsafe.Slice((*uint8)(inputImgHost), len(pixels))
	copy(himg, pixels)

	// Create input image handle
	inputImgHandle, err := gpu.ImageCreateUnorm(0, uint64(srcW), uint32(srcH))
	if err != nil {
		return nil, err
	}
	defer inputImgHandle.Destroy()

	// Create sampler with user-specified filter
	smp, err := gpu.SamplerCreateNormalizedClamp(filter)
	if err != nil {
		return nil, err
	}
	defer smp.Destroy()

	// Allocate output image buffer (host + device)
	dstSize := uintptr(dstW*dstH) * unsafe.Sizeof(color.RGBA{})
	outputImgHost, outputImgDevice, err := gpu.MemAllocHostDevicePair(dstSize, unsafe.Sizeof(color.RGBA{}))
	if err != nil {
		return nil, err
	}
	defer gpu.MemFree(outputImgHost)
	defer gpu.MemFree(outputImgDevice)

	// Create output image handle (kernel-writable)
	outputImgHandle, err := gpu.ImageCreateUnorm(gozel.ZE_IMAGE_FLAG_KERNEL_WRITE, uint64(dstW), uint32(dstH))
	if err != nil {
		return nil, err
	}
	defer outputImgHandle.Destroy()

	// Set kernel arguments
	err = krn.SetArgumentValue(0, inputImgHandle)
	if err != nil {
		return nil, err
	}
	err = krn.SetArgumentValue(1, smp)
	if err != nil {
		return nil, err
	}
	err = krn.SetArgumentValue(2, outputImgHandle)
	if err != nil {
		return nil, err
	}

	// Determine group size
	gX, gY, _, err := krn.SuggestGroupSize(uint32(dstW), uint32(dstH), 1)
	if err != nil {
		return nil, err
	}
	err = krn.SetGroupSize(gX, gY, 1)
	if err != nil {
		return nil, err
	}
	gcx := uint32(math.Ceil(float64(dstW) / float64(gX)))
	gcy := uint32(math.Ceil(float64(dstH) / float64(gY)))

	// Build command list
	lst, err := gpu.CommandListCreate()
	if err != nil {
		return nil, err
	}
	defer lst.Destroy()

	// Event: input image copy done
	inpcpev, cl, err := gpu.EventCreate(gozel.ZE_EVENT_SCOPE_FLAG_HOST, 0)
	if err != nil {
		return nil, err
	}
	defer cl()
	defer inpcpev.Destroy()

	// Copy input image: host -> device -> image
	cl2, err := gpu.ImageCopyFromHostBuffer(
		lst, inputImgHost, inputImgDevice, srcSize,
		inputImgHandle, inpcpev,
	)
	if err != nil {
		return nil, err
	}
	defer cl2()

	// Event: kernel done
	kev, cl3, err := gpu.EventCreate(gozel.ZE_EVENT_SCOPE_FLAG_HOST, 0)
	if err != nil {
		return nil, err
	}
	defer cl3()
	defer kev.Destroy()

	// Launch kernel (wait for input copy)
	err = lst.AppendLaunchKernel(krn, &gozel.ZeGroupCount{
		Groupcountx: gcx, Groupcounty: gcy, Groupcountz: 1,
	}, kev, inpcpev)
	if err != nil {
		return nil, err
	}

	// Event: output copy done
	outcpev, cl4, err := gpu.EventCreate(gozel.ZE_EVENT_SCOPE_FLAG_HOST, 0)
	if err != nil {
		return nil, err
	}
	defer cl4()
	defer outcpev.Destroy()

	// Copy output image: image -> device -> host (wait for kernel)
	cl5, err := gpu.ImageCopyToHostBuffer(
		lst, outputImgHost, outputImgDevice, dstSize,
		outputImgHandle, outcpev, kev,
	)
	if err != nil {
		return nil, err
	}
	defer cl5()

	// Close and execute
	err = lst.Close()
	if err != nil {
		return nil, err
	}

	err = gpu.ExecCommandLists(lst)
	if err != nil {
		return nil, err
	}

	// Wait for output copy to complete
	err = outcpev.HostSynchronize(math.MaxUint64)
	if err != nil {
		return nil, err
	}

	// Build result NRGBA image from output buffer
	dst := image.NewNRGBA(image.Rect(0, 0, dstW, dstH))
	outPix := unsafe.Slice((*uint8)(outputImgHost), dstW*dstH*4)
	copy(dst.Pix, outPix)
	return dst, nil
}
