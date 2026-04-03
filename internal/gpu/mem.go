package gpu

import "unsafe"

// MemAllocHostDevicePair allocs both host and device mem with the same size.
func MemAllocHostDevicePair(size uintptr, alignment uintptr) (h, d unsafe.Pointer, err error) {
	h, err = ctx.MemAllocHost(size, alignment)
	if err != nil {
		return
	}
	d, err = ctx.MemAllocDevice(dev, size, alignment)
	if err != nil {
		_ = ctx.MemFree(h)
		return
	}
	return
}

// MemCopyGo2Host copies go array to host mem and returns the go array repr of dst.
func MemCopyGo2Host[T any](dst unsafe.Pointer, src []T) []T {
	dstSlice := unsafe.Slice((*T)(dst), len(src))
	copy(dstSlice, src)
	return dstSlice
}

// MemFree frees memory previously allocated with MemAllocDevice or MemAllocHost.
func MemFree(ptr unsafe.Pointer) error {
	return ctx.MemFree(ptr)
}
