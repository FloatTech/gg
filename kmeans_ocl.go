package gg

import (
	_ "embed"
	"image/color"
	"math"
	"unsafe"

	"github.com/FloatTech/gg/internal/gpu"
	"github.com/fumiama/gozel/gozel"
	"github.com/fumiama/gozel/ze"
)

//go:generate ocloc compile -file internal/build/kmeans_ocl.cl -spv_only -options "-cl-mad-enable -cl-fast-relaxed-math -cl-finite-math-only -cl-single-precision-constant" -internal_options "-O3" -output internal/build/kmeans_ocl
//go:generate llvm-spirv -to-text internal/build/kmeans_ocl_.spv -o internal/build/kmeans_ocl.spt

//go:embed internal/build/kmeans_ocl_.spv
var kmeansspv []byte

var (
	canUseKmeansKernel = false
	kmeansModel        ze.ModuleHandle
)

func init() {
	if !gpu.IsAvailable {
		return
	}

	var err error
	kmeansModel, err = gpu.ModuleCreateAndCheckKernels(kmeansspv)
	if err != nil {
		return
	}

	canUseKmeansKernel = true
}

func (ki *kmeansImage) gpuInit() error {
	width := ki.bounds.Dx()
	height := ki.bounds.Dy()
	ki.bounds = ImageBoundsBelow(ki.bounds, 512, 512)
	dstw, dsth := ki.bounds.Dx(), ki.bounds.Dy()

	krn1st, err := kmeansModel.KernelCreate("assign_first_iter")
	if err != nil {
		canUseKmeansKernel = false
		return err
	}
	krnrem, err := kmeansModel.KernelCreate("assign_remaining_iter")
	if err != nil {
		canUseKmeansKernel = false
		_ = krn1st.Destroy()
		return err
	}
	ki.krn1st = krn1st
	ki.krnrem = krnrem

	inputImgHost, inputImgDevice, err := gpu.MemAllocHostDevicePair(
		uintptr(len(ki.pixels))*unsafe.Sizeof(color.RGBA{}),
		unsafe.Sizeof(color.RGBA{}),
	)
	if err != nil {
		canUseKmeansKernel = false
		ki.gpuDestroy(true)
		return err
	}
	_ = gpu.MemCopyGo2Host(inputImgHost, ki.pixels)
	ki.inputImgHost = inputImgHost
	ki.inputImgDevice = inputImgDevice

	inputImgHandle, err := gpu.ImageCreateUnorm(0, uint64(width), uint32(height))
	if err != nil {
		canUseKmeansKernel = false
		ki.gpuDestroy(true)
		return err
	}
	ki.inputImgHandle = inputImgHandle

	smp, err := gpu.SamplerCreateNormalizedLinearClamp()
	if err != nil {
		canUseKmeansKernel = false
		ki.gpuDestroy(true)
		return err
	}
	ki.smp = smp

	clustersHost, clustersDevice, err := gpu.MemAllocHostDevicePair(
		uintptr(len(ki.clusters))*unsafe.Sizeof(color.RGBA{}),
		unsafe.Sizeof(color.RGBA{}),
	)
	if err != nil {
		canUseKmeansKernel = false
		ki.gpuDestroy(true)
		return err
	}
	ki.clusters = gpu.MemCopyGo2Host(clustersHost, ki.clusters)
	ki.clustersHost = clustersHost
	ki.clustersDevice = clustersDevice

	clustersImgHandle, err := gpu.ImageCreateUnorm(
		gozel.ZE_IMAGE_FLAG_KERNEL_WRITE, uint64(len(ki.clusters)), 1,
	)
	if err != nil {
		canUseKmeansKernel = false
		ki.gpuDestroy(true)
		return err
	}
	ki.clustersImgHandle = clustersImgHandle

	clusterAssignmentsHost, clusterAssignmentsDevice, err := gpu.MemAllocHostDevicePair(
		uintptr(dstw*dsth)*unsafe.Sizeof(uint16(0)),
		unsafe.Sizeof(uint16(0)),
	)
	if err != nil {
		canUseKmeansKernel = false
		ki.gpuDestroy(true)
		return err
	}
	ki.clusterAssignments = unsafe.Slice((*uint16)(clusterAssignmentsHost), dstw*dsth)
	ki.clusterAssignmentsHost = clusterAssignmentsHost
	ki.clusterAssignmentsDevice = clusterAssignmentsDevice

	sampleResultHost, sampleResultDevice, err := gpu.MemAllocHostDevicePair(
		uintptr(dstw*dsth)*unsafe.Sizeof(color.RGBA{}),
		unsafe.Sizeof(color.RGBA{}),
	)
	if err != nil {
		canUseKmeansKernel = false
		ki.gpuDestroy(true)
		return err
	}
	ki.sampleResultHost = sampleResultHost
	ki.sampleResultDevice = sampleResultDevice

	sampleResult, err := gpu.ImageCreateUnorm(
		gozel.ZE_IMAGE_FLAG_KERNEL_WRITE, uint64(dstw), uint32(dsth),
	)
	if err != nil {
		canUseKmeansKernel = false
		ki.gpuDestroy(true)
		return err
	}
	ki.sampleResult = sampleResult

	err = krn1st.SetArgumentValue(0, ki.inputImgHandle)
	if err != nil {
		canUseKmeansKernel = false
		ki.gpuDestroy(true)
		return err
	}
	err = krn1st.SetArgumentValue(1, ki.smp)
	if err != nil {
		canUseKmeansKernel = false
		ki.gpuDestroy(true)
		return err
	}
	err = krn1st.SetArgumentValue(2, ki.clustersImgHandle)
	if err != nil {
		canUseKmeansKernel = false
		ki.gpuDestroy(true)
		return err
	}
	err = krn1st.SetArgumentValue(3, &ki.clusterAssignmentsDevice)
	if err != nil {
		canUseKmeansKernel = false
		ki.gpuDestroy(true)
		return err
	}
	err = krn1st.SetArgumentValue(4, ki.sampleResult)
	if err != nil {
		canUseKmeansKernel = false
		ki.gpuDestroy(true)
		return err
	}

	err = krnrem.SetArgumentValue(0, ki.sampleResult)
	if err != nil {
		canUseKmeansKernel = false
		ki.gpuDestroy(true)
		return err
	}
	err = krnrem.SetArgumentValue(1, ki.clustersImgHandle)
	if err != nil {
		canUseKmeansKernel = false
		ki.gpuDestroy(true)
		return err
	}
	err = krnrem.SetArgumentValue(2, &ki.clusterAssignmentsDevice)
	if err != nil {
		canUseKmeansKernel = false
		ki.gpuDestroy(true)
		return err
	}

	gX, gY, _, err := krn1st.SuggestGroupSize(uint32(dstw), uint32(dsth), 1)
	if err != nil {
		canUseKmeansKernel = false
		ki.gpuDestroy(true)
		return err
	}
	err = krn1st.SetGroupSize(gX, gY, 1)
	if err != nil {
		canUseKmeansKernel = false
		ki.gpuDestroy(true)
		return err
	}
	err = krnrem.SetGroupSize(gX, gY, 1)
	if err != nil {
		canUseKmeansKernel = false
		ki.gpuDestroy(true)
		return err
	}
	ki.gcx = uint32(math.Ceil(float64(dstw) / float64(gX)))
	ki.gcy = uint32(math.Ceil(float64(dsth) / float64(gY)))

	ki.canUseGPU = true
	return nil
}

func (ki *kmeansImage) gpuAssign() error {
	var (
		srcN     = uintptr(len(ki.pixels))
		srcbufsz = srcN * unsafe.Sizeof(color.RGBA{})
	)

	lst, err := gpu.CommandListCreate()
	if err != nil {
		return err
	}
	defer lst.Destroy()

	kev, cl, err := gpu.EventCreate(gozel.ZE_EVENT_SCOPE_FLAG_HOST, 0)
	if err != nil {
		return err
	}
	defer cl()
	defer kev.Destroy()

	clacpev, cl, err := gpu.EventCreate(gozel.ZE_EVENT_SCOPE_FLAG_HOST, 0)
	if err != nil {
		return err
	}
	defer cl()
	defer clacpev.Destroy()

	waitev := ze.EventHandle(0)

	cluscpev, cl, err := gpu.EventCreate(gozel.ZE_EVENT_SCOPE_FLAG_HOST, 0)
	if err != nil {
		return err
	}
	defer cl()
	defer cluscpev.Destroy()

	cl, err = gpu.ImageCopyFromHostBuffer(
		lst, ki.clustersHost, ki.clustersDevice,
		uintptr(len(ki.clusters))*unsafe.Sizeof(color.RGBA{}),
		ki.clustersImgHandle, cluscpev,
	)
	if err != nil {
		return err
	}
	defer cl()

	if !ki.isRemainingAssign {
		ki.isRemainingAssign = true

		inpcpev, cl, err := gpu.EventCreate(gozel.ZE_EVENT_SCOPE_FLAG_HOST, 0)
		if err != nil {
			return err
		}
		defer cl()
		defer inpcpev.Destroy()

		cl, err = gpu.ImageCopyFromHostBuffer(
			lst, ki.inputImgHost, ki.inputImgDevice,
			srcbufsz, ki.inputImgHandle, inpcpev,
		)
		if err != nil {
			return err
		}
		defer cl()

		err = lst.AppendLaunchKernel(ki.krn1st, &gozel.ZeGroupCount{
			Groupcountx: ki.gcx, Groupcounty: ki.gcy, Groupcountz: 1,
		}, kev, inpcpev, cluscpev)
		if err != nil {
			return err
		}

		smpcpev, cl, err := gpu.EventCreate(gozel.ZE_EVENT_SCOPE_FLAG_HOST, 0)
		if err != nil {
			return err
		}
		defer cl()
		defer smpcpev.Destroy()

		cl, err = gpu.ImageCopyToHostBuffer(
			lst, ki.sampleResultHost, ki.sampleResultDevice,
			uintptr(len(ki.clusterAssignments))*unsafe.Sizeof(color.RGBA{}),
			ki.sampleResult, smpcpev, kev,
		)
		if err != nil {
			return err
		}
		defer cl()

		waitev = smpcpev
	} else {
		err = lst.AppendLaunchKernel(ki.krnrem, &gozel.ZeGroupCount{
			Groupcountx: ki.gcx, Groupcounty: ki.gcy, Groupcountz: 1,
		}, kev, cluscpev)
	}

	if err != nil {
		return err
	}

	cl, err = gpu.ImageCopyToHostBuffer(
		lst, ki.clustersHost, ki.clustersDevice,
		uintptr(len(ki.clusters))*unsafe.Sizeof(color.RGBA{}),
		ki.clustersImgHandle, cluscpev, kev,
	)
	if err != nil {
		return err
	}
	defer cl()

	err = lst.AppendMemoryCopy(
		ki.clusterAssignmentsHost, ki.clusterAssignmentsDevice,
		uintptr(len(ki.clusterAssignments))*unsafe.Sizeof(uint16(0)),
		clacpev, kev,
	)
	if err != nil {
		return err
	}

	err = lst.Close()
	if err != nil {
		return err
	}

	err = gpu.ExecCommandLists(lst)
	if err != nil {
		return err
	}

	if waitev != 0 {
		err = waitev.HostSynchronize(math.MaxUint64)
		if err != nil {
			return err
		}
		ki.pixels = unsafe.Slice(
			(*color.RGBA)(ki.sampleResultHost),
			uintptr(len(ki.clusterAssignments)),
		)
	}
	err = cluscpev.HostSynchronize(math.MaxUint64)
	if err != nil {
		return err
	}
	return clacpev.HostSynchronize(math.MaxUint64)
}

func (ki *kmeansImage) gpuDestroy(isFailed bool) {
	ki.canUseGPU = false

	// Copy data back to Go-managed slices before freeing host memory,
	// so the CPU fallback path won't access freed memory.
	if isFailed && ki.clustersHost != nil {
		goSlice := make([]color.RGBA, len(ki.clusters))
		copy(goSlice, ki.clusters)
		ki.clusters = goSlice
	}
	if isFailed && ki.clusterAssignmentsHost != nil {
		goSlice := make([]uint16, len(ki.clusterAssignments))
		copy(goSlice, ki.clusterAssignments)
		ki.clusterAssignments = goSlice
	}
	if isFailed && ki.sampleResultHost != nil {
		goSlice := make([]color.RGBA, len(ki.pixels))
		copy(goSlice, ki.pixels)
		ki.pixels = goSlice
	}

	if ki.inputImgHost != nil {
		_ = gpu.MemFree(ki.inputImgHost)
	}
	if ki.inputImgDevice != nil {
		_ = gpu.MemFree(ki.inputImgDevice)
	}
	if ki.inputImgHandle != 0 {
		_ = ki.inputImgHandle.Destroy()
	}
	if ki.smp != 0 {
		_ = ki.smp.Destroy()
	}
	if ki.clustersHost != nil {
		_ = gpu.MemFree(ki.clustersHost)
	}
	if ki.clustersDevice != nil {
		_ = gpu.MemFree(ki.clustersDevice)
	}
	if ki.clustersImgHandle != 0 {
		_ = ki.clustersImgHandle.Destroy()
	}
	if ki.clusterAssignmentsHost != nil {
		_ = gpu.MemFree(ki.clusterAssignmentsHost)
	}
	if ki.clusterAssignmentsDevice != nil {
		_ = gpu.MemFree(ki.clusterAssignmentsDevice)
	}
	if ki.sampleResultHost != nil {
		_ = gpu.MemFree(ki.sampleResultHost)
	}
	if ki.sampleResultDevice != nil {
		_ = gpu.MemFree(ki.sampleResultDevice)
	}
	if ki.sampleResult != 0 {
		_ = ki.sampleResult.Destroy()
	}
	if ki.krn1st != 0 {
		_ = ki.krn1st.Destroy()
	}
	if ki.krnrem != 0 {
		_ = ki.krnrem.Destroy()
	}
}
