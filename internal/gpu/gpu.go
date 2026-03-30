// Package gpu use gozel to run some heavy jobs on Intel GPUs.
package gpu

import (
	"errors"
	"math"
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
	dh    ze.DriverHandle
	ctx   ze.ContextHandle
	dev   ze.DeviceHandle
	dcp   gozel.ZeDeviceComputeProperties
	q     ze.CommandQueueHandle
	evids eventIDsTable
	evph  ze.EventPoolHandle
)

// IsAvailable shows that GPU is available for calling.
var IsAvailable = func() bool {
	gpus, err := ze.InitGPUDrivers()
	if err != nil || len(gpus) == 0 {
		return false
	}
	dh = gpus[0]

	ctx, err = dh.ContextCreate()
	if err != nil {
		Destroy()
		return false
	}

	devs, err := dh.DeviceGet()
	if err != nil || len(devs) == 0 {
		Destroy()
		return false
	}
	dev = devs[0]

	dcp, err = dev.DeviceGetComputeProperties()
	if err != nil {
		Destroy()
		return false
	}

	q, err = ctx.CommandQueueCreate(dev, gozel.ZE_COMMAND_QUEUE_MODE_ASYNCHRONOUS)
	if err != nil {
		Destroy()
		return false
	}

	evph, err = ctx.EventPoolCreate(gpuEventSize, dev)
	if err != nil {
		Destroy()
		return false
	}

	return true
}()

// Destroy GPU instance.
func Destroy() {
	if evph != 0 {
		_ = evph.Destroy()
	}
	if q != 0 {
		_ = q.Destroy()
	}
	if ctx != 0 {
		_ = ctx.Destroy()
	}
}

// CreateModuleAndCheckKernels loads module from spv and check kernel names' exisitance.
func CreateModuleAndCheckKernels(spv []byte, names ...string) (ze.ModuleHandle, error) {
	mod, err := ctx.ModuleCreate(dev, spv)
	if err != nil {
		return 0, err
	}
	for _, name := range names {
		krn, err := mod.KernelCreate(name)
		if err != nil {
			_ = mod.Destroy()
			return 0, err
		}
		_ = krn.Destroy()
	}
	return mod, nil
}

// Exec1D execute __sycl_kernel_name(&arg1, &arg2, ..., &p)
func Exec1D[T any](mod ze.ModuleHandle, name string, p []T, args ...any) error {
	if len(p) == 0 {
		return nil
	}

	n := uintptr(len(p))
	sz := n * unsafe.Sizeof(p[0])
	gsz := min(dcp.Maxgroupsizex, dcp.Maxtotalgroupsize)
	gc := sz / uintptr(gsz)
	if sz%uintptr(gsz) != 0 {
		gc++
	}

	hbuf, err := ctx.MemAllocHost(sz, 1)
	if err != nil {
		return err
	}
	defer ctx.MemFree(hbuf)

	dbuf, err := ctx.MemAllocDevice(dev, sz, 1)
	if err != nil {
		return err
	}
	defer ctx.MemFree(dbuf)

	hp := unsafe.Slice((*T)(hbuf), n)
	copy(hp, p)

	krn, err := mod.KernelCreate(name)
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

	eid, err := evids.get()
	if err != nil {
		return err
	}
	defer evids.put(eid)
	evcph2d, err := evph.EventCreate(eid, gozel.ZE_EVENT_SCOPE_FLAG_HOST, 0)
	if err != nil {
		return err
	}
	defer evcph2d.Destroy()

	eid, err = evids.get()
	if err != nil {
		return err
	}
	defer evids.put(eid)
	evk, err := evph.EventCreate(eid, gozel.ZE_EVENT_SCOPE_FLAG_HOST, 0)
	if err != nil {
		return err
	}
	defer evk.Destroy()

	eid, err = evids.get()
	if err != nil {
		return err
	}
	defer evids.put(eid)
	evcpd2h, err := evph.EventCreate(eid, gozel.ZE_EVENT_SCOPE_FLAG_HOST, 0)
	if err != nil {
		return err
	}
	defer evcpd2h.Destroy()

	lst, err := ctx.CommandListCreate(dev)
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

	err = q.ExecuteCommandLists(lst)
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
