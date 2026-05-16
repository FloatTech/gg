package gpu

import "unsafe"

// MemAllocHostDevicePair allocs both host and device mem with the same size.
func MemAllocHostDevicePair(size uintptr, alignment uintptr) (h, d unsafe.Pointer, err error) {
	h, err = g().ctx.MemAllocHost(size, alignment)
	if err != nil {
		return
	}
	d, err = g().ctx.MemAllocDevice(g().dev, size, alignment)
	if err != nil {
		_ = g().ctx.MemFree(h)
		return
	}
	return
}

// MemAllocHost allocs host mem.
func MemAllocHost(size uintptr, alignment uintptr) (unsafe.Pointer, error) {
	return g().ctx.MemAllocHost(size, alignment)
}

// MemAllocDevice allocs device mem.
func MemAllocDevice(size uintptr, alignment uintptr) (unsafe.Pointer, error) {
	return g().ctx.MemAllocDevice(g().dev, size, alignment)
}

// MemCopyGo2Host copies go array to host mem and returns the go array repr of dst.
func MemCopyGo2Host[T any](dst unsafe.Pointer, src []T) []T {
	dstSlice := unsafe.Slice((*T)(dst), len(src))
	copy(dstSlice, src)
	return dstSlice
}

// MemFree frees memory previously allocated with MemAllocDevice or MemAllocHost.
func MemFree(ptr unsafe.Pointer) error {
	return g().ctx.MemFree(ptr)
}
