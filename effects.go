package gg

func gateN100P100(per int) int {
	if per > 100 {
		per = 100
	}
	if per < -100 {
		per = -100
	}
	return per
}

func gate0P255(n int) int {
	if n < 0 {
		n = 0
	}
	if n > 255 {
		n = 255
	}
	return n
}

// Brightness 调整亮度 范围：±100%
func (dc *Context) Brightness(per int) {
	if per == 0 {
		return
	}
	per = gateN100P100(per)
	gain := 255 * per / 100
	for i, v := range dc.im.Pix {
		if i%4 == 3 { // alpha
			continue
		}
		dc.im.Pix[i] = uint8(gate0P255(int(v) + gain))
	}
}

// Contrast 调整对比度 范围：±100%
func (dc *Context) Contrast(per int) {
	if per == 0 {
		return
	}
	per = gateN100P100(per) + 100
	switch {
	case 0 <= per && per < 100: // 损益
		gain := per
		for i, v := range dc.im.Pix {
			if i%4 == 3 { // alpha
				continue
			}
			dc.im.Pix[i] = uint8(gate0P255(int(v) * gain / 100))
		}
	case 100 < per && per <= 200: // 增益
		gain := 200 - per
		for i, v := range dc.im.Pix {
			if i%4 == 3 { // alpha
				continue
			}
			dc.im.Pix[i] = uint8(gate0P255(int(v) * 100 / gain))
		}
	default:
		panic("unreachable")
	}
}
