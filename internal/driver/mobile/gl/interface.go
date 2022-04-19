// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

// Context is an OpenGL ES context.
//
// A Context has a method for every GL function supported by ES 2 or later.
// In a program compiled with ES 3 support.
//
// Calls are not safe for concurrent use. However calls can be made from
// any goroutine, the gl package removes the notion of thread-local
// context.
//
// Contexts are independent. Two contexts can be used concurrently.
type Context interface {
	// ActiveTexture sets the active texture unit.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glActiveTexture.xhtml
	ActiveTexture(texture Enum)

	// AttachShader attaches a shader to a program.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glAttachShader.xhtml
	AttachShader(p Program, s Shader)

	// BindBuffer binds a buffer.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindBuffer.xhtml
	BindBuffer(target Enum, b Buffer)
	// BindTexture binds a texture.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindTexture.xhtml
	BindTexture(target Enum, t Texture)

	// BindVertexArray binds a vertex array.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glBindVertexArray.xhtml
	BindVertexArray(rb VertexArray)

	// BlendColor sets the blend color.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendColor.xhtml
	BlendColor(red, green, blue, alpha float32)

	// BlendFunc sets the pixel blending factors.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glBlendFunc.xhtml
	BlendFunc(sfactor, dfactor Enum)

	// BufferData creates a new data store for the bound buffer object.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glBufferData.xhtml
	BufferData(target Enum, src []byte, usage Enum)
	// Clear clears the window.
	//
	// The behavior of Clear is influenced by the pixel ownership test,
	// the scissor test, dithering, and the buffer writemasks.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glClear.xhtml
	Clear(mask Enum)

	// ClearColor specifies the RGBA values used to clear color buffers.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glClearColor.xhtml
	ClearColor(red, green, blue, alpha float32)

	// CompileShader compiles the source code of s.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glCompileShader.xhtml
	CompileShader(s Shader)

	// CreateBuffer creates a buffer object.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenBuffers.xhtml
	CreateBuffer() Buffer

	// CreateProgram creates a new empty program object.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glCreateProgram.xhtml
	CreateProgram() Program

	// CreateShader creates a new empty shader object.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glCreateShader.xhtml
	CreateShader(ty Enum) Shader

	// CreateTexture creates a texture object.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenTextures.xhtml
	CreateTexture() Texture

	// CreateTVertexArray creates a vertex array.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGenVertexArrays.xhtml
	CreateVertexArray() VertexArray
	// DeleteBuffer deletes the given buffer object.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteBuffers.xhtml
	DeleteBuffer(v Buffer)

	// DeleteTexture deletes the given texture object.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glDeleteTextures.xhtml
	DeleteTexture(v Texture)

	// Disable disables various GL capabilities.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glDisable.xhtml
	Disable(cap Enum)

	// DrawArrays renders geometric primitives from the bound data.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glDrawArrays.xhtml
	DrawArrays(mode Enum, first, count int)

	// Enable enables various GL capabilities.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glEnable.xhtml
	Enable(cap Enum)

	// EnableVertexAttribArray enables a vertex attribute array.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glEnableVertexAttribArray.xhtml
	EnableVertexAttribArray(a Attrib)
	// Flush empties all buffers. It does not block.
	//
	// An OpenGL implementation may buffer network communication,
	// the command stream, or data inside the graphics accelerator.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glFlush.xhtml
	Flush()

	// GetAttribLocation returns the location of an attribute variable.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetAttribLocation.xhtml
	GetAttribLocation(p Program, name string) Attrib

	// GetError returns the next error.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetError.xhtml
	GetError() Enum

	// GetProgrami returns a parameter value for a shader.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetProgramiv.xhtml
	GetProgrami(p Program, pname Enum) int

	// GetProgramInfoLog returns the information log for a shader.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetProgramInfoLog.xhtml
	GetProgramInfoLog(p Program) string

	// GetShaderi returns a parameter value for a shader.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetShaderiv.xhtml
	GetShaderi(s Shader, pname Enum) int

	// GetShaderInfoLog returns the information log for a shader.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetShaderInfoLog.xhtml
	GetShaderInfoLog(s Shader) string

	// GetUniformLocation returns the location of a uniform variable.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glGetUniformLocation.xhtml
	GetUniformLocation(p Program, name string) Uniform

	// LinkProgram links the specified program.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glLinkProgram.xhtml
	LinkProgram(p Program)

	// ReadPixels returns pixel data from a buffer.
	//
	// In GLES 3, the source buffer is controlled with ReadBuffer.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glReadPixels.xhtml
	ReadPixels(dst []byte, x, y, width, height int, format, ty Enum)

	// Scissor defines the scissor box rectangle, in window coordinates.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glScissor.xhtml
	Scissor(x, y, width, height int32)

	// ShaderSource sets the source code of s to the given source code.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glShaderSource.xhtml
	ShaderSource(s Shader, src string)
	// TexImage2D writes a 2D texture image.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexImage2D.xhtml
	TexImage2D(target Enum, level int, internalFormat int, width, height int, format Enum, ty Enum, data []byte)

	// TexParameteri sets an integer texture parameter.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glTexParameter.xhtml
	TexParameteri(target, pname Enum, param int)

	// Uniform1f writes a float uniform variable.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
	Uniform1f(dst Uniform, v float32)

	// Uniform4f writes a vec4 uniform variable.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
	Uniform4f(dst Uniform, v0, v1, v2, v3 float32)

	// Uniform4fv writes a vec4 uniform array of len(src)/4 elements.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glUniform.xhtml
	Uniform4fv(dst Uniform, src []float32)
	// UseProgram sets the active program.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glUseProgram.xhtml
	UseProgram(p Program)

	// VertexAttribPointer uses a bound buffer to define vertex attribute data.
	//
	// Direct use of VertexAttribPointer to load data into OpenGL is not
	// supported via the Go bindings. Instead, use BindBuffer with an
	// ARRAY_BUFFER and then fill it using BufferData.
	//
	// The size argument specifies the number of components per attribute,
	// between 1-4. The stride argument specifies the byte offset between
	// consecutive vertex attributes.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glVertexAttribPointer.xhtml
	VertexAttribPointer(dst Attrib, size int, ty Enum, normalized bool, stride, offset int)

	// Viewport sets the viewport, an affine transformation that
	// normalizes device coordinates to window coordinates.
	//
	// http://www.khronos.org/opengles/sdk/docs/man3/html/glViewport.xhtml
	Viewport(x, y, width, height int)
}

// Worker is used by display driver code to execute OpenGL calls.
//
// Typically display driver code creates a gl.Context for an application,
// and along with it establishes a locked OS thread to execute the cgo
// calls:
//
//	go func() {
//		runtime.LockOSThread()
//		// ... platform-specific cgo call to bind a C OpenGL context
//		// into thread-local storage.
//
//		glctx, worker := gl.NewContext()
//		workAvailable := worker.WorkAvailable()
//		go userAppCode(glctx)
//		for {
//			select {
//			case <-workAvailable:
//				worker.DoWork()
//			case <-drawEvent:
//				// ... platform-specific cgo call to draw screen
//			}
//		}
//	}()
//
// This interface is an internal implementation detail and should only be used
// by the package responsible for managing the screen.
type Worker interface {
	// WorkAvailable returns a channel that communicates when DoWork should be
	// called.
	WorkAvailable() <-chan struct{}

	// DoWork performs any pending OpenGL calls.
	DoWork()
}
