package gg

import (
	_ "embed"

	"github.com/fumiama/gozel/ze"
)

//go:generate clang++ -fsycl -fsycl-device-only -fno-sycl-instrument-device-code -fsycl-targets=spirv64 -Xclang -emit-llvm-bc build/bezier_sycl.cpp -o build/device_bezier_kern.bc
//go:generate sycl-post-link -symbols -split=auto -emit-param-info -properties -o build/device_bezier_kern.table build/device_bezier_kern.bc
//go:generate llvm-spirv --sycl-opt -o build/bezier_sycl.spv build/device_bezier_kern_0.bc
//go:generate clang++ -target spirv64-unknown-unknown -S -emit-llvm -x ir build/device_bezier_kern_0.bc -o main.ll
//go:generate llvm-spirv -to-text build/bezier_sycl.spv -o build/bezier_sycl.spt

//go:embed build/bezier_sycl.spv
var bezierspv []byte

var (
	canUseBezierKernel = false
	bezierModel        ze.ModuleHandle
)

func init() {
	if !canUseGPU {
		return
	}

	var err error
	bezierModel, err = gpuCreateModWithKernels(bezierspv)
	if err != nil {
		return
	}

	canUseBezierKernel = true
}

func quadraticBezeirGPU(x0, y0, x1, y1, x2, y2, ds float64, p []Point) error {
	return gpuExec1DKernelWithArgs("__sycl_kernel_quadratic", p,
		x0, y0, x1, y1, x2, y2, ds,
	)
}

func cubicBezeirGPU(x0, y0, x1, y1, x2, y2, x3, y3, ds float64, p []Point) error {
	return gpuExec1DKernelWithArgs("__sycl_kernel_cubic", p,
		x0, y0, x1, y1, x2, y2, x3, y3, ds,
	)
}
