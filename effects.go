package gg

// Brightness 调整亮度 范围：±100%
func (dc *Context) Brightness(per int) {
	if per == 0 {
		return
	}
	per = clamp(per, -100, 100)
	gain := 255 * per / 100
	for i, v := range dc.im.Pix {
		if i%4 == 3 { // alpha
			continue
		}
		dc.im.Pix[i] = uint8(clamp(int(v)+gain, 0, 255))
	}
}

// Contrast 调整对比度 范围：±100%
func (dc *Context) Contrast(per int) {
	if per == 0 {
		return
	}
	per = clamp(per, -100, 100) + 100
	switch {
	case 0 <= per && per < 100: // 损益
		gain := per
		for i, v := range dc.im.Pix {
			if i%4 == 3 { // alpha
				continue
			}
			dc.im.Pix[i] = uint8(clamp(int(v)*gain/100, 0, 255))
		}
	case 100 < per && per <= 200: // 增益
		gain := 200 - per
		if gain == 0 {
			gain = 1
		}
		for i, v := range dc.im.Pix {
			if i%4 == 3 { // alpha
				continue
			}
			dc.im.Pix[i] = uint8(clamp(int(v)*100/gain, 0, 255))
		}
	default:
		panic("unreachable")
	}
}
