package gpu

import (
	"unsafe"

	"github.com/fumiama/gozel/gozel"
	"github.com/fumiama/gozel/ze"
)

// ImageCreateUnorm create image that can load from go image.RGBA.Pix and is
// able to be used in linear interpolation.
func ImageCreateUnorm(flags gozel.ZeImageFlags, width uint64, height uint32) (ze.ImageHandle, error) {
	return ctx.ImageCreate(dev, flags, gozel.ZeImageFormat{
		Layout: gozel.ZE_IMAGE_FORMAT_LAYOUT_8_8_8_8,
		Type:   gozel.ZE_IMAGE_FORMAT_TYPE_UNORM, // UNORM: bilinear sampling returns float [0,1]
		X:      gozel.ZE_IMAGE_FORMAT_SWIZZLE_R,
		Y:      gozel.ZE_IMAGE_FORMAT_SWIZZLE_G,
		Z:      gozel.ZE_IMAGE_FORMAT_SWIZZLE_B,
		W:      gozel.ZE_IMAGE_FORMAT_SWIZZLE_A,
	}, width, height)
}

// ImageCopyFromHostBuffer wraps the calls that copies img data through host -> device -> image
func ImageCopyFromHostBuffer(
	lst ze.CommandListHandle, hbuf, dbuf unsafe.Pointer, sz uintptr,
	hDstImage ze.ImageHandle, hSignalEvent ze.EventHandle,
	waitEvents ...ze.EventHandle,
) (cl func(), err error) {
	evcph2did, err := evids.get()
	if err != nil {
		return
	}
	evcph2d, err := evph.EventCreate(evcph2did, gozel.ZE_EVENT_SCOPE_FLAG_HOST, 0)
	if err != nil {
		evids.put(evcph2did)
		return
	}
	err = lst.AppendMemoryCopy(dbuf, hbuf, sz, evcph2d, waitEvents...)
	if err != nil {
		_ = evcph2d.Destroy()
		evids.put(evcph2did)
		return
	}
	err = lst.AppendImageCopyFromMemory(hDstImage, dbuf, nil, hSignalEvent, evcph2d)
	if err != nil {
		_ = evcph2d.Destroy()
		evids.put(evcph2did)
		return
	}
	return func() {
		_ = evcph2d.Destroy()
		evids.put(evcph2did)
	}, nil
}

// ImageCopyToHostBuffer wraps the calls that copies img data through image -> device -> host
func ImageCopyToHostBuffer(
	lst ze.CommandListHandle, hbuf, dbuf unsafe.Pointer, sz uintptr,
	hSrcImage ze.ImageHandle, hSignalEvent ze.EventHandle,
	waitEvents ...ze.EventHandle,
) (cl func(), err error) {
	eid, err := evids.get()
	if err != nil {
		return
	}
	evcpim2d, err := evph.EventCreate(eid, gozel.ZE_EVENT_SCOPE_FLAG_HOST, 0)
	if err != nil {
		evids.put(eid)
		return
	}
	err = lst.AppendImageCopyToMemory(dbuf, hSrcImage, nil, evcpim2d, waitEvents...)
	if err != nil {
		_ = evcpim2d.Destroy()
		evids.put(eid)
		return
	}
	err = lst.AppendMemoryCopy(hbuf, dbuf, sz, hSignalEvent, evcpim2d)
	if err != nil {
		_ = evcpim2d.Destroy()
		evids.put(eid)
		return
	}
	return func() {
		_ = evcpim2d.Destroy()
		evids.put(eid)
	}, nil
}
