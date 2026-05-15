package gg

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

	"github.com/FloatTech/gg/internal/gpu"
)

//go:generate ocloc compile -file internal/build/darkimg_ocl.cl -spv_only -options "-cl-mad-enable -cl-fast-relaxed-math -cl-finite-math-only -cl-single-precision-constant" -internal_options "-O3" -output internal/build/darkimg_ocl
//go:generate llvm-spirv -to-text internal/build/darkimg_ocl_.spv -o internal/build/darkimg_ocl.spt

//go:embed internal/build/darkimg_ocl_.spv
var darkimgspv []byte

var (
	canUseDarkimgKernel = false
	darkimgModel        ze.ModuleHandle
)

func init() {
	if !gpu.IsAvailable() {
		return
	}

	var err error
	darkimgModel, err = gpu.ModuleCreateAndCheckKernels(darkimgspv, "isdark")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[gg.darkimg_ocl] gpu init err:", err)
		return
	}

	canUseDarkimgKernel = true
}

func gpuIsDarkimg(img image.Image, scale float32) (bool, error) {
	rgbaimg := ImageToRGBA(img)
	pixels := rgbaimg.Pix
	srcW, srcH := img.Bounds().Dx(), img.Bounds().Dy()
	dstW, dstH := uint32(float32(srcW)*scale), uint32(float32(srcH)*scale)

	krn, err := darkimgModel.KernelCreate("isdark")
	if err != nil {
		return false, err
	}
	defer krn.Destroy()

	// Allocate input image buffer (host + device)
	srcSize := uintptr(srcW*srcH) * unsafe.Sizeof(color.RGBA{})
	inputImgHost, inputImgDevice, err := gpu.MemAllocHostDevicePair(srcSize, unsafe.Sizeof(color.RGBA{}))
	if err != nil {
		return false, err
	}
	defer gpu.MemFree(inputImgHost)
	defer gpu.MemFree(inputImgDevice)

	// Copy source pixels to host buffer
	himg := unsafe.Slice((*uint8)(inputImgHost), len(pixels))
	copy(himg, pixels)

	// Create input image handle
	inputImgHandle, err := gpu.ImageCreateUnorm(0, uint64(srcW), uint32(srcH))
	if err != nil {
		return false, err
	}
	defer inputImgHandle.Destroy()

	// Create sampler with user-specified filter
	smp, err := gpu.SamplerCreateNormalizedClamp(gozel.ZE_SAMPLER_FILTER_MODE_LINEAR)
	if err != nil {
		return false, err
	}
	defer smp.Destroy()

	// Allocate output buffer (host + device)
	dstSize := uintptr(dstW*dstH) * unsafe.Sizeof(uint8(0))
	outputImgHost, outputImgDevice, err := gpu.MemAllocHostDevicePair(dstSize, unsafe.Sizeof(uint8(0)))
	if err != nil {
		return false, err
	}
	defer gpu.MemFree(outputImgHost)
	defer gpu.MemFree(outputImgDevice)

	// Set kernel arguments
	err = krn.SetArgumentValue(0, inputImgHandle)
	if err != nil {
		return false, err
	}
	err = krn.SetArgumentValue(1, smp)
	if err != nil {
		return false, err
	}
	err = krn.SetArgumentValue(2, dstW)
	if err != nil {
		return false, err
	}
	err = krn.SetArgumentValue(3, dstH)
	if err != nil {
		return false, err
	}
	err = krn.SetArgumentValue(4, outputImgDevice)
	if err != nil {
		return false, err
	}

	// Determine group size
	gX, gY, _, err := krn.SuggestGroupSize(dstW, dstH, 1)
	if err != nil {
		return false, err
	}
	err = krn.SetGroupSize(gX, gY, 1)
	if err != nil {
		return false, err
	}
	gcx := uint32(math.Ceil(float64(dstW) / float64(gX)))
	gcy := uint32(math.Ceil(float64(dstH) / float64(gY)))

	// Build command list
	lst, err := gpu.CommandListCreate()
	if err != nil {
		return false, err
	}
	defer lst.Destroy()

	// Event: input image copy done
	inpcpev, cl, err := gpu.EventCreate(gozel.ZE_EVENT_SCOPE_FLAG_HOST, 0)
	if err != nil {
		return false, err
	}
	defer cl()
	defer inpcpev.Destroy()

	// Copy input image: host -> device -> image
	cl2, err := gpu.ImageCopyFromHostBuffer(
		lst, inputImgHost, inputImgDevice, srcSize,
		inputImgHandle, inpcpev,
	)
	if err != nil {
		return false, err
	}
	defer cl2()

	// Event: kernel done
	kev, cl3, err := gpu.EventCreate(gozel.ZE_EVENT_SCOPE_FLAG_HOST, 0)
	if err != nil {
		return false, err
	}
	defer cl3()
	defer kev.Destroy()

	// Launch kernel (wait for input copy)
	err = lst.AppendLaunchKernel(krn, &gozel.ZeGroupCount{
		Groupcountx: gcx, Groupcounty: gcy, Groupcountz: 1,
	}, kev, inpcpev)
	if err != nil {
		return false, err
	}

	// Event: output copy done
	outcpev, cl4, err := gpu.EventCreate(gozel.ZE_EVENT_SCOPE_FLAG_HOST, 0)
	if err != nil {
		return false, err
	}
	defer cl4()
	defer outcpev.Destroy()

	// Copy output data: image -> device -> host (wait for kernel)
	err = lst.AppendMemoryCopy(outputImgHost, outputImgDevice, dstSize, outcpev, kev)
	if err != nil {
		return false, err
	}

	// Close and execute
	err = lst.Close()
	if err != nil {
		return false, err
	}

	err = gpu.ExecCommandLists(lst)
	if err != nil {
		return false, err
	}

	// Wait for output copy to complete
	err = outcpev.HostSynchronize(math.MaxUint64)
	if err != nil {
		return false, err
	}

	// Build result from output buffer
	output := unsafe.Slice((*uint8)(outputImgHost), dstSize)
	visibleCount := 0
	for _, vis := range output {
		visibleCount += int(vis)
	}

	// 若不到 5% 的像素有肉眼可见细节，则认为几乎全黑
	return visibleCount*100/int(dstSize) < 5, nil
}
