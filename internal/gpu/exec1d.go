package gpu

import (
	"math"
	"unsafe"

	"github.com/fumiama/gozel/gozel"
	"github.com/fumiama/gozel/ze"
)

// Exec1D1Buf execute __sycl_kernel_name(&arg1, &arg2, ..., &p),
// p is the only buf passing to and reading back from GPU.
//
// Equivalent pseudo code:
//
//	for i := range p {
//		go gpuThread(arg1, arg2, ..., p, i)
//	}
func Exec1D1Buf[T any](mod ze.ModuleHandle, name string, p []T, args ...any) error {
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

	hbuf, dbuf, err := MemAllocHostDevicePair(sz, 1)
	if err != nil {
		return err
	}
	defer ctx.MemFree(hbuf)
	defer ctx.MemFree(dbuf)

	hp := MemCopyGo2Host(hbuf, p)

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

	evcph2d, cl, err := EventCreate(gozel.ZE_EVENT_SCOPE_FLAG_HOST, 0)
	if err != nil {
		return err
	}
	defer cl()
	defer evcph2d.Destroy()

	evk, cl, err := EventCreate(gozel.ZE_EVENT_SCOPE_FLAG_HOST, 0)
	if err != nil {
		return err
	}
	defer cl()
	defer evk.Destroy()

	evcpd2h, cl, err := EventCreate(gozel.ZE_EVENT_SCOPE_FLAG_HOST, 0)
	if err != nil {
		return err
	}
	defer cl()
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
