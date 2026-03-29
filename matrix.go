package gg

import "math"

// Matrix represents a 3x2 affine transformation matrix.
//
// Matrix 表示一个 3x2 仿射变换矩阵。
type Matrix struct {
	XX, YX, XY, YY, X0, Y0 float64
}

// Identity returns the identity transformation matrix.
//
// Identity 返回单位变换矩阵。
func Identity() Matrix {
	return Matrix{
		1, 0,
		0, 1,
		0, 0,
	}
}

// Translate returns a translation matrix.
//
// Translate 返回平移变换矩阵。
func Translate(x, y float64) Matrix {
	return Matrix{
		1, 0,
		0, 1,
		x, y,
	}
}

// Scale returns a scaling matrix.
//
// Scale 返回缩放变换矩阵。
func Scale(x, y float64) Matrix {
	return Matrix{
		x, 0,
		0, y,
		0, 0,
	}
}

// Rotate returns a rotation matrix for the given angle in radians.
//
// Rotate 返回指定弧度角的旋转矩阵。
func Rotate(angle float64) Matrix {
	c := math.Cos(angle)
	s := math.Sin(angle)
	return Matrix{
		c, s,
		-s, c,
		0, 0,
	}
}

// Shear returns a shearing matrix.
//
// Shear 返回剪切变换矩阵。
func Shear(x, y float64) Matrix {
	return Matrix{
		1, y,
		x, 1,
		0, 0,
	}
}

// Multiply multiplies two matrices and returns the result.
//
// Multiply 将两个矩阵相乘并返回结果。
func (a Matrix) Multiply(b Matrix) Matrix {
	return Matrix{
		a.XX*b.XX + a.YX*b.XY,
		a.XX*b.YX + a.YX*b.YY,
		a.XY*b.XX + a.YY*b.XY,
		a.XY*b.YX + a.YY*b.YY,
		a.X0*b.XX + a.Y0*b.XY + b.X0,
		a.X0*b.YX + a.Y0*b.YY + b.Y0,
	}
}

// TransformVector transforms a vector (without translation).
//
// TransformVector 变换向量（不包含平移）。
func (a Matrix) TransformVector(x, y float64) (tx, ty float64) {
	tx = a.XX*x + a.XY*y
	ty = a.YX*x + a.YY*y
	return
}

// TransformPoint transforms a point (with translation).
//
// TransformPoint 变换点（包含平移）。
func (a Matrix) TransformPoint(x, y float64) (tx, ty float64) {
	tx = a.XX*x + a.XY*y + a.X0
	ty = a.YX*x + a.YY*y + a.Y0
	return
}

// Translate returns a new matrix with a translation applied.
//
// Translate 返回应用了平移的新矩阵。
func (a Matrix) Translate(x, y float64) Matrix {
	return Translate(x, y).Multiply(a)
}

// Scale returns a new matrix with a scaling applied.
//
// Scale 返回应用了缩放的新矩阵。
func (a Matrix) Scale(x, y float64) Matrix {
	return Scale(x, y).Multiply(a)
}

// Rotate returns a new matrix with a rotation applied.
//
// Rotate 返回应用了旋转的新矩阵。
func (a Matrix) Rotate(angle float64) Matrix {
	return Rotate(angle).Multiply(a)
}

// Shear returns a new matrix with a shear applied.
//
// Shear 返回应用了剪切的新矩阵。
func (a Matrix) Shear(x, y float64) Matrix {
	return Shear(x, y).Multiply(a)
}
