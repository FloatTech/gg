// Package path is previously in gg main package.
package path

import (
	"math"

	"github.com/FloatTech/gg/2d/bezeir"
	"github.com/FloatTech/gg/2d/point"
	"github.com/golang/freetype/raster"
	"golang.org/x/image/math/fixed"
)

// Flatten converts a raster path into point-only subpaths by approximating
// quadratic and cubic segments with sampled points.
//
// Flatten 将 raster 路径转换为仅含点的子路径，并用采样点近似二次与三次曲线段。
func Flatten(p raster.Path) [][]point.Point {
	var result = make([][]point.Point, len(p)*2)
	var path = make([]point.Point, 0, len(p)*2)
	var cx, cy float64
	for i := 0; i < len(p); {
		switch p[i] {
		case 0:
			if len(path) > 0 {
				result = append(result, path)
				path = nil
			}
			x := point.DeFixed(p[i+1])
			y := point.DeFixed(p[i+2])
			path = append(path, point.Point{X: x, Y: y})
			cx, cy = x, y
			i += 4
		case 1:
			x := point.DeFixed(p[i+1])
			y := point.DeFixed(p[i+2])
			path = append(path, point.Point{X: x, Y: y})
			cx, cy = x, y
			i += 4
		case 2:
			x1 := point.DeFixed(p[i+1])
			y1 := point.DeFixed(p[i+2])
			x2 := point.DeFixed(p[i+3])
			y2 := point.DeFixed(p[i+4])
			n := bezeir.QuadraticBezierLen(cx, cy, x1, y1, x2, y2)
			a := len(path)
			path = append(path, make([]point.Point, n)...)
			bezeir.QuadraticBezierInplace(cx, cy, x1, y1, x2, y2, float64(n)-1, path[a:])
			cx, cy = x2, y2
			i += 6
		case 3:
			x1 := point.DeFixed(p[i+1])
			y1 := point.DeFixed(p[i+2])
			x2 := point.DeFixed(p[i+3])
			y2 := point.DeFixed(p[i+4])
			x3 := point.DeFixed(p[i+5])
			y3 := point.DeFixed(p[i+6])
			n := bezeir.CubicBezierLen(cx, cy, x1, y1, x2, y2, x3, y3)
			a := len(path)
			path = append(path, make([]point.Point, n)...)
			bezeir.CubicBezierInplace(cx, cy, x1, y1, x2, y2, x3, y3, float64(n)-1, path[a:])
			cx, cy = x3, y3
			i += 8
		default:
			panic("bad path")
		}
	}
	if len(path) > 0 {
		result = append(result, path)
	}
	return result
}

// Dash applies a dash pattern and offset to flattened paths, returning only
// the visible stroke segments.
//
// Dash 将虚线模式及偏移应用到展平路径上，并返回可见的描边线段。
func Dash(paths [][]point.Point, dashes []float64, offset float64) [][]point.Point {
	if len(dashes) == 0 {
		return paths
	}
	if len(dashes) == 1 {
		dashes = append(dashes, dashes[0])
	}
	var result = make([][]point.Point, 0, len(paths)*2)
	for _, path := range paths {
		if len(path) < 2 {
			continue
		}
		previous := path[0]
		pathIndex := 1
		dashIndex := 0
		segmentLength := 0.0

		// offset
		if offset != 0 {
			var totalLength float64
			for _, dashLength := range dashes {
				totalLength += dashLength
			}
			offset = math.Mod(offset, totalLength)
			if offset < 0 {
				offset += totalLength
			}
			for i, dashLength := range dashes {
				offset -= dashLength
				if offset < 0 {
					dashIndex = i
					segmentLength = dashLength + offset
					break
				}
			}
		}

		var segment = make([]point.Point, 0, (len(path)-pathIndex)*2)
		segment = append(segment, previous)
		for pathIndex < len(path) {
			dashLength := dashes[dashIndex]
			pidx := path[pathIndex]
			d := previous.Distance(pidx)
			maxd := dashLength - segmentLength
			if d > maxd {
				t := maxd / d
				p := previous.Interpolate(pidx, t)
				segment = append(segment, p)
				if dashIndex%2 == 0 && len(segment) > 1 {
					result = append(result, segment)
				}
				segment = make([]point.Point, 0, (len(path)-pathIndex)*2)
				segment = append(segment, p)
				segmentLength = 0
				previous = p
				dashIndex = (dashIndex + 1) % len(dashes)
			} else {
				segment = append(segment, pidx)
				previous = pidx
				segmentLength += d
				pathIndex++
			}
		}
		if dashIndex%2 == 0 && len(segment) > 1 {
			result = append(result, segment)
		}
	}
	return result
}

// Raster converts point subpaths back to a raster path and skips extremely
// short edges to avoid rendering artifacts on joins and caps.
//
// Raster 将点子路径转换回 raster 路径，并跳过极短边以减少连接和端点渲染伪影。
func Raster(paths [][]point.Point) raster.Path {
	var result = make(raster.Path, 0, len(paths))
	for _, path := range paths {
		var previous fixed.Point26_6
		for i, point := range path {
			f := point.Fixed()
			if i == 0 {
				result.Start(f)
			} else {
				dx := f.X - previous.X
				dy := f.Y - previous.Y
				if dx < 0 {
					dx = -dx
				}
				if dy < 0 {
					dy = -dy
				}
				if dx+dy > 8 {
					// TODO: this is a hack for cases where two points are
					// too close - causes rendering issues with joins / caps
					result.Add1(f)
				}
			}
			previous = f
		}
	}
	return result
}

// Dashed returns a raster path with the dash pattern applied by flattening,
// dashing, and rasterizing in sequence.
//
// Dashed 按顺序执行展平、虚线切分与栅格化，返回应用虚线模式后的 raster 路径。
func Dashed(path raster.Path, dashes []float64, offset float64) raster.Path {
	return Raster(Dash(Flatten(path), dashes, offset))
}
