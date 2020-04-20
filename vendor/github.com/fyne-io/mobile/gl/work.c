// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin linux openbsd

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
static void glUniformMatrix2x3fv(GLint location, GLsizei count, GLboolean transpose, const GLfloat *value) { gles3missing(); }
static void glUniformMatrix3x2fv(GLint location, GLsizei count, GLboolean transpose, const GLfloat *value) { gles3missing(); }
static void glUniformMatrix2x4fv(GLint location, GLsizei count, GLboolean transpose, const GLfloat *value) { gles3missing(); }
static void glUniformMatrix4x2fv(GLint location, GLsizei count, GLboolean transpose, const GLfloat *value) { gles3missing(); }
static void glUniformMatrix3x4fv(GLint location, GLsizei count, GLboolean transpose, const GLfloat *value) { gles3missing(); }
static void glUniformMatrix4x3fv(GLint location, GLsizei count, GLboolean transpose, const GLfloat *value) { gles3missing(); }
static void glBlitFramebuffer(GLint srcX0, GLint srcY0, GLint srcX1, GLint srcY1, GLint dstX0, GLint dstY0, GLint dstX1, GLint dstY1, GLbitfield mask, GLenum filter) { gles3missing(); }
static void glUniform1ui(GLint location, GLuint v0) { gles3missing(); }
static void glUniform2ui(GLint location, GLuint v0, GLuint v1) { gles3missing(); }
static void glUniform3ui(GLint location, GLuint v0, GLuint v1, GLuint v2) { gles3missing(); }
static void glUniform4ui(GLint location, GLuint v0, GLuint v1, GLuint v2, GLuint v3) { gles3missing(); }
static void glUniform1uiv(GLint location, GLsizei count, const GLuint *value) { gles3missing(); }
static void glUniform2uiv(GLint location, GLsizei count, const GLuint *value) { gles3missing(); }
static void glUniform3uiv(GLint location, GLsizei count, const GLuint *value) { gles3missing(); }
static void glUniform4uiv(GLint location, GLsizei count, const GLuint *value) { gles3missing(); }
static void glBindVertexArray(GLuint array) { gles3missing(); }
static void glGenVertexArrays(GLsizei n, GLuint *arrays) { gles3missing(); }
static void glDeleteVertexArrays(GLsizei n, const GLuint *arrays) { gles3missing(); }
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
	case glfnBindAttribLocation:
		glBindAttribLocation((GLint)args->a0, (GLint)args->a1, (GLchar*)args->a2);
		break;
	case glfnBindBuffer:
		glBindBuffer((GLenum)args->a0, (GLuint)args->a1);
		break;
	case glfnBindFramebuffer:
		glBindFramebuffer((GLenum)args->a0, (GLint)args->a1);
		break;
	case glfnBindRenderbuffer:
		glBindRenderbuffer((GLenum)args->a0, (GLint)args->a1);
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
	case glfnBlendEquation:
		glBlendEquation((GLenum)args->a0);
		break;
	case glfnBlendEquationSeparate:
		glBlendEquationSeparate((GLenum)args->a0, (GLenum)args->a1);
		break;
	case glfnBlendFunc:
		glBlendFunc((GLenum)args->a0, (GLenum)args->a1);
		break;
	case glfnBlendFuncSeparate:
		glBlendFuncSeparate((GLenum)args->a0, (GLenum)args->a1, (GLenum)args->a2, (GLenum)args->a3);
		break;
	case glfnBlitFramebuffer:
		glBlitFramebuffer((GLint)args->a0, (GLint)args->a1, (GLint)args->a2, (GLint)args->a3, (GLint)args->a4, (GLint)args->a5, (GLint)args->a6, (GLint)args->a7, (GLbitfield)args->a8, (GLenum)args->a9);
		break;
	case glfnBufferData:
		glBufferData((GLenum)args->a0, (GLsizeiptr)args->a1, (GLvoid*)parg, (GLenum)args->a2);
		break;
	case glfnBufferSubData:
		glBufferSubData((GLenum)args->a0, (GLint)args->a1, (GLsizeiptr)args->a2, (GLvoid*)parg);
		break;
	case glfnCheckFramebufferStatus:
		ret = glCheckFramebufferStatus((GLenum)args->a0);
		break;
	case glfnClear:
		glClear((GLenum)args->a0);
		break;
	case glfnClearColor:
		glClearColor(*(GLfloat*)&args->a0, *(GLfloat*)&args->a1, *(GLfloat*)&args->a2, *(GLfloat*)&args->a3);
		break;
	case glfnClearDepthf:
		glClearDepthf(*(GLfloat*)&args->a0);
		break;
	case glfnClearStencil:
		glClearStencil((GLint)args->a0);
		break;
	case glfnColorMask:
		glColorMask((GLboolean)args->a0, (GLboolean)args->a1, (GLboolean)args->a2, (GLboolean)args->a3);
		break;
	case glfnCompileShader:
		glCompileShader((GLint)args->a0);
		break;
	case glfnCompressedTexImage2D:
		glCompressedTexImage2D((GLenum)args->a0, (GLint)args->a1, (GLenum)args->a2, (GLint)args->a3, (GLint)args->a4, (GLint)args->a5, (GLsizeiptr)args->a6, (GLvoid*)parg);
		break;
	case glfnCompressedTexSubImage2D:
		glCompressedTexSubImage2D((GLenum)args->a0, (GLint)args->a1, (GLint)args->a2, (GLint)args->a3, (GLint)args->a4, (GLint)args->a5, (GLenum)args->a6, (GLsizeiptr)args->a7, (GLvoid*)parg);
		break;
	case glfnCopyTexImage2D:
		glCopyTexImage2D((GLenum)args->a0, (GLint)args->a1, (GLenum)args->a2, (GLint)args->a3, (GLint)args->a4, (GLint)args->a5, (GLint)args->a6, (GLint)args->a7);
		break;
	case glfnCopyTexSubImage2D:
		glCopyTexSubImage2D((GLenum)args->a0, (GLint)args->a1, (GLint)args->a2, (GLint)args->a3, (GLint)args->a4, (GLint)args->a5, (GLint)args->a6, (GLint)args->a7);
		break;
	case glfnCreateProgram:
		ret = glCreateProgram();
		break;
	case glfnCreateShader:
		ret = glCreateShader((GLenum)args->a0);
		break;
	case glfnCullFace:
		glCullFace((GLenum)args->a0);
		break;
	case glfnDeleteBuffer:
		glDeleteBuffers(1, (const GLuint*)(&args->a0));
		break;
	case glfnDeleteFramebuffer:
		glDeleteFramebuffers(1, (const GLuint*)(&args->a0));
		break;
	case glfnDeleteProgram:
		glDeleteProgram((GLint)args->a0);
		break;
	case glfnDeleteRenderbuffer:
		glDeleteRenderbuffers(1, (const GLuint*)(&args->a0));
		break;
	case glfnDeleteShader:
		glDeleteShader((GLint)args->a0);
		break;
	case glfnDeleteTexture:
		glDeleteTextures(1, (const GLuint*)(&args->a0));
		break;
	case glfnDeleteVertexArray:
		glDeleteVertexArrays(1, (const GLuint*)(&args->a0));
		break;
	case glfnDepthFunc:
		glDepthFunc((GLenum)args->a0);
		break;
	case glfnDepthMask:
		glDepthMask((GLboolean)args->a0);
		break;
	case glfnDepthRangef:
		glDepthRangef(*(GLfloat*)&args->a0, *(GLfloat*)&args->a1);
		break;
	case glfnDetachShader:
		glDetachShader((GLint)args->a0, (GLint)args->a1);
		break;
	case glfnDisable:
		glDisable((GLenum)args->a0);
		break;
	case glfnDisableVertexAttribArray:
		glDisableVertexAttribArray((GLint)args->a0);
		break;
	case glfnDrawArrays:
		glDrawArrays((GLenum)args->a0, (GLint)args->a1, (GLint)args->a2);
		break;
	case glfnDrawElements:
		glDrawElements((GLenum)args->a0, (GLint)args->a1, (GLenum)args->a2, (void*)args->a3);
		break;
	case glfnEnable:
		glEnable((GLenum)args->a0);
		break;
	case glfnEnableVertexAttribArray:
		glEnableVertexAttribArray((GLint)args->a0);
		break;
	case glfnFinish:
		glFinish();
		break;
	case glfnFlush:
		glFlush();
		break;
	case glfnFramebufferRenderbuffer:
		glFramebufferRenderbuffer((GLenum)args->a0, (GLenum)args->a1, (GLenum)args->a2, (GLint)args->a3);
		break;
	case glfnFramebufferTexture2D:
		glFramebufferTexture2D((GLenum)args->a0, (GLenum)args->a1, (GLenum)args->a2, (GLint)args->a3, (GLint)args->a4);
		break;
	case glfnFrontFace:
		glFrontFace((GLenum)args->a0);
		break;
	case glfnGenBuffer:
		glGenBuffers(1, (GLuint*)&ret);
		break;
	case glfnGenFramebuffer:
		glGenFramebuffers(1, (GLuint*)&ret);
		break;
	case glfnGenRenderbuffer:
		glGenRenderbuffers(1, (GLuint*)&ret);
		break;
	case glfnGenTexture:
		glGenTextures(1, (GLuint*)&ret);
		break;
	case glfnGenVertexArray:
		glGenVertexArrays(1, (GLuint*)&ret);
		break;
	case glfnGenerateMipmap:
		glGenerateMipmap((GLenum)args->a0);
		break;
	case glfnGetActiveAttrib:
		glGetActiveAttrib(
			(GLuint)args->a0,
			(GLuint)args->a1,
			(GLsizei)args->a2,
			NULL,
			(GLint*)&ret,
			(GLenum*)args->a3,
			(GLchar*)parg);
		break;
	case glfnGetActiveUniform:
		glGetActiveUniform(
			(GLuint)args->a0,
			(GLuint)args->a1,
			(GLsizei)args->a2,
			NULL,
			(GLint*)&ret,
			(GLenum*)args->a3,
			(GLchar*)parg);
		break;
	case glfnGetAttachedShaders:
		glGetAttachedShaders((GLuint)args->a0, (GLsizei)args->a1, (GLsizei*)&ret, (GLuint*)parg);
		break;
	case glfnGetAttribLocation:
		ret = glGetAttribLocation((GLint)args->a0, (GLchar*)args->a1);
		break;
	case glfnGetBooleanv:
		glGetBooleanv((GLenum)args->a0, (GLboolean*)parg);
		break;
	case glfnGetBufferParameteri:
		glGetBufferParameteriv((GLenum)args->a0, (GLenum)args->a1, (GLint*)&ret);
		break;
	case glfnGetFloatv:
		glGetFloatv((GLenum)args->a0, (GLfloat*)parg);
		break;
	case glfnGetIntegerv:
		glGetIntegerv((GLenum)args->a0, (GLint*)parg);
		break;
	case glfnGetError:
		ret = glGetError();
		break;
	case glfnGetFramebufferAttachmentParameteriv:
		glGetFramebufferAttachmentParameteriv((GLenum)args->a0, (GLenum)args->a1, (GLenum)args->a2, (GLint*)&ret);
		break;
	case glfnGetProgramiv:
		glGetProgramiv((GLint)args->a0, (GLenum)args->a1, (GLint*)&ret);
		break;
	case glfnGetProgramInfoLog:
		glGetProgramInfoLog((GLuint)args->a0, (GLsizei)args->a1, 0, (GLchar*)parg);
		break;
	case glfnGetRenderbufferParameteriv:
		glGetRenderbufferParameteriv((GLenum)args->a0, (GLenum)args->a1, (GLint*)&ret);
		break;
	case glfnGetShaderiv:
		glGetShaderiv((GLint)args->a0, (GLenum)args->a1, (GLint*)&ret);
		break;
	case glfnGetShaderInfoLog:
		glGetShaderInfoLog((GLuint)args->a0, (GLsizei)args->a1, 0, (GLchar*)parg);
		break;
	case glfnGetShaderPrecisionFormat:
		glGetShaderPrecisionFormat((GLenum)args->a0, (GLenum)args->a1, (GLint*)parg, &((GLint*)parg)[2]);
		break;
	case glfnGetShaderSource:
		glGetShaderSource((GLuint)args->a0, (GLsizei)args->a1, 0, (GLchar*)parg);
		break;
	case glfnGetString:
		ret = (uintptr_t)glGetString((GLenum)args->a0);
		break;
	case glfnGetTexParameterfv:
		glGetTexParameterfv((GLenum)args->a0, (GLenum)args->a1, (GLfloat*)parg);
		break;
	case glfnGetTexParameteriv:
		glGetTexParameteriv((GLenum)args->a0, (GLenum)args->a1, (GLint*)parg);
		break;
	case glfnGetUniformfv:
		glGetUniformfv((GLuint)args->a0, (GLint)args->a1, (GLfloat*)parg);
		break;
	case glfnGetUniformiv:
		glGetUniformiv((GLuint)args->a0, (GLint)args->a1, (GLint*)parg);
		break;
	case glfnGetUniformLocation:
		ret = glGetUniformLocation((GLint)args->a0, (GLchar*)args->a1);
		break;
	case glfnGetVertexAttribfv:
		glGetVertexAttribfv((GLuint)args->a0, (GLenum)args->a1, (GLfloat*)parg);
		break;
	case glfnGetVertexAttribiv:
		glGetVertexAttribiv((GLuint)args->a0, (GLenum)args->a1, (GLint*)parg);
		break;
	case glfnHint:
		glHint((GLenum)args->a0, (GLenum)args->a1);
		break;
	case glfnIsBuffer:
		ret = glIsBuffer((GLint)args->a0);
		break;
	case glfnIsEnabled:
		ret = glIsEnabled((GLenum)args->a0);
		break;
	case glfnIsFramebuffer:
		ret = glIsFramebuffer((GLint)args->a0);
		break;
	case glfnIsProgram:
		ret = glIsProgram((GLint)args->a0);
		break;
	case glfnIsRenderbuffer:
		ret = glIsRenderbuffer((GLint)args->a0);
		break;
	case glfnIsShader:
		ret = glIsShader((GLint)args->a0);
		break;
	case glfnIsTexture:
		ret = glIsTexture((GLint)args->a0);
		break;
	case glfnLineWidth:
		glLineWidth(*(GLfloat*)&args->a0);
		break;
	case glfnLinkProgram:
		glLinkProgram((GLint)args->a0);
		break;
	case glfnPixelStorei:
		glPixelStorei((GLenum)args->a0, (GLint)args->a1);
		break;
	case glfnPolygonOffset:
		glPolygonOffset(*(GLfloat*)&args->a0, *(GLfloat*)&args->a1);
		break;
	case glfnReadPixels:
		glReadPixels((GLint)args->a0, (GLint)args->a1, (GLsizei)args->a2, (GLsizei)args->a3, (GLenum)args->a4, (GLenum)args->a5, (void*)parg);
		break;
	case glfnReleaseShaderCompiler:
		glReleaseShaderCompiler();
		break;
	case glfnRenderbufferStorage:
		glRenderbufferStorage((GLenum)args->a0, (GLenum)args->a1, (GLint)args->a2, (GLint)args->a3);
		break;
	case glfnSampleCoverage:
		glSampleCoverage(*(GLfloat*)&args->a0, (GLboolean)args->a1);
		break;
	case glfnScissor:
		glScissor((GLint)args->a0, (GLint)args->a1, (GLint)args->a2, (GLint)args->a3);
		break;
	case glfnShaderSource:
#if defined(os_ios) || defined(os_osx)
		glShaderSource((GLuint)args->a0, (GLsizei)args->a1, (const GLchar *const *)args->a2, NULL);
#else
		glShaderSource((GLuint)args->a0, (GLsizei)args->a1, (const GLchar **)args->a2, NULL);
#endif
		break;
	case glfnStencilFunc:
		glStencilFunc((GLenum)args->a0, (GLint)args->a1, (GLuint)args->a2);
		break;
	case glfnStencilFuncSeparate:
		glStencilFuncSeparate((GLenum)args->a0, (GLenum)args->a1, (GLint)args->a2, (GLuint)args->a3);
		break;
	case glfnStencilMask:
		glStencilMask((GLuint)args->a0);
		break;
	case glfnStencilMaskSeparate:
		glStencilMaskSeparate((GLenum)args->a0, (GLuint)args->a1);
		break;
	case glfnStencilOp:
		glStencilOp((GLenum)args->a0, (GLenum)args->a1, (GLenum)args->a2);
		break;
	case glfnStencilOpSeparate:
		glStencilOpSeparate((GLenum)args->a0, (GLenum)args->a1, (GLenum)args->a2, (GLenum)args->a3);
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
	case glfnTexSubImage2D:
		glTexSubImage2D(
			(GLenum)args->a0,
			(GLint)args->a1,
			(GLint)args->a2,
			(GLint)args->a3,
			(GLsizei)args->a4,
			(GLsizei)args->a5,
			(GLenum)args->a6,
			(GLenum)args->a7,
			(const GLvoid*)parg);
		break;
	case glfnTexParameterf:
		glTexParameterf((GLenum)args->a0, (GLenum)args->a1, *(GLfloat*)&args->a2);
		break;
	case glfnTexParameterfv:
		glTexParameterfv((GLenum)args->a0, (GLenum)args->a1, (GLfloat*)parg);
		break;
	case glfnTexParameteri:
		glTexParameteri((GLenum)args->a0, (GLenum)args->a1, (GLint)args->a2);
		break;
	case glfnTexParameteriv:
		glTexParameteriv((GLenum)args->a0, (GLenum)args->a1, (GLint*)parg);
		break;
	case glfnUniform1f:
		glUniform1f((GLint)args->a0, *(GLfloat*)&args->a1);
		break;
	case glfnUniform1fv:
		glUniform1fv((GLint)args->a0, (GLsizeiptr)args->a1, (GLvoid*)parg);
		break;
	case glfnUniform1i:
		glUniform1i((GLint)args->a0, (GLint)args->a1);
		break;
	case glfnUniform1ui:
		glUniform1ui((GLint)args->a0, (GLuint)args->a1);
		break;
	case glfnUniform1uiv:
		glUniform1uiv((GLint)args->a0, (GLsizeiptr)args->a1, (GLuint*)parg);
		break;
	case glfnUniform1iv:
		glUniform1iv((GLint)args->a0, (GLsizeiptr)args->a1, (GLvoid*)parg);
		break;
	case glfnUniform2f:
		glUniform2f((GLint)args->a0, *(GLfloat*)&args->a1, *(GLfloat*)&args->a2);
		break;
	case glfnUniform2fv:
		glUniform2fv((GLint)args->a0, (GLsizeiptr)args->a1, (GLvoid*)parg);
		break;
	case glfnUniform2i:
		glUniform2i((GLint)args->a0, (GLint)args->a1, (GLint)args->a2);
		break;
	case glfnUniform2ui:
		glUniform2ui((GLint)args->a0, (GLuint)args->a1, (GLuint)args->a2);
		break;
	case glfnUniform2uiv:
		glUniform2uiv((GLint)args->a0, (GLsizeiptr)args->a1, (GLuint*)parg);
		break;
	case glfnUniform2iv:
		glUniform2iv((GLint)args->a0, (GLsizeiptr)args->a1, (GLvoid*)parg);
		break;
	case glfnUniform3f:
		glUniform3f((GLint)args->a0, *(GLfloat*)&args->a1, *(GLfloat*)&args->a2, *(GLfloat*)&args->a3);
		break;
	case glfnUniform3fv:
		glUniform3fv((GLint)args->a0, (GLsizeiptr)args->a1, (GLvoid*)parg);
		break;
	case glfnUniform3i:
		glUniform3i((GLint)args->a0, (GLint)args->a1, (GLint)args->a2, (GLint)args->a3);
		break;
	case glfnUniform3ui:
		glUniform3ui((GLint)args->a0, (GLuint)args->a1, (GLuint)args->a2, (GLuint)args->a3);
		break;
	case glfnUniform3uiv:
		glUniform3uiv((GLint)args->a0, (GLsizeiptr)args->a1, (GLuint*)parg);
		break;
	case glfnUniform3iv:
		glUniform3iv((GLint)args->a0, (GLsizeiptr)args->a1, (GLvoid*)parg);
		break;
	case glfnUniform4f:
		glUniform4f((GLint)args->a0, *(GLfloat*)&args->a1, *(GLfloat*)&args->a2, *(GLfloat*)&args->a3, *(GLfloat*)&args->a4);
		break;
	case glfnUniform4fv:
		glUniform4fv((GLint)args->a0, (GLsizeiptr)args->a1, (GLvoid*)parg);
		break;
	case glfnUniform4i:
		glUniform4i((GLint)args->a0, (GLint)args->a1, (GLint)args->a2, (GLint)args->a3, (GLint)args->a4);
		break;
	case glfnUniform4ui:
		glUniform4ui((GLint)args->a0, (GLuint)args->a1, (GLuint)args->a2, (GLuint)args->a3, (GLuint)args->a4);
		break;
	case glfnUniform4uiv:
		glUniform4uiv((GLint)args->a0, (GLsizeiptr)args->a1, (GLuint*)parg);
		break;
	case glfnUniform4iv:
		glUniform4iv((GLint)args->a0, (GLsizeiptr)args->a1, (GLvoid*)parg);
		break;
	case glfnUniformMatrix2fv:
		glUniformMatrix2fv((GLint)args->a0, (GLsizeiptr)args->a1, 0, (GLvoid*)parg);
		break;
	case glfnUniformMatrix3fv:
		glUniformMatrix3fv((GLint)args->a0, (GLsizeiptr)args->a1, 0, (GLvoid*)parg);
		break;
	case glfnUniformMatrix4fv:
		glUniformMatrix4fv((GLint)args->a0, (GLsizeiptr)args->a1, 0, (GLvoid*)parg);
		break;
	case glfnUniformMatrix2x3fv:
		glUniformMatrix2x3fv((GLint)args->a0, (GLsizeiptr)args->a1, 0, (GLvoid*)parg);
		break;
	case glfnUniformMatrix3x2fv:
		glUniformMatrix3x2fv((GLint)args->a0, (GLsizeiptr)args->a1, 0, (GLvoid*)parg);
		break;
	case glfnUniformMatrix2x4fv:
		glUniformMatrix2x4fv((GLint)args->a0, (GLsizeiptr)args->a1, 0, (GLvoid*)parg);
		break;
	case glfnUniformMatrix4x2fv:
		glUniformMatrix4x2fv((GLint)args->a0, (GLsizeiptr)args->a1, 0, (GLvoid*)parg);
		break;
	case glfnUniformMatrix3x4fv:
		glUniformMatrix3x4fv((GLint)args->a0, (GLsizeiptr)args->a1, 0, (GLvoid*)parg);
		break;
	case glfnUniformMatrix4x3fv:
		glUniformMatrix4x3fv((GLint)args->a0, (GLsizeiptr)args->a1, 0, (GLvoid*)parg);
		break;
	case glfnUseProgram:
		glUseProgram((GLint)args->a0);
		break;
	case glfnValidateProgram:
		glValidateProgram((GLint)args->a0);
		break;
	case glfnVertexAttrib1f:
		glVertexAttrib1f((GLint)args->a0, *(GLfloat*)&args->a1);
		break;
	case glfnVertexAttrib1fv:
		glVertexAttrib1fv((GLint)args->a0, (GLfloat*)parg);
		break;
	case glfnVertexAttrib2f:
		glVertexAttrib2f((GLint)args->a0, *(GLfloat*)&args->a1, *(GLfloat*)&args->a2);
		break;
	case glfnVertexAttrib2fv:
		glVertexAttrib2fv((GLint)args->a0, (GLfloat*)parg);
		break;
	case glfnVertexAttrib3f:
		glVertexAttrib3f((GLint)args->a0, *(GLfloat*)&args->a1, *(GLfloat*)&args->a2, *(GLfloat*)&args->a3);
		break;
	case glfnVertexAttrib3fv:
		glVertexAttrib3fv((GLint)args->a0, (GLfloat*)parg);
		break;
	case glfnVertexAttrib4f:
		glVertexAttrib4f((GLint)args->a0, *(GLfloat*)&args->a1, *(GLfloat*)&args->a2, *(GLfloat*)&args->a3, *(GLfloat*)&args->a4);
		break;
	case glfnVertexAttrib4fv:
		glVertexAttrib4fv((GLint)args->a0, (GLfloat*)parg);
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
