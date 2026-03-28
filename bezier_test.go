package gg

import (
	"math"
	"math/rand"
	"testing"
)

func TestQuadraticBezeirCalc(t *testing.T) {
	p1 := make([]Point, 4096)
	p2 := make([]Point, 4096)
	for range 4096 {
		x0, y0, x1, y1, x2, y2, ds := rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64()
		quadraticBezier(x0, y0, x1, y1, x2, y2, ds, p1)
		quadraticBezierPure(x0, y0, x1, y1, x2, y2, ds, p2)
		for j := range 4096 {
			tolX := 0.00001 * (1 + math.Abs(p2[j].X))
			tolY := 0.00001 * (1 + math.Abs(p2[j].Y))
			if math.Abs(p1[j].X-p2[j].X) >= tolX || math.Abs(p1[j].Y-p2[j].Y) >= tolY {
				t.Fatalf("No.%d expect (%.2f, %.2f) but got (%.2f, %.2f)", j, p2[j].X, p2[j].Y, p1[j].X, p1[j].Y)
			}
		}
	}
}

func BenchmarkQuadraticBezeirCalc(b *testing.B) {
	b.Run("16", func(b *testing.B) {
		benchmarkQuadraticBezeir(b, 16)
	})
	b.Run("256", func(b *testing.B) {
		benchmarkQuadraticBezeir(b, 256)
	})
	b.Run("512", func(b *testing.B) {
		benchmarkQuadraticBezeir(b, 512)
	})
	b.Run("1024", func(b *testing.B) {
		benchmarkQuadraticBezeir(b, 1024)
	})
	b.Run("2048", func(b *testing.B) {
		benchmarkQuadraticBezeir(b, 2048)
	})
	b.Run("4K", func(b *testing.B) {
		benchmarkQuadraticBezeir(b, 1024*4)
	})
	b.Run("32K", func(b *testing.B) {
		benchmarkQuadraticBezeir(b, 1024*32)
	})
}

func BenchmarkQuadraticBezeirPureCalc(b *testing.B) {
	b.Run("16", func(b *testing.B) {
		benchmarkQuadraticBezeirPure(b, 16)
	})
	b.Run("256", func(b *testing.B) {
		benchmarkQuadraticBezeirPure(b, 256)
	})
	b.Run("512", func(b *testing.B) {
		benchmarkQuadraticBezeirPure(b, 512)
	})
	b.Run("1024", func(b *testing.B) {
		benchmarkQuadraticBezeirPure(b, 1024)
	})
	b.Run("2048", func(b *testing.B) {
		benchmarkQuadraticBezeirPure(b, 2048)
	})
	b.Run("4K", func(b *testing.B) {
		benchmarkQuadraticBezeirPure(b, 1024*4)
	})
	b.Run("32K", func(b *testing.B) {
		benchmarkQuadraticBezeirPure(b, 1024*32)
	})
}

func benchmarkQuadraticBezeir(b *testing.B, plen int) {
	p := make([]Point, plen)
	x0, y0, x1, y1, x2, y2, ds := rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64()
	b.ResetTimer()
	for range b.N {
		quadraticBezier(x0, y0, x1, y1, x2, y2, ds, p)
	}
}

func benchmarkQuadraticBezeirPure(b *testing.B, plen int) {
	p := make([]Point, plen)
	x0, y0, x1, y1, x2, y2, ds := rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64()
	b.ResetTimer()
	for range b.N {
		quadraticBezierPure(x0, y0, x1, y1, x2, y2, ds, p)
	}
}

func TestCubicBezeirCalc(t *testing.T) {
	p1 := make([]Point, 4096)
	p2 := make([]Point, 4096)
	for range 4096 {
		x0, y0, x1, y1, x2, y2, x3, y3, ds := rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64()
		cubicBezier(x0, y0, x1, y1, x2, y2, x3, y3, ds, p1)
		cubicBezierPure(x0, y0, x1, y1, x2, y2, x3, y3, ds, p2)
		for j := range 4096 {
			tolX := 0.00001 * (1 + math.Abs(p2[j].X))
			tolY := 0.00001 * (1 + math.Abs(p2[j].Y))
			if math.Abs(p1[j].X-p2[j].X) >= tolX || math.Abs(p1[j].Y-p2[j].Y) >= tolY {
				t.Fatalf("No.%d expect (%.2f, %.2f) but got (%.2f, %.2f)", j, p2[j].X, p2[j].Y, p1[j].X, p1[j].Y)
			}
		}
	}
}

func BenchmarkCubicBezeirCalc(b *testing.B) {
	b.Run("16", func(b *testing.B) {
		benchmarkCubicBezeir(b, 16)
	})
	b.Run("256", func(b *testing.B) {
		benchmarkCubicBezeir(b, 256)
	})
	b.Run("512", func(b *testing.B) {
		benchmarkCubicBezeir(b, 512)
	})
	b.Run("1024", func(b *testing.B) {
		benchmarkCubicBezeir(b, 1024)
	})
	b.Run("2048", func(b *testing.B) {
		benchmarkCubicBezeir(b, 2048)
	})
	b.Run("4K", func(b *testing.B) {
		benchmarkCubicBezeir(b, 1024*4)
	})
	b.Run("32K", func(b *testing.B) {
		benchmarkCubicBezeir(b, 1024*32)
	})
}

func BenchmarkCubicBezeirPureCalc(b *testing.B) {
	b.Run("16", func(b *testing.B) {
		benchmarkCubicBezeirPure(b, 16)
	})
	b.Run("256", func(b *testing.B) {
		benchmarkCubicBezeirPure(b, 256)
	})
	b.Run("512", func(b *testing.B) {
		benchmarkCubicBezeirPure(b, 512)
	})
	b.Run("1024", func(b *testing.B) {
		benchmarkCubicBezeirPure(b, 1024)
	})
	b.Run("2048", func(b *testing.B) {
		benchmarkCubicBezeirPure(b, 2048)
	})
	b.Run("4K", func(b *testing.B) {
		benchmarkCubicBezeirPure(b, 1024*4)
	})
	b.Run("32K", func(b *testing.B) {
		benchmarkCubicBezeirPure(b, 1024*32)
	})
}

func benchmarkCubicBezeir(b *testing.B, plen int) {
	p := make([]Point, plen)
	x0, y0, x1, y1, x2, y2, x3, y3, ds := rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64()
	b.ResetTimer()
	for range b.N {
		cubicBezier(x0, y0, x1, y1, x2, y2, x3, y3, ds, p)
	}
}

func benchmarkCubicBezeirPure(b *testing.B, plen int) {
	p := make([]Point, plen)
	x0, y0, x1, y1, x2, y2, x3, y3, ds := rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64()
	b.ResetTimer()
	for range b.N {
		cubicBezierPure(x0, y0, x1, y1, x2, y2, x3, y3, ds, p)
	}
}
