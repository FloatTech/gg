//go:build !amd64

package gg

func quadratic(x0, y0, x1, y1, x2, y2, ds float64, p []Point) {
	var u, a, b, c, t float64
	for i := 0; i < len(p); i++ {
		t = float64(i) / ds
		u = 1 - t
		a = u * u
		b = 2 * u * t
		c = t * t
		p[i].X, p[i].Y = a*x0+b*x1+c*x2, a*y0+b*y1+c*y2
	}
}

func cubic(x0, y0, x1, y1, x2, y2, x3, y3, ds float64, p []Point) {
	var u, a, b, c, d, t float64
	for i := 0; i < len(p); i++ {
		t = float64(i) / ds
		u = 1 - t
		a = u * u * u
		b = 3 * u * u * t
		c = 3 * u * t * t
		d = t * t * t
		p[i].X, p[i].Y = a*x0+b*x1+c*x2+d*x3, a*y0+b*y1+c*y2+d*y3
	}
}
