package gpu

import (
	"github.com/fumiama/gozel/gozel"
	"github.com/fumiama/gozel/ze"
)

const gpuEventSize = 1024

// EventCreate automatically get an empty event id and create event.
func EventCreate(signal gozel.ZeEventScopeFlags, wait gozel.ZeEventScopeFlags) (ze.EventHandle, func(), error) {
	eid, err := evids.get()
	if err != nil {
		return 0, nil, err
	}
	ev, err := evph.EventCreate(eid, signal, wait)
	if err != nil {
		evids.put(eid)
		return 0, nil, err
	}
	return ev, func() { evids.put(eid) }, nil
}
