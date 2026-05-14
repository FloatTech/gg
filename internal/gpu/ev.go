package gpu

import (
	"github.com/fumiama/gozel/gozel"
	"github.com/fumiama/gozel/ze"
)

const gpuEventSize = 1024

// EventCreate automatically get an empty event id and create event.
func EventCreate(signal gozel.ZeEventScopeFlags, wait gozel.ZeEventScopeFlags) (ze.EventHandle, func(), error) {
	eid, err := g().evids.get()
	if err != nil {
		return 0, nil, err
	}
	ev, err := g().evph.EventCreate(eid, signal, wait)
	if err != nil {
		g().evids.put(eid)
		return 0, nil, err
	}
	return ev, func() { g().evids.put(eid) }, nil
}
