// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package f32

import "fmt"

type Vec4 [4]float32

func (v Vec4) String() string {
	return fmt.Sprintf("Vec4[% 0.3f, % 0.3f, % 0.3f, % 0.3f]", v[0], v[1], v[2], v[3])
}

func (v *Vec4) Normalize() {
	sq := v.Dot(v)
	inv := 1 / Sqrt(sq)
	v[0] *= inv
	v[1] *= inv
	v[2] *= inv
	v[3] *= inv
}

func (v *Vec4) Sub(v0, v1 *Vec4) {
	v[0] = v0[0] - v1[0]
	v[1] = v0[1] - v1[1]
	v[2] = v0[2] - v1[2]
	v[3] = v0[3] - v1[3]
}

func (v *Vec4) Add(v0, v1 *Vec4) {
	v[0] = v0[0] + v1[0]
	v[1] = v0[1] + v1[1]
	v[2] = v0[2] + v1[2]
	v[3] = v0[3] + v1[3]
}

func (v *Vec4) Mul(v0, v1 *Vec4) {
	v[0] = v0[0] * v1[0]
	v[1] = v0[1] * v1[1]
	v[2] = v0[2] * v1[2]
	v[3] = v0[3] * v1[3]
}

func (v *Vec4) Dot(v1 *Vec4) float32 {
	return v[0]*v1[0] + v[1]*v1[1] + v[2]*v1[2] + v[3]*v1[3]
}
