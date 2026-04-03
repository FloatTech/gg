// Package gpu use gozel to run some heavy jobs on Intel GPUs.
package gpu

import (
	"errors"

	"github.com/fumiama/gozel/gozel"
	"github.com/fumiama/gozel/ze"
)

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
