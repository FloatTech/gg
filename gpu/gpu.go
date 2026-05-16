// Package gpu use gozel to run some heavy jobs on Intel GPUs.
package gpu

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"sync/atomic"

	"github.com/fumiama/gozel/gozel"
	"github.com/fumiama/gozel/ze"
)

var (
	// ErrGPUIsBusy is returned when a worker cannot get a event ID
	ErrGPUIsBusy = errors.New("gpu is busy")
	// ErrEmptyGPUList when driver handles returns empty list
	ErrEmptyGPUList = errors.New("empty gpu list")
)

// defaultInstance is the default gpu instance.
var defaultInstance atomic.Pointer[instance]

func init() {
	ins, err := newInstance()
	if err != nil {
		fmt.Fprintln(os.Stderr, "[gg.gpu] init err:", err)
		return
	}
	defaultInstance.Store(ins)
}

// IsAvailable shows that a valid GPU is ready to use.
func IsAvailable() bool {
	return g() != nil
}

// Reset re-init GPU instance.
func Reset() {
	ins, err := newInstance()
	if err != nil {
		fmt.Fprintln(os.Stderr, "[gg.gpu] re-init err:", err)
		return
	}
	defaultInstance.Store(ins)
}

func g() *instance {
	return defaultInstance.Load()
}

// instance is for internal use.
type instance struct {
	dh    ze.DriverHandle
	ctx   ze.ContextHandle
	dev   ze.DeviceHandle
	dcp   gozel.ZeDeviceComputeProperties
	q     ze.CommandQueueHandle
	evids eventIDsTable
	evph  ze.EventPoolHandle
}

// newInstance init new GPU instance.
func newInstance() (inst *instance, err error) {
	ins := new(instance)
	gpus, err := ze.InitGPUDrivers()
	if err != nil {
		return
	}
	if len(gpus) == 0 {
		err = ErrEmptyGPUList
		return
	}
	ins.dh = gpus[0]

	ins.ctx, err = ins.dh.ContextCreate()
	if err != nil {
		ins.Destroy()
		return
	}

	devs, err := ins.dh.DeviceGet()
	if err != nil || len(devs) == 0 {
		ins.Destroy()
		return
	}
	ins.dev = devs[0]

	ins.dcp, err = ins.dev.DeviceGetComputeProperties()
	if err != nil {
		ins.Destroy()
		return
	}

	ins.q, err = ins.ctx.CommandQueueCreate(ins.dev, gozel.ZE_COMMAND_QUEUE_MODE_ASYNCHRONOUS)
	if err != nil {
		ins.Destroy()
		return
	}

	ins.evph, err = ins.ctx.EventPoolCreate(gpuEventSize, ins.dev)
	if err != nil {
		ins.Destroy()
		return
	}

	runtime.SetFinalizer(ins, func(ins *instance) {
		ins.Destroy()
	})

	return ins, nil
}

// Destroy GPU instance.
func (ins *instance) Destroy() {
	if ins.evph != 0 {
		_ = ins.evph.Destroy()
	}
	if ins.q != 0 {
		_ = ins.q.Destroy()
	}
	if ins.ctx != 0 {
		_ = ins.ctx.Destroy()
	}
}
