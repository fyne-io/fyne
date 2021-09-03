// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package f32

import "fmt"

// A Mat3 is a 3x3 matrix of float32 values.
// Elements are indexed first by row then column, i.e. m[row][column].
type Mat3 [3]Vec3

func (m Mat3) String() string {
	return fmt.Sprintf(`Mat3[% 0.3f, % 0.3f, % 0.3f,
     % 0.3f, % 0.3f, % 0.3f,
     % 0.3f, % 0.3f, % 0.3f]`,
		m[0][0], m[0][1], m[0][2],
		m[1][0], m[1][1], m[1][2],
		m[2][0], m[2][1], m[2][2])
}

func (m *Mat3) Identity() {
	*m = Mat3{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
	}
}

func (m *Mat3) Eq(n *Mat3, epsilon float32) bool {
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

// Mul stores a × b in m.
func (m *Mat3) Mul(a, b *Mat3) {
	// Store the result in local variables, in case m == a || m == b.
	m00 := a[0][0]*b[0][0] + a[0][1]*b[1][0] + a[0][2]*b[2][0]
	m01 := a[0][0]*b[0][1] + a[0][1]*b[1][1] + a[0][2]*b[2][1]
	m02 := a[0][0]*b[0][2] + a[0][1]*b[1][2] + a[0][2]*b[2][2]
	m10 := a[1][0]*b[0][0] + a[1][1]*b[1][0] + a[1][2]*b[2][0]
	m11 := a[1][0]*b[0][1] + a[1][1]*b[1][1] + a[1][2]*b[2][1]
	m12 := a[1][0]*b[0][2] + a[1][1]*b[1][2] + a[1][2]*b[2][2]
	m20 := a[2][0]*b[0][0] + a[2][1]*b[1][0] + a[2][2]*b[2][0]
	m21 := a[2][0]*b[0][1] + a[2][1]*b[1][1] + a[2][2]*b[2][1]
	m22 := a[2][0]*b[0][2] + a[2][1]*b[1][2] + a[2][2]*b[2][2]
	m[0][0] = m00
	m[0][1] = m01
	m[0][2] = m02
	m[1][0] = m10
	m[1][1] = m11
	m[1][2] = m12
	m[2][0] = m20
	m[2][1] = m21
	m[2][2] = m22
}
