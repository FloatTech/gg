package gg

import (
	"errors"
	"math"
	"sync/atomic"
	"unsafe"

	"github.com/fumiama/gozel/gozel"
	"github.com/fumiama/gozel/ze"
)

const gpuEventSize = 32

var (
	// ErrGPUIsBusy is returned when a worker cannot get a event ID
	ErrGPUIsBusy = errors.New("gpu is busy")
)

var (
	gpu      ze.DriverHandle
	gpuCtx   ze.ContextHandle
	gpuDev   ze.DeviceHandle
	gpuCProp gozel.ZeDeviceComputeProperties
	gpuQ     ze.CommandQueueHandle
	gpuEvIDs [gpuEventSize]atomic.Bool
	gpuEvP   ze.EventPoolHandle

	canUseGPU = func() bool {
		gpus, err := ze.InitGPUDrivers()
		if err != nil || len(gpus) == 0 {
			return false
		}
		gpu = gpus[0]

		gpuCtx, err = gpu.ContextCreate()
		if err != nil {
			gpuDestroyAll()
			return false
		}

		devs, err := gpu.DeviceGet()
		if err != nil || len(devs) == 0 {
			gpuDestroyAll()
			return false
		}
		gpuDev = devs[0]

		gpuCProp, err = gpuDev.DeviceGetComputeProperties()
		if err != nil {
			gpuDestroyAll()
			return false
		}

		gpuQ, err = gpuCtx.CommandQueueCreate(gpuDev, gozel.ZE_COMMAND_QUEUE_MODE_ASYNCHRONOUS)
		if err != nil {
			gpuDestroyAll()
			return false
		}

		gpuEvP, err = gpuCtx.EventPoolCreate(gpuEventSize, gpuDev)
		if err != nil {
			gpuDestroyAll()
			return false
		}

		return true
	}()
)

func gpuDestroyAll() {
	if gpuEvP != 0 {
		_ = gpuEvP.Destroy()
	}
	if gpuQ != 0 {
		_ = gpuQ.Destroy()
	}
	if gpuCtx != 0 {
		_ = gpuCtx.Destroy()
	}
}

func gpuCreateModWithKernels(spv []byte, names ...string) (ze.ModuleHandle, error) {
	mod, err := gpuCtx.ModuleCreate(gpuDev, spv)
	if err != nil {
		return 0, err
	}
	for _, name := range names {
		krn, err := mod.KernelCreate(name)
		if err != nil {
			mod.Destroy()
			return 0, err
		}
		_ = krn.Destroy()
	}
	return mod, nil
}

func gpuLockEventID() (uint32, error) {
	for i := range gpuEventSize {
		if (&gpuEvIDs[i]).CompareAndSwap(false, true) {
			return uint32(i), nil
		}
	}
	return 0, ErrGPUIsBusy
}

func gpuUnlockEventID(id uint32) {
	(&gpuEvIDs[id]).Store(false)
}

// gpuExec1DKernelWithArgs execure __sycl_kernel_name(&arg1, &arg2, ..., &p)
func gpuExec1DKernelWithArgs[T any](name string, p []T, args ...any) error {
	if len(p) == 0 {
		return nil
	}

	n := uintptr(len(p))
	sz := n * unsafe.Sizeof(p[0])
	gsz := min(gpuCProp.Maxgroupsizex, gpuCProp.Maxtotalgroupsize)
	gc := sz / uintptr(gsz)
	if sz%uintptr(gsz) != 0 {
		gc++
	}

	hbuf, err := gpuCtx.MemAllocHost(sz, 1)
	if err != nil {
		return err
	}
	defer gpuCtx.MemFree(hbuf)

	dbuf, err := gpuCtx.MemAllocDevice(gpuDev, sz, 1)
	if err != nil {
		return err
	}
	defer gpuCtx.MemFree(dbuf)

	hp := unsafe.Slice((*T)(hbuf), n)
	copy(hp, p)

	krn, err := bezierModel.KernelCreate(name)
	if err != nil {
		return err
	}
	defer krn.Destroy()

	for i, arg := range args {
		err = krn.SetArgumentValue(uint32(i), arg)
		if err != nil {
			return err
		}
	}
	err = krn.SetArgumentValue(uint32(len(args)), &dbuf)
	if err != nil {
		return err
	}

	err = krn.SetGroupSize(gsz, 1, 1)
	if err != nil {
		return err
	}

	eid, err := gpuLockEventID()
	if err != nil {
		return err
	}
	defer gpuUnlockEventID(eid)
	evcph2d, err := gpuEvP.EventCreate(eid, gozel.ZE_EVENT_SCOPE_FLAG_HOST, 0)
	if err != nil {
		return err
	}
	defer evcph2d.Destroy()

	eid, err = gpuLockEventID()
	if err != nil {
		return err
	}
	defer gpuUnlockEventID(eid)
	evk, err := gpuEvP.EventCreate(eid, gozel.ZE_EVENT_SCOPE_FLAG_HOST, 0)
	if err != nil {
		return err
	}
	defer evk.Destroy()

	eid, err = gpuLockEventID()
	if err != nil {
		return err
	}
	defer gpuUnlockEventID(eid)
	evcpd2h, err := gpuEvP.EventCreate(eid, gozel.ZE_EVENT_SCOPE_FLAG_HOST, 0)
	if err != nil {
		return err
	}
	defer evcpd2h.Destroy()

	lst, err := gpuCtx.CommandListCreate(gpuDev)
	if err != nil {
		return err
	}
	defer lst.Destroy()
	err = lst.AppendMemoryCopy(dbuf, hbuf, sz, evcph2d)
	if err != nil {
		return err
	}
	err = lst.AppendLaunchKernel(krn, &gozel.ZeGroupCount{
		Groupcountx: uint32(gc), Groupcounty: 1, Groupcountz: 1,
	}, evk, evcph2d)
	if err != nil {
		return err
	}
	err = lst.AppendMemoryCopy(hbuf, dbuf, sz, evcpd2h, evk)
	if err != nil {
		return err
	}
	err = lst.Close()
	if err != nil {
		return err
	}

	err = gpuQ.ExecuteCommandLists(lst)
	if err != nil {
		return err
	}

	err = evcpd2h.HostSynchronize(math.MaxUint64)
	if err != nil {
		return err
	}

	copy(p, hp)

	return nil
}
