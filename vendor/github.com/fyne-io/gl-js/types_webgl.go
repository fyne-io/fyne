// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build js && !wasm
// +build js,!wasm

package gl

import "github.com/gopherjs/gopherjs/js"

type Enum int

type Attrib struct {
	Value int
}

type Program struct {
	*js.Object
}

type Shader struct {
	*js.Object
}

type Buffer struct {
	*js.Object
}

type Framebuffer struct {
	*js.Object
}

type Renderbuffer struct {
	*js.Object
}

type Texture struct {
	*js.Object
}

type Uniform struct {
	*js.Object
}

var NoAttrib Attrib
var NoProgram Program
var NoShader Shader
var NoBuffer Buffer
var NoFramebuffer Framebuffer
var NoRenderbuffer Renderbuffer
var NoTexture Texture
var NoUniform Uniform
