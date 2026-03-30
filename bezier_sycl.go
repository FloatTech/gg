package gg

import (
	_ "embed"

	"github.com/FloatTech/gg/internal/gpu"
	"github.com/fumiama/gozel/ze"
)

//go:generate clang++ -fsycl -fsycl-device-only -fno-sycl-instrument-device-code -fsycl-targets=spirv64 -Xclang -emit-llvm-bc internal/build/bezier_sycl.cpp -o internal/build/device_bezier_kern.bc
//go:generate sycl-post-link -symbols -split=auto -emit-param-info -properties -o internal/build/device_bezier_kern.table internal/build/device_bezier_kern.bc
//go:generate llvm-spirv --sycl-opt -o internal/build/bezier_sycl.spv internal/build/device_bezier_kern_0.bc
//go:generate clang++ -target spirv64-unknown-unknown -S -emit-llvm -x ir internal/build/device_bezier_kern_0.bc -o main.ll
//go:generate llvm-spirv -to-text internal/build/bezier_sycl.spv -o internal/build/bezier_sycl.spt

//go:embed internal/build/bezier_sycl.spv
var bezierspv []byte

var (
	canUseBezierKernel = false
	bezierModel        ze.ModuleHandle
)

func init() {
	if !gpu.IsAvailable {
		return
	}

	var err error
	bezierModel, err = gpu.CreateModuleAndCheckKernels(bezierspv)
	if err != nil {
		return
	}

	canUseBezierKernel = true
}

func quadraticBezeirGPU(x0, y0, x1, y1, x2, y2, ds float64, p []Point) error {
	return gpu.Exec1D(bezierModel, "__sycl_kernel_quadratic", p,
		x0, y0, x1, y1, x2, y2, ds,
	)
}

func cubicBezeirGPU(x0, y0, x1, y1, x2, y2, x3, y3, ds float64, p []Point) error {
	return gpu.Exec1D(bezierModel, "__sycl_kernel_cubic", p,
		x0, y0, x1, y1, x2, y2, x3, y3, ds,
	)
}
