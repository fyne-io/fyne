// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package f32

import "fmt"

type Vec3 [3]float32

func (v Vec3) String() string {
	return fmt.Sprintf("Vec3[% 0.3f, % 0.3f, % 0.3f]", v[0], v[1], v[2])
}

func (v *Vec3) Normalize() {
	sq := v.Dot(v)
	inv := 1 / Sqrt(sq)
	v[0] *= inv
	v[1] *= inv
	v[2] *= inv
}

func (v *Vec3) Sub(v0, v1 *Vec3) {
	v[0] = v0[0] - v1[0]
	v[1] = v0[1] - v1[1]
	v[2] = v0[2] - v1[2]
}

func (v *Vec3) Add(v0, v1 *Vec3) {
	v[0] = v0[0] + v1[0]
	v[1] = v0[1] + v1[1]
	v[2] = v0[2] + v1[2]
}

func (v *Vec3) Mul(v0, v1 *Vec3) {
	v[0] = v0[0] * v1[0]
	v[1] = v0[1] * v1[1]
	v[2] = v0[2] * v1[2]
}

func (v *Vec3) Cross(v0, v1 *Vec3) {
	v[0] = v0[1]*v1[2] - v0[2]*v1[1]
	v[1] = v0[2]*v1[0] - v0[0]*v1[2]
	v[2] = v0[0]*v1[1] - v0[1]*v1[0]
}

func (v *Vec3) Dot(v1 *Vec3) float32 {
	return v[0]*v1[0] + v[1]*v1[1] + v[2]*v1[2]
}
