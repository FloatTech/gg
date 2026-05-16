package gpu

import "sync/atomic"

type eventIDsTable [gpuEventSize]atomic.Bool

func (evids *eventIDsTable) get() (uint32, error) {
	for i := range gpuEventSize {
		if (&evids[i]).CompareAndSwap(false, true) {
			return uint32(i), nil
		}
	}
	return 0, ErrGPUIsBusy
}

func (evids *eventIDsTable) put(id uint32) {
	(&evids[id]).Store(false)
}
