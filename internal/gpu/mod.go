package gpu

import "github.com/fumiama/gozel/ze"

// ModuleCreateAndCheckKernels loads module from spv and check kernel names' exisitance.
func ModuleCreateAndCheckKernels(spv []byte, names ...string) (ze.ModuleHandle, error) {
	mod, err := ctx.ModuleCreate(dev, spv)
	if err != nil {
		return 0, err
	}
	for _, name := range names {
		krn, err := mod.KernelCreate(name)
		if err != nil {
			_ = mod.Destroy()
			return 0, err
		}
		_ = krn.Destroy()
	}
	return mod, nil
}
