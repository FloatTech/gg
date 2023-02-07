package gg

func limitper(per int) int {
	if per > 100 {
		per = 100
	}
	if per < -100 {
		per = -100
	}
	return per
}

func limit2uint8(n int) int {
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
	per = limitper(per)
	gain := 255 * per / 100
	for i, v := range dc.im.Pix {
		if i%4 == 3 { // alpha
			continue
		}
		dc.im.Pix[i] = uint8(limit2uint8(int(v) + gain))
	}
}

// Contrast 调整对比度 范围：±100%
func (dc *Context) Contrast(per int) {
	if per == 0 {
		return
	}
	per = limitper(per) + 100
	gain := 0
	switch {
	case 0 <= per && per <= 100: // 损益
		gain = per
		for i, v := range dc.im.Pix {
			if i%4 == 3 { // alpha
				continue
			}
			dc.im.Pix[i] = uint8(limit2uint8(int(v) * gain / 100))
		}
	case 1 < per && per < 2: // 增益
		gain = 200 - per
		for i, v := range dc.im.Pix {
			if i%4 == 3 { // alpha
				continue
			}
			dc.im.Pix[i] = uint8(limit2uint8(int(v) * 100 / gain))
		}
	default:
		return
	}
}
