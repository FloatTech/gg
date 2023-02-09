package gg

import (
	"math/rand"
	"testing"
)

func quadraticnoasm(x0, y0, x1, y1, x2, y2, ds float64, p []Point) {
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

func cubicnoasm(x0, y0, x1, y1, x2, y2, x3, y3, ds float64, p []Point) {
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

func TestQuadraticCalc(t *testing.T) {
	p1 := make([]Point, 4096)
	p2 := make([]Point, 4096)
	for i := 0; i < 4096; i++ {
		x0, y0, x1, y1, x2, y2, ds := rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64()
		quadratic(x0, y0, x1, y1, x2, y2, ds, p1)
		quadraticnoasm(x0, y0, x1, y1, x2, y2, ds, p2)
		for j := 0; j < 4096; j++ {
			if p1[j].X != p2[j].X || p1[j].Y != p2[j].Y {
				t.Fatal()
			}
		}
	}
}

func BenchmarkQuadraticCalc(b *testing.B) {
	b.Run("16", func(b *testing.B) {
		benchmarkQuadratic(b, 16)
	})
	b.Run("256", func(b *testing.B) {
		benchmarkQuadratic(b, 256)
	})
	b.Run("512", func(b *testing.B) {
		benchmarkQuadratic(b, 512)
	})
	b.Run("1024", func(b *testing.B) {
		benchmarkQuadratic(b, 1024)
	})
	b.Run("2048", func(b *testing.B) {
		benchmarkQuadratic(b, 2048)
	})
	b.Run("4K", func(b *testing.B) {
		benchmarkQuadratic(b, 1024*4)
	})
	b.Run("32K", func(b *testing.B) {
		benchmarkQuadratic(b, 1024*32)
	})
}

func BenchmarkQuadraticNoASMCalc(b *testing.B) {
	b.Run("16", func(b *testing.B) {
		benchmarkQuadraticNoASM(b, 16)
	})
	b.Run("256", func(b *testing.B) {
		benchmarkQuadraticNoASM(b, 256)
	})
	b.Run("512", func(b *testing.B) {
		benchmarkQuadraticNoASM(b, 512)
	})
	b.Run("1024", func(b *testing.B) {
		benchmarkQuadraticNoASM(b, 1024)
	})
	b.Run("2048", func(b *testing.B) {
		benchmarkQuadraticNoASM(b, 2048)
	})
	b.Run("4K", func(b *testing.B) {
		benchmarkQuadraticNoASM(b, 1024*4)
	})
	b.Run("32K", func(b *testing.B) {
		benchmarkQuadraticNoASM(b, 1024*32)
	})
}

func benchmarkQuadratic(b *testing.B, plen int) {
	p := make([]Point, plen)
	x0, y0, x1, y1, x2, y2, ds := rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		quadratic(x0, y0, x1, y1, x2, y2, ds, p)
	}
}

func benchmarkQuadraticNoASM(b *testing.B, plen int) {
	p := make([]Point, plen)
	x0, y0, x1, y1, x2, y2, ds := rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		quadraticnoasm(x0, y0, x1, y1, x2, y2, ds, p)
	}
}

func TestCubicCalc(t *testing.T) {
	p1 := make([]Point, 4096)
	p2 := make([]Point, 4096)
	for i := 0; i < 4096; i++ {
		x0, y0, x1, y1, x2, y2, x3, y3, ds := rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64()
		cubic(x0, y0, x1, y1, x2, y2, x3, y3, ds, p1)
		cubicnoasm(x0, y0, x1, y1, x2, y2, x3, y3, ds, p2)
		for j := 0; j < 4096; j++ {
			if p1[j].X != p2[j].X || p1[j].Y != p2[j].Y {
				t.Fatal()
			}
		}
	}
}

func BenchmarkCubicCalc(b *testing.B) {
	b.Run("16", func(b *testing.B) {
		benchmarkCubic(b, 16)
	})
	b.Run("256", func(b *testing.B) {
		benchmarkCubic(b, 256)
	})
	b.Run("512", func(b *testing.B) {
		benchmarkCubic(b, 512)
	})
	b.Run("1024", func(b *testing.B) {
		benchmarkCubic(b, 1024)
	})
	b.Run("2048", func(b *testing.B) {
		benchmarkCubic(b, 2048)
	})
	b.Run("4K", func(b *testing.B) {
		benchmarkCubic(b, 1024*4)
	})
	b.Run("32K", func(b *testing.B) {
		benchmarkCubic(b, 1024*32)
	})
}

func BenchmarkCubicNoASMCalc(b *testing.B) {
	b.Run("16", func(b *testing.B) {
		benchmarkCubicNoASM(b, 16)
	})
	b.Run("256", func(b *testing.B) {
		benchmarkCubicNoASM(b, 256)
	})
	b.Run("512", func(b *testing.B) {
		benchmarkCubicNoASM(b, 512)
	})
	b.Run("1024", func(b *testing.B) {
		benchmarkCubicNoASM(b, 1024)
	})
	b.Run("2048", func(b *testing.B) {
		benchmarkCubicNoASM(b, 2048)
	})
	b.Run("4K", func(b *testing.B) {
		benchmarkCubicNoASM(b, 1024*4)
	})
	b.Run("32K", func(b *testing.B) {
		benchmarkCubicNoASM(b, 1024*32)
	})
}

func benchmarkCubic(b *testing.B, plen int) {
	p := make([]Point, plen)
	x0, y0, x1, y1, x2, y2, x3, y3, ds := rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cubic(x0, y0, x1, y1, x2, y2, x3, y3, ds, p)
	}
}

func benchmarkCubicNoASM(b *testing.B, plen int) {
	p := make([]Point, plen)
	x0, y0, x1, y1, x2, y2, x3, y3, ds := rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cubicnoasm(x0, y0, x1, y1, x2, y2, x3, y3, ds, p)
	}
}
