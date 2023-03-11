// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin || linux || openbsd || freebsd
// +build darwin linux openbsd freebsd

#include <stdlib.h>
#include "_cgo_export.h"
#include "work.h"

#if defined(GL_ES_VERSION_3_0) && GL_ES_VERSION_3_0
#else
#include <stdio.h>
static void gles3missing() {
	printf("GLES3 function is missing\n");
	exit(2);
}
static void glBindVertexArray(GLuint array) { gles3missing(); }
static void glGenVertexArrays(GLsizei n, GLuint *arrays) { gles3missing(); }
#endif

uintptr_t processFn(struct fnargs* args, char* parg) {
	uintptr_t ret = 0;
	switch (args->fn) {
	case glfnUNDEFINED:
		abort(); // bad glfn
		break;
	case glfnActiveTexture:
		glActiveTexture((GLenum)args->a0);
		break;
	case glfnAttachShader:
		glAttachShader((GLint)args->a0, (GLint)args->a1);
		break;
	case glfnBindBuffer:
		glBindBuffer((GLenum)args->a0, (GLuint)args->a1);
		break;
	case glfnBindTexture:
		glBindTexture((GLenum)args->a0, (GLint)args->a1);
		break;
	case glfnBindVertexArray:
		glBindVertexArray((GLenum)args->a0);
		break;
	case glfnBlendColor:
		glBlendColor(*(GLfloat*)&args->a0, *(GLfloat*)&args->a1, *(GLfloat*)&args->a2, *(GLfloat*)&args->a3);
		break;
	case glfnBlendFunc:
		glBlendFunc((GLenum)args->a0, (GLenum)args->a1);
		break;
	case glfnBufferData:
		glBufferData((GLenum)args->a0, (GLsizeiptr)args->a1, (GLvoid*)parg, (GLenum)args->a2);
		break;
	case glfnClear:
		glClear((GLenum)args->a0);
		break;
	case glfnClearColor:
		glClearColor(*(GLfloat*)&args->a0, *(GLfloat*)&args->a1, *(GLfloat*)&args->a2, *(GLfloat*)&args->a3);
		break;
	case glfnCompileShader:
		glCompileShader((GLint)args->a0);
		break;
	case glfnCreateProgram:
		ret = glCreateProgram();
		break;
	case glfnCreateShader:
		ret = glCreateShader((GLenum)args->a0);
		break;
	case glfnDeleteBuffer:
		glDeleteBuffers(1, (const GLuint*)(&args->a0));
		break;
	case glfnDeleteTexture:
		glDeleteTextures(1, (const GLuint*)(&args->a0));
		break;
	case glfnDisable:
		glDisable((GLenum)args->a0);
		break;
	case glfnDrawArrays:
		glDrawArrays((GLenum)args->a0, (GLint)args->a1, (GLint)args->a2);
		break;
	case glfnEnable:
		glEnable((GLenum)args->a0);
		break;
	case glfnEnableVertexAttribArray:
		glEnableVertexAttribArray((GLint)args->a0);
		break;
	case glfnFlush:
		glFlush();
		break;
	case glfnGenBuffer:
		glGenBuffers(1, (GLuint*)&ret);
		break;
	case glfnGenTexture:
		glGenTextures(1, (GLuint*)&ret);
		break;
	case glfnGenVertexArray:
		glGenVertexArrays(1, (GLuint*)&ret);
		break;
	case glfnGetAttribLocation:
		ret = glGetAttribLocation((GLint)args->a0, (GLchar*)args->a1);
		break;
	case glfnGetError:
		ret = glGetError();
		break;
	case glfnGetProgramiv:
		glGetProgramiv((GLint)args->a0, (GLenum)args->a1, (GLint*)&ret);
		break;
	case glfnGetProgramInfoLog:
		glGetProgramInfoLog((GLuint)args->a0, (GLsizei)args->a1, 0, (GLchar*)parg);
		break;
	case glfnGetShaderiv:
		glGetShaderiv((GLint)args->a0, (GLenum)args->a1, (GLint*)&ret);
		break;
	case glfnGetShaderInfoLog:
		glGetShaderInfoLog((GLuint)args->a0, (GLsizei)args->a1, 0, (GLchar*)parg);
		break;
	case glfnGetShaderSource:
		glGetShaderSource((GLuint)args->a0, (GLsizei)args->a1, 0, (GLchar*)parg);
		break;
	case glfnGetTexParameteriv:
		glGetTexParameteriv((GLenum)args->a0, (GLenum)args->a1, (GLint*)parg);
		break;
	case glfnGetUniformLocation:
		ret = glGetUniformLocation((GLint)args->a0, (GLchar*)args->a1);
		break;
	case glfnLinkProgram:
		glLinkProgram((GLint)args->a0);
		break;
	case glfnReadPixels:
		glReadPixels((GLint)args->a0, (GLint)args->a1, (GLsizei)args->a2, (GLsizei)args->a3, (GLenum)args->a4, (GLenum)args->a5, (void*)parg);
		break;
	case glfnScissor:
		glScissor((GLint)args->a0, (GLint)args->a1, (GLint)args->a2, (GLint)args->a3);
		break;
	case glfnShaderSource:
#if defined(os_ios) || defined(os_macos)
		glShaderSource((GLuint)args->a0, (GLsizei)args->a1, (const GLchar *const *)args->a2, NULL);
#else
		glShaderSource((GLuint)args->a0, (GLsizei)args->a1, (const GLchar **)args->a2, NULL);
#endif
		break;
	case glfnTexImage2D:
		glTexImage2D(
			(GLenum)args->a0,
			(GLint)args->a1,
			(GLint)args->a2,
			(GLsizei)args->a3,
			(GLsizei)args->a4,
			0, // border
			(GLenum)args->a5,
			(GLenum)args->a6,
			(const GLvoid*)parg);
		break;
	case glfnTexParameteri:
		glTexParameteri((GLenum)args->a0, (GLenum)args->a1, (GLint)args->a2);
		break;
	case glfnUniform1f:
		glUniform1f((GLint)args->a0, *(GLfloat*)&args->a1);
		break;
	case glfnUniform4f:
		glUniform4f((GLint)args->a0, *(GLfloat*)&args->a1, *(GLfloat*)&args->a2, *(GLfloat*)&args->a3, *(GLfloat*)&args->a4);
		break;
	case glfnUniform4fv:
		glUniform4fv((GLint)args->a0, (GLsizeiptr)args->a1, (GLvoid*)parg);
		break;
	case glfnUseProgram:
		glUseProgram((GLint)args->a0);
		break;
	case glfnVertexAttribPointer:
		glVertexAttribPointer((GLuint)args->a0, (GLint)args->a1, (GLenum)args->a2, (GLboolean)args->a3, (GLsizei)args->a4, (const GLvoid*)args->a5);
		break;
	case glfnViewport:
		glViewport((GLint)args->a0, (GLint)args->a1, (GLint)args->a2, (GLint)args->a3);
		break;
	}
	return ret;
}
