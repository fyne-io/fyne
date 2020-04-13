// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin linux openbsd windows
// +build gldebug

package gl

// Alternate versions of the types defined in types.go with extra
// debugging information attached. For documentation, see types.go.

import "fmt"

type Enum uint32

type Attrib struct {
	Value uint
	name  string
}

type Program struct {
	Init  bool
	Value uint32
}

type Shader struct {
	Value uint32
}

type Buffer struct {
	Value uint32
}

type Framebuffer struct {
	Value uint32
}

type Renderbuffer struct {
	Value uint32
}

type Texture struct {
	Value uint32
}

type Uniform struct {
	Value int32
	name  string
}

type VertexArray struct {
	Value uint32
}

func (v Attrib) c() uintptr { return uintptr(v.Value) }
func (v Enum) c() uintptr   { return uintptr(v) }
func (v Program) c() uintptr {
	if !v.Init {
		ret := uintptr(0)
		ret--
		return ret
	}
	return uintptr(v.Value)
}
func (v Shader) c() uintptr       { return uintptr(v.Value) }
func (v Buffer) c() uintptr       { return uintptr(v.Value) }
func (v Framebuffer) c() uintptr  { return uintptr(v.Value) }
func (v Renderbuffer) c() uintptr { return uintptr(v.Value) }
func (v Texture) c() uintptr      { return uintptr(v.Value) }
func (v Uniform) c() uintptr      { return uintptr(v.Value) }
func (v VertexArray) c() uintptr  { return uintptr(v.Value) }

func (v Attrib) String() string       { return fmt.Sprintf("Attrib(%d:%s)", v.Value, v.name) }
func (v Program) String() string      { return fmt.Sprintf("Program(%d)", v.Value) }
func (v Shader) String() string       { return fmt.Sprintf("Shader(%d)", v.Value) }
func (v Buffer) String() string       { return fmt.Sprintf("Buffer(%d)", v.Value) }
func (v Framebuffer) String() string  { return fmt.Sprintf("Framebuffer(%d)", v.Value) }
func (v Renderbuffer) String() string { return fmt.Sprintf("Renderbuffer(%d)", v.Value) }
func (v Texture) String() string      { return fmt.Sprintf("Texture(%d)", v.Value) }
func (v Uniform) String() string      { return fmt.Sprintf("Uniform(%d:%s)", v.Value, v.name) }
func (v VertexArray) String() string  { return fmt.Sprintf("VertexArray(%d)", v.Value) }
