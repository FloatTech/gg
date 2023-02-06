package gg

import (
	"math/rand"
	"testing"
)

func TestInnerQuadratic(t *testing.T) {
	var x0, y0, x1, y1, x2, y2, tf float64
	for i := 0; i < 4096; i++ {
		x0, y0, x1, y1, x2, y2, tf = rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64()
		xn, yn := quadraticnoasm(x0, y0, x1, y1, x2, y2, tf)
		x, y := quadratic(x0, y0, x1, y1, x2, y2, tf)
		if xn != x || yn != y {
			t.Fatal("[", i, "]", x0, y0, x1, y1, x2, y2, tf, "=", x, "(", xn, ")", y, "(", yn, ")")
		}
	}
}

func BenchmarkInnerQuadratic(b *testing.B) {
	x0, y0, x1, y1, x2, y2, tf := rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = quadratic(x0, y0, x1, y1, x2, y2, tf)
	}
}

func BenchmarkInnerQuadraticNoASM(b *testing.B) {
	x0, y0, x1, y1, x2, y2, tf := rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = quadraticnoasm(x0, y0, x1, y1, x2, y2, tf)
	}
}

func TestInnerCubic(t *testing.T) {
	var x0, y0, x1, y1, x2, y2, x3, y3, tf float64
	for i := 0; i < 4096; i++ {
		x0, y0, x1, y1, x2, y2, x3, y3, tf = rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64()
		xn, yn := cubicnoasm(x0, y0, x1, y1, x2, y2, x3, y3, tf)
		x, y := cubic(x0, y0, x1, y1, x2, y2, x3, y3, tf)
		if xn != x || yn != y {
			t.Fatal("[", i, "]", x0, y0, x1, y1, x2, y2, x3, y3, tf, "=", x, "(", xn, ")", y, "(", yn, ")")
		}
	}
}

func BenchmarkInnerCubic(b *testing.B) {
	x0, y0, x1, y1, x2, y2, x3, y3, tf := rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cubic(x0, y0, x1, y1, x2, y2, x3, y3, tf)
	}
}

func BenchmarkInnerCubicNoASM(b *testing.B) {
	x0, y0, x1, y1, x2, y2, x3, y3, tf := rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cubicnoasm(x0, y0, x1, y1, x2, y2, x3, y3, tf)
	}
}

func quadraticnoasm(x0, y0, x1, y1, x2, y2, t float64) (x, y float64) {
	u := 1 - t
	a := u * u
	b := 2 * u * t
	c := t * t
	x = a*x0 + b*x1 + c*x2
	y = a*y0 + b*y1 + c*y2
	return
}

func cubicnoasm(x0, y0, x1, y1, x2, y2, x3, y3, t float64) (x, y float64) {
	u := 1 - t
	a := u * u * u
	b := 3 * u * u * t
	c := 3 * u * t * t
	d := t * t * t
	x = a*x0 + b*x1 + c*x2 + d*x3
	y = a*y0 + b*y1 + c*y2 + d*y3
	return
}
