// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package f32

import "fmt"

// An Affine is a 3x3 matrix of float32 values for which the bottom row is
// implicitly always equal to [0 0 1].
// Elements are indexed first by row then column, i.e. m[row][column].
type Affine [2]Vec3

func (m Affine) String() string {
	return fmt.Sprintf(`Affine[% 0.3f, % 0.3f, % 0.3f,
     % 0.3f, % 0.3f, % 0.3f]`,
		m[0][0], m[0][1], m[0][2],
		m[1][0], m[1][1], m[1][2])
}

// Identity sets m to be the identity transform.
func (m *Affine) Identity() {
	*m = Affine{
		{1, 0, 0},
		{0, 1, 0},
	}
}

// Eq reports whether each component of m is within epsilon of the same
// component in n.
func (m *Affine) Eq(n *Affine, epsilon float32) bool {
	for i := range m {
		for j := range m[i] {
			diff := m[i][j] - n[i][j]
			if diff < -epsilon || +epsilon < diff {
				return false
			}
		}
	}
	return true
}

// Mul sets m to be p × q.
func (m *Affine) Mul(p, q *Affine) {
	// Store the result in local variables, in case m == a || m == b.
	m00 := p[0][0]*q[0][0] + p[0][1]*q[1][0]
	m01 := p[0][0]*q[0][1] + p[0][1]*q[1][1]
	m02 := p[0][0]*q[0][2] + p[0][1]*q[1][2] + p[0][2]
	m10 := p[1][0]*q[0][0] + p[1][1]*q[1][0]
	m11 := p[1][0]*q[0][1] + p[1][1]*q[1][1]
	m12 := p[1][0]*q[0][2] + p[1][1]*q[1][2] + p[1][2]
	m[0][0] = m00
	m[0][1] = m01
	m[0][2] = m02
	m[1][0] = m10
	m[1][1] = m11
	m[1][2] = m12
}

// Inverse sets m to be the inverse of p.
func (m *Affine) Inverse(p *Affine) {
	m00 := p[1][1]
	m01 := -p[0][1]
	m02 := p[1][2]*p[0][1] - p[1][1]*p[0][2]
	m10 := -p[1][0]
	m11 := p[0][0]
	m12 := p[1][0]*p[0][2] - p[1][2]*p[0][0]

	det := m00*m11 - m10*m01

	m[0][0] = m00 / det
	m[0][1] = m01 / det
	m[0][2] = m02 / det
	m[1][0] = m10 / det
	m[1][1] = m11 / det
	m[1][2] = m12 / det
}

// Scale sets m to be a scale followed by p.
// It is equivalent to m.Mul(p, &Affine{{x,0,0}, {0,y,0}}).
func (m *Affine) Scale(p *Affine, x, y float32) {
	m[0][0] = p[0][0] * x
	m[0][1] = p[0][1] * y
	m[0][2] = p[0][2]
	m[1][0] = p[1][0] * x
	m[1][1] = p[1][1] * y
	m[1][2] = p[1][2]
}

// Translate sets m to be a translation followed by p.
// It is equivalent to m.Mul(p, &Affine{{1,0,x}, {0,1,y}}).
func (m *Affine) Translate(p *Affine, x, y float32) {
	m[0][0] = p[0][0]
	m[0][1] = p[0][1]
	m[0][2] = p[0][0]*x + p[0][1]*y + p[0][2]
	m[1][0] = p[1][0]
	m[1][1] = p[1][1]
	m[1][2] = p[1][0]*x + p[1][1]*y + p[1][2]
}

// Rotate sets m to a rotation in radians followed by p.
// It is equivalent to m.Mul(p, affineRotation).
func (m *Affine) Rotate(p *Affine, radians float32) {
	s, c := Sin(radians), Cos(radians)
	m.Mul(p, &Affine{
		{+c, +s, 0},
		{-s, +c, 0},
	})
}
