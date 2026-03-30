package gg

// sq square
//
// sq 计算平方
func sq(n float64) float64 {
	return n * n
}

// clamp controls n in [a, b]
//
// clamp 将 n 的范围限制在 [a, b]
func clamp(n, a, b int) int {
	if n > b {
		n = b
	}
	if n < a {
		n = a
	}
	return n
}
