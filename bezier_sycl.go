package gg

import (
	_ "embed"

	"github.com/fumiama/gozel/ze"
)

//go:generate clang++ -fsycl -fsycl-device-only -fsycl-targets=spirv64 -Xclang -emit-llvm-bc bezier_sycl.cpp -o device_bezier_kern.bc
//go:generate sycl-post-link -symbols -split=auto -o device_bezier_kern.table device_bezier_kern.bc
//go:generate llvm-spirv -o bezier_sycl.spv device_bezier_kern_0.bc
//go:generate clang++ -target spirv64-unknown-unknown -S -emit-llvm -x ir device_bezier_kern_0.bc -o bezier_sycl.ll

//go:embed bezier_sycl.spv
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
