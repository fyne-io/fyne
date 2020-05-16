// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build js,wasm

package gl

import "syscall/js"

type Enum int

type Attrib struct {
	Value int
}

type Program struct {
	js.Value
}

type Shader struct {
	js.Value
}

type Buffer struct {
	js.Value
}

type Framebuffer struct {
	js.Value
}

type Renderbuffer struct {
	js.Value
}

type Texture struct {
	js.Value
}

type Uniform struct {
	js.Value
}

var NoAttrib = Attrib{0}
var NoProgram = Program{js.Null()}
var NoShader = Shader{js.Null()}
var NoBuffer = Buffer{js.Null()}
var NoFramebuffer = Framebuffer{js.Null()}
var NoRenderbuffer = Renderbuffer{js.Null()}
var NoTexture = Texture{js.Null()}
var NoUniform = Uniform{js.Null()}
