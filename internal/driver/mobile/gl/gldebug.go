// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated from gl.go using go generate. DO NOT EDIT.
// See doc.go for details.

// +build darwin linux openbsd freebsd windows
// +build gldebug

package gl

import (
	"fmt"
	"log"
	"math"
	"sync/atomic"
	"unsafe"
)

func (ctx *context) errDrain() string {
	var errs []Enum
	for {
		e := ctx.GetError()
		if e == 0 {
			break
		}
		errs = append(errs, e)
	}
	if len(errs) > 0 {
		return fmt.Sprintf(" error: %v", errs)
	}
	return ""
}

func (ctx *context) enqueueDebug(c call) uintptr {
	numCalls := atomic.AddInt32(&ctx.debug, 1)
	if numCalls > 1 {
		panic("concurrent calls made to the same GL context")
	}
	defer func() {
		if atomic.AddInt32(&ctx.debug, -1) > 0 {
			select {} // block so you see us in the panic
		}
	}()

	return ctx.enqueue(c)
}

var enumString = map[Enum]string{
	0x0:        "0",
	0x1:        "1",
	0x2:        "2",
	0x3:        "LINE_STRIP",
	0x4:        "4",
	0x5:        "TRIANGLE_STRIP",
	0x6:        "TRIANGLE_FAN",
	0x300:      "SRC_COLOR",
	0x301:      "ONE_MINUS_SRC_COLOR",
	0x302:      "SRC_ALPHA",
	0x303:      "ONE_MINUS_SRC_ALPHA",
	0x304:      "DST_ALPHA",
	0x305:      "ONE_MINUS_DST_ALPHA",
	0x306:      "DST_COLOR",
	0x307:      "ONE_MINUS_DST_COLOR",
	0x308:      "SRC_ALPHA_SATURATE",
	0x8006:     "FUNC_ADD",
	0x8009:     "32777",
	0x883d:     "BLEND_EQUATION_ALPHA",
	0x800a:     "FUNC_SUBTRACT",
	0x800b:     "FUNC_REVERSE_SUBTRACT",
	0x80c8:     "BLEND_DST_RGB",
	0x80c9:     "BLEND_SRC_RGB",
	0x80ca:     "BLEND_DST_ALPHA",
	0x80cb:     "BLEND_SRC_ALPHA",
	0x8001:     "CONSTANT_COLOR",
	0x8002:     "ONE_MINUS_CONSTANT_COLOR",
	0x8003:     "CONSTANT_ALPHA",
	0x8004:     "ONE_MINUS_CONSTANT_ALPHA",
	0x8005:     "BLEND_COLOR",
	0x8892:     "ARRAY_BUFFER",
	0x8893:     "ELEMENT_ARRAY_BUFFER",
	0x8894:     "ARRAY_BUFFER_BINDING",
	0x8895:     "ELEMENT_ARRAY_BUFFER_BINDING",
	0x88e0:     "STREAM_DRAW",
	0x88e4:     "STATIC_DRAW",
	0x88e8:     "DYNAMIC_DRAW",
	0x8764:     "BUFFER_SIZE",
	0x8765:     "BUFFER_USAGE",
	0x8626:     "CURRENT_VERTEX_ATTRIB",
	0x404:      "FRONT",
	0x405:      "BACK",
	0x408:      "FRONT_AND_BACK",
	0xde1:      "TEXTURE_2D",
	0xb44:      "CULL_FACE",
	0xbe2:      "BLEND",
	0xbd0:      "DITHER",
	0xb90:      "STENCIL_TEST",
	0xb71:      "DEPTH_TEST",
	0xc11:      "SCISSOR_TEST",
	0x8037:     "POLYGON_OFFSET_FILL",
	0x809e:     "SAMPLE_ALPHA_TO_COVERAGE",
	0x80a0:     "SAMPLE_COVERAGE",
	0x500:      "INVALID_ENUM",
	0x501:      "INVALID_VALUE",
	0x502:      "INVALID_OPERATION",
	0x505:      "OUT_OF_MEMORY",
	0x900:      "CW",
	0x901:      "CCW",
	0xb21:      "LINE_WIDTH",
	0x846d:     "ALIASED_POINT_SIZE_RANGE",
	0x846e:     "ALIASED_LINE_WIDTH_RANGE",
	0xb45:      "CULL_FACE_MODE",
	0xb46:      "FRONT_FACE",
	0xb70:      "DEPTH_RANGE",
	0xb72:      "DEPTH_WRITEMASK",
	0xb73:      "DEPTH_CLEAR_VALUE",
	0xb74:      "DEPTH_FUNC",
	0xb91:      "STENCIL_CLEAR_VALUE",
	0xb92:      "STENCIL_FUNC",
	0xb94:      "STENCIL_FAIL",
	0xb95:      "STENCIL_PASS_DEPTH_FAIL",
	0xb96:      "STENCIL_PASS_DEPTH_PASS",
	0xb97:      "STENCIL_REF",
	0xb93:      "STENCIL_VALUE_MASK",
	0xb98:      "STENCIL_WRITEMASK",
	0x8800:     "STENCIL_BACK_FUNC",
	0x8801:     "STENCIL_BACK_FAIL",
	0x8802:     "STENCIL_BACK_PASS_DEPTH_FAIL",
	0x8803:     "STENCIL_BACK_PASS_DEPTH_PASS",
	0x8ca3:     "STENCIL_BACK_REF",
	0x8ca4:     "STENCIL_BACK_VALUE_MASK",
	0x8ca5:     "STENCIL_BACK_WRITEMASK",
	0xba2:      "VIEWPORT",
	0xc10:      "SCISSOR_BOX",
	0xc22:      "COLOR_CLEAR_VALUE",
	0xc23:      "COLOR_WRITEMASK",
	0xcf5:      "UNPACK_ALIGNMENT",
	0xd05:      "PACK_ALIGNMENT",
	0xd33:      "MAX_TEXTURE_SIZE",
	0xd3a:      "MAX_VIEWPORT_DIMS",
	0xd50:      "SUBPIXEL_BITS",
	0xd52:      "RED_BITS",
	0xd53:      "GREEN_BITS",
	0xd54:      "BLUE_BITS",
	0xd55:      "ALPHA_BITS",
	0xd56:      "DEPTH_BITS",
	0xd57:      "STENCIL_BITS",
	0x2a00:     "POLYGON_OFFSET_UNITS",
	0x8038:     "POLYGON_OFFSET_FACTOR",
	0x8069:     "TEXTURE_BINDING_2D",
	0x80a8:     "SAMPLE_BUFFERS",
	0x80a9:     "SAMPLES",
	0x80aa:     "SAMPLE_COVERAGE_VALUE",
	0x80ab:     "SAMPLE_COVERAGE_INVERT",
	0x86a2:     "NUM_COMPRESSED_TEXTURE_FORMATS",
	0x86a3:     "COMPRESSED_TEXTURE_FORMATS",
	0x1100:     "DONT_CARE",
	0x1101:     "FASTEST",
	0x1102:     "NICEST",
	0x8192:     "GENERATE_MIPMAP_HINT",
	0x1400:     "BYTE",
	0x1401:     "UNSIGNED_BYTE",
	0x1402:     "SHORT",
	0x1403:     "UNSIGNED_SHORT",
	0x1404:     "INT",
	0x1405:     "UNSIGNED_INT",
	0x1406:     "FLOAT",
	0x140c:     "FIXED",
	0x1902:     "DEPTH_COMPONENT",
	0x1906:     "ALPHA",
	0x1907:     "RGB",
	0x1908:     "RGBA",
	0x1909:     "LUMINANCE",
	0x190a:     "LUMINANCE_ALPHA",
	0x8033:     "UNSIGNED_SHORT_4_4_4_4",
	0x8034:     "UNSIGNED_SHORT_5_5_5_1",
	0x8363:     "UNSIGNED_SHORT_5_6_5",
	0x8869:     "MAX_VERTEX_ATTRIBS",
	0x8dfb:     "MAX_VERTEX_UNIFORM_VECTORS",
	0x8dfc:     "MAX_VARYING_VECTORS",
	0x8b4d:     "MAX_COMBINED_TEXTURE_IMAGE_UNITS",
	0x8b4c:     "MAX_VERTEX_TEXTURE_IMAGE_UNITS",
	0x8872:     "MAX_TEXTURE_IMAGE_UNITS",
	0x8dfd:     "MAX_FRAGMENT_UNIFORM_VECTORS",
	0x8b4f:     "SHADER_TYPE",
	0x8b80:     "DELETE_STATUS",
	0x8b82:     "LINK_STATUS",
	0x8b83:     "VALIDATE_STATUS",
	0x8b85:     "ATTACHED_SHADERS",
	0x8b86:     "ACTIVE_UNIFORMS",
	0x8b87:     "ACTIVE_UNIFORM_MAX_LENGTH",
	0x8b89:     "ACTIVE_ATTRIBUTES",
	0x8b8a:     "ACTIVE_ATTRIBUTE_MAX_LENGTH",
	0x8b8c:     "SHADING_LANGUAGE_VERSION",
	0x8b8d:     "CURRENT_PROGRAM",
	0x200:      "NEVER",
	0x201:      "LESS",
	0x202:      "EQUAL",
	0x203:      "LEQUAL",
	0x204:      "GREATER",
	0x205:      "NOTEQUAL",
	0x206:      "GEQUAL",
	0x207:      "ALWAYS",
	0x1e00:     "KEEP",
	0x1e01:     "REPLACE",
	0x1e02:     "INCR",
	0x1e03:     "DECR",
	0x150a:     "INVERT",
	0x8507:     "INCR_WRAP",
	0x8508:     "DECR_WRAP",
	0x1f00:     "VENDOR",
	0x1f01:     "RENDERER",
	0x1f02:     "VERSION",
	0x1f03:     "EXTENSIONS",
	0x2600:     "NEAREST",
	0x2601:     "LINEAR",
	0x2700:     "NEAREST_MIPMAP_NEAREST",
	0x2701:     "LINEAR_MIPMAP_NEAREST",
	0x2702:     "NEAREST_MIPMAP_LINEAR",
	0x2703:     "LINEAR_MIPMAP_LINEAR",
	0x2800:     "TEXTURE_MAG_FILTER",
	0x2801:     "TEXTURE_MIN_FILTER",
	0x2802:     "TEXTURE_WRAP_S",
	0x2803:     "TEXTURE_WRAP_T",
	0x1702:     "TEXTURE",
	0x8513:     "TEXTURE_CUBE_MAP",
	0x8514:     "TEXTURE_BINDING_CUBE_MAP",
	0x8515:     "TEXTURE_CUBE_MAP_POSITIVE_X",
	0x8516:     "TEXTURE_CUBE_MAP_NEGATIVE_X",
	0x8517:     "TEXTURE_CUBE_MAP_POSITIVE_Y",
	0x8518:     "TEXTURE_CUBE_MAP_NEGATIVE_Y",
	0x8519:     "TEXTURE_CUBE_MAP_POSITIVE_Z",
	0x851a:     "TEXTURE_CUBE_MAP_NEGATIVE_Z",
	0x851c:     "MAX_CUBE_MAP_TEXTURE_SIZE",
	0x84c0:     "TEXTURE0",
	0x84c1:     "TEXTURE1",
	0x84c2:     "TEXTURE2",
	0x84c3:     "TEXTURE3",
	0x84c4:     "TEXTURE4",
	0x84c5:     "TEXTURE5",
	0x84c6:     "TEXTURE6",
	0x84c7:     "TEXTURE7",
	0x84c8:     "TEXTURE8",
	0x84c9:     "TEXTURE9",
	0x84ca:     "TEXTURE10",
	0x84cb:     "TEXTURE11",
	0x84cc:     "TEXTURE12",
	0x84cd:     "TEXTURE13",
	0x84ce:     "TEXTURE14",
	0x84cf:     "TEXTURE15",
	0x84d0:     "TEXTURE16",
	0x84d1:     "TEXTURE17",
	0x84d2:     "TEXTURE18",
	0x84d3:     "TEXTURE19",
	0x84d4:     "TEXTURE20",
	0x84d5:     "TEXTURE21",
	0x84d6:     "TEXTURE22",
	0x84d7:     "TEXTURE23",
	0x84d8:     "TEXTURE24",
	0x84d9:     "TEXTURE25",
	0x84da:     "TEXTURE26",
	0x84db:     "TEXTURE27",
	0x84dc:     "TEXTURE28",
	0x84dd:     "TEXTURE29",
	0x84de:     "TEXTURE30",
	0x84df:     "TEXTURE31",
	0x84e0:     "ACTIVE_TEXTURE",
	0x2901:     "REPEAT",
	0x812f:     "CLAMP_TO_EDGE",
	0x8370:     "MIRRORED_REPEAT",
	0x8622:     "VERTEX_ATTRIB_ARRAY_ENABLED",
	0x8623:     "VERTEX_ATTRIB_ARRAY_SIZE",
	0x8624:     "VERTEX_ATTRIB_ARRAY_STRIDE",
	0x8625:     "VERTEX_ATTRIB_ARRAY_TYPE",
	0x886a:     "VERTEX_ATTRIB_ARRAY_NORMALIZED",
	0x8645:     "VERTEX_ATTRIB_ARRAY_POINTER",
	0x889f:     "VERTEX_ATTRIB_ARRAY_BUFFER_BINDING",
	0x8b9a:     "IMPLEMENTATION_COLOR_READ_TYPE",
	0x8b9b:     "IMPLEMENTATION_COLOR_READ_FORMAT",
	0x8b81:     "COMPILE_STATUS",
	0x8b84:     "INFO_LOG_LENGTH",
	0x8b88:     "SHADER_SOURCE_LENGTH",
	0x8dfa:     "SHADER_COMPILER",
	0x8df8:     "SHADER_BINARY_FORMATS",
	0x8df9:     "NUM_SHADER_BINARY_FORMATS",
	0x8df0:     "LOW_FLOAT",
	0x8df1:     "MEDIUM_FLOAT",
	0x8df2:     "HIGH_FLOAT",
	0x8df3:     "LOW_INT",
	0x8df4:     "MEDIUM_INT",
	0x8df5:     "HIGH_INT",
	0x8d40:     "FRAMEBUFFER",
	0x8d41:     "RENDERBUFFER",
	0x8056:     "RGBA4",
	0x8057:     "RGB5_A1",
	0x8d62:     "RGB565",
	0x81a5:     "DEPTH_COMPONENT16",
	0x8d48:     "STENCIL_INDEX8",
	0x8d42:     "RENDERBUFFER_WIDTH",
	0x8d43:     "RENDERBUFFER_HEIGHT",
	0x8d44:     "RENDERBUFFER_INTERNAL_FORMAT",
	0x8d50:     "RENDERBUFFER_RED_SIZE",
	0x8d51:     "RENDERBUFFER_GREEN_SIZE",
	0x8d52:     "RENDERBUFFER_BLUE_SIZE",
	0x8d53:     "RENDERBUFFER_ALPHA_SIZE",
	0x8d54:     "RENDERBUFFER_DEPTH_SIZE",
	0x8d55:     "RENDERBUFFER_STENCIL_SIZE",
	0x8cd0:     "FRAMEBUFFER_ATTACHMENT_OBJECT_TYPE",
	0x8cd1:     "FRAMEBUFFER_ATTACHMENT_OBJECT_NAME",
	0x8cd2:     "FRAMEBUFFER_ATTACHMENT_TEXTURE_LEVEL",
	0x8cd3:     "FRAMEBUFFER_ATTACHMENT_TEXTURE_CUBE_MAP_FACE",
	0x8ce0:     "COLOR_ATTACHMENT0",
	0x8d00:     "DEPTH_ATTACHMENT",
	0x8d20:     "STENCIL_ATTACHMENT",
	0x8cd5:     "FRAMEBUFFER_COMPLETE",
	0x8cd6:     "FRAMEBUFFER_INCOMPLETE_ATTACHMENT",
	0x8cd7:     "FRAMEBUFFER_INCOMPLETE_MISSING_ATTACHMENT",
	0x8cd9:     "FRAMEBUFFER_INCOMPLETE_DIMENSIONS",
	0x8cdd:     "FRAMEBUFFER_UNSUPPORTED",
	0x8ca6:     "36006",
	0x8ca7:     "RENDERBUFFER_BINDING",
	0x84e8:     "MAX_RENDERBUFFER_SIZE",
	0x506:      "INVALID_FRAMEBUFFER_OPERATION",
	0x100:      "DEPTH_BUFFER_BIT",
	0x400:      "STENCIL_BUFFER_BIT",
	0x4000:     "COLOR_BUFFER_BIT",
	0x8b50:     "FLOAT_VEC2",
	0x8b51:     "FLOAT_VEC3",
	0x8b52:     "FLOAT_VEC4",
	0x8b53:     "INT_VEC2",
	0x8b54:     "INT_VEC3",
	0x8b55:     "INT_VEC4",
	0x8b56:     "BOOL",
	0x8b57:     "BOOL_VEC2",
	0x8b58:     "BOOL_VEC3",
	0x8b59:     "BOOL_VEC4",
	0x8b5a:     "FLOAT_MAT2",
	0x8b5b:     "FLOAT_MAT3",
	0x8b5c:     "FLOAT_MAT4",
	0x8b5e:     "SAMPLER_2D",
	0x8b60:     "SAMPLER_CUBE",
	0x8b30:     "FRAGMENT_SHADER",
	0x8b31:     "VERTEX_SHADER",
	0x8a35:     "ACTIVE_UNIFORM_BLOCK_MAX_NAME_LENGTH",
	0x8a36:     "ACTIVE_UNIFORM_BLOCKS",
	0x911a:     "ALREADY_SIGNALED",
	0x8c2f:     "ANY_SAMPLES_PASSED",
	0x8d6a:     "ANY_SAMPLES_PASSED_CONSERVATIVE",
	0x1905:     "BLUE",
	0x911f:     "BUFFER_ACCESS_FLAGS",
	0x9120:     "BUFFER_MAP_LENGTH",
	0x9121:     "BUFFER_MAP_OFFSET",
	0x88bc:     "BUFFER_MAPPED",
	0x88bd:     "BUFFER_MAP_POINTER",
	0x1800:     "COLOR",
	0x8cea:     "COLOR_ATTACHMENT10",
	0x8ce1:     "COLOR_ATTACHMENT1",
	0x8ceb:     "COLOR_ATTACHMENT11",
	0x8cec:     "COLOR_ATTACHMENT12",
	0x8ced:     "COLOR_ATTACHMENT13",
	0x8cee:     "COLOR_ATTACHMENT14",
	0x8cef:     "COLOR_ATTACHMENT15",
	0x8ce2:     "COLOR_ATTACHMENT2",
	0x8ce3:     "COLOR_ATTACHMENT3",
	0x8ce4:     "COLOR_ATTACHMENT4",
	0x8ce5:     "COLOR_ATTACHMENT5",
	0x8ce6:     "COLOR_ATTACHMENT6",
	0x8ce7:     "COLOR_ATTACHMENT7",
	0x8ce8:     "COLOR_ATTACHMENT8",
	0x8ce9:     "COLOR_ATTACHMENT9",
	0x884e:     "COMPARE_REF_TO_TEXTURE",
	0x9270:     "COMPRESSED_R11_EAC",
	0x9272:     "COMPRESSED_RG11_EAC",
	0x9274:     "COMPRESSED_RGB8_ETC2",
	0x9276:     "COMPRESSED_RGB8_PUNCHTHROUGH_ALPHA1_ETC2",
	0x9278:     "COMPRESSED_RGBA8_ETC2_EAC",
	0x9271:     "COMPRESSED_SIGNED_R11_EAC",
	0x9273:     "COMPRESSED_SIGNED_RG11_EAC",
	0x9279:     "COMPRESSED_SRGB8_ALPHA8_ETC2_EAC",
	0x9275:     "COMPRESSED_SRGB8_ETC2",
	0x9277:     "COMPRESSED_SRGB8_PUNCHTHROUGH_ALPHA1_ETC2",
	0x911c:     "CONDITION_SATISFIED",
	0x8f36:     "36662",
	0x8f37:     "36663",
	0x8865:     "CURRENT_QUERY",
	0x1801:     "DEPTH",
	0x88f0:     "DEPTH24_STENCIL8",
	0x8cad:     "DEPTH32F_STENCIL8",
	0x81a6:     "DEPTH_COMPONENT24",
	0x8cac:     "DEPTH_COMPONENT32F",
	0x84f9:     "DEPTH_STENCIL",
	0x821a:     "DEPTH_STENCIL_ATTACHMENT",
	0x8825:     "DRAW_BUFFER0",
	0x882f:     "DRAW_BUFFER10",
	0x8826:     "DRAW_BUFFER1",
	0x8830:     "DRAW_BUFFER11",
	0x8831:     "DRAW_BUFFER12",
	0x8832:     "DRAW_BUFFER13",
	0x8833:     "DRAW_BUFFER14",
	0x8834:     "DRAW_BUFFER15",
	0x8827:     "DRAW_BUFFER2",
	0x8828:     "DRAW_BUFFER3",
	0x8829:     "DRAW_BUFFER4",
	0x882a:     "DRAW_BUFFER5",
	0x882b:     "DRAW_BUFFER6",
	0x882c:     "DRAW_BUFFER7",
	0x882d:     "DRAW_BUFFER8",
	0x882e:     "DRAW_BUFFER9",
	0x8ca9:     "DRAW_FRAMEBUFFER",
	0x88ea:     "DYNAMIC_COPY",
	0x88e9:     "DYNAMIC_READ",
	0x8dad:     "FLOAT_32_UNSIGNED_INT_24_8_REV",
	0x8b65:     "FLOAT_MAT2x3",
	0x8b66:     "FLOAT_MAT2x4",
	0x8b67:     "FLOAT_MAT3x2",
	0x8b68:     "FLOAT_MAT3x4",
	0x8b69:     "FLOAT_MAT4x2",
	0x8b6a:     "FLOAT_MAT4x3",
	0x8b8b:     "FRAGMENT_SHADER_DERIVATIVE_HINT",
	0x8215:     "FRAMEBUFFER_ATTACHMENT_ALPHA_SIZE",
	0x8214:     "FRAMEBUFFER_ATTACHMENT_BLUE_SIZE",
	0x8210:     "FRAMEBUFFER_ATTACHMENT_COLOR_ENCODING",
	0x8211:     "FRAMEBUFFER_ATTACHMENT_COMPONENT_TYPE",
	0x8216:     "FRAMEBUFFER_ATTACHMENT_DEPTH_SIZE",
	0x8213:     "FRAMEBUFFER_ATTACHMENT_GREEN_SIZE",
	0x8212:     "FRAMEBUFFER_ATTACHMENT_RED_SIZE",
	0x8217:     "FRAMEBUFFER_ATTACHMENT_STENCIL_SIZE",
	0x8cd4:     "FRAMEBUFFER_ATTACHMENT_TEXTURE_LAYER",
	0x8218:     "FRAMEBUFFER_DEFAULT",
	0x8d56:     "FRAMEBUFFER_INCOMPLETE_MULTISAMPLE",
	0x8219:     "FRAMEBUFFER_UNDEFINED",
	0x1904:     "GREEN",
	0x140b:     "HALF_FLOAT",
	0x8d9f:     "INT_2_10_10_10_REV",
	0x8c8c:     "INTERLEAVED_ATTRIBS",
	0x8dca:     "INT_SAMPLER_2D",
	0x8dcf:     "INT_SAMPLER_2D_ARRAY",
	0x8dcb:     "INT_SAMPLER_3D",
	0x8dcc:     "INT_SAMPLER_CUBE",
	0xffffffff: "INVALID_INDEX",
	0x821b:     "MAJOR_VERSION",
	0x10:       "MAP_FLUSH_EXPLICIT_BIT",
	0x8:        "MAP_INVALIDATE_BUFFER_BIT",
	0x20:       "MAP_UNSYNCHRONIZED_BIT",
	0x8008:     "MAX",
	0x8073:     "MAX_3D_TEXTURE_SIZE",
	0x88ff:     "MAX_ARRAY_TEXTURE_LAYERS",
	0x8cdf:     "MAX_COLOR_ATTACHMENTS",
	0x8a33:     "MAX_COMBINED_FRAGMENT_UNIFORM_COMPONENTS",
	0x8a2e:     "MAX_COMBINED_UNIFORM_BLOCKS",
	0x8a31:     "MAX_COMBINED_VERTEX_UNIFORM_COMPONENTS",
	0x8824:     "MAX_DRAW_BUFFERS",
	0x8d6b:     "MAX_ELEMENT_INDEX",
	0x80e9:     "MAX_ELEMENTS_INDICES",
	0x80e8:     "MAX_ELEMENTS_VERTICES",
	0x9125:     "MAX_FRAGMENT_INPUT_COMPONENTS",
	0x8a2d:     "MAX_FRAGMENT_UNIFORM_BLOCKS",
	0x8b49:     "MAX_FRAGMENT_UNIFORM_COMPONENTS",
	0x8905:     "MAX_PROGRAM_TEXEL_OFFSET",
	0x8d57:     "MAX_SAMPLES",
	0x9111:     "MAX_SERVER_WAIT_TIMEOUT",
	0x84fd:     "MAX_TEXTURE_LOD_BIAS",
	0x8c8a:     "MAX_TRANSFORM_FEEDBACK_INTERLEAVED_COMPONENTS",
	0x8c8b:     "MAX_TRANSFORM_FEEDBACK_SEPARATE_ATTRIBS",
	0x8c80:     "MAX_TRANSFORM_FEEDBACK_SEPARATE_COMPONENTS",
	0x8a30:     "MAX_UNIFORM_BLOCK_SIZE",
	0x8a2f:     "MAX_UNIFORM_BUFFER_BINDINGS",
	0x8b4b:     "MAX_VARYING_COMPONENTS",
	0x9122:     "MAX_VERTEX_OUTPUT_COMPONENTS",
	0x8a2b:     "MAX_VERTEX_UNIFORM_BLOCKS",
	0x8b4a:     "MAX_VERTEX_UNIFORM_COMPONENTS",
	0x8007:     "MIN",
	0x821c:     "MINOR_VERSION",
	0x8904:     "MIN_PROGRAM_TEXEL_OFFSET",
	0x821d:     "NUM_EXTENSIONS",
	0x87fe:     "NUM_PROGRAM_BINARY_FORMATS",
	0x9380:     "NUM_SAMPLE_COUNTS",
	0x9112:     "OBJECT_TYPE",
	0xd02:      "PACK_ROW_LENGTH",
	0xd04:      "PACK_SKIP_PIXELS",
	0xd03:      "PACK_SKIP_ROWS",
	0x88eb:     "PIXEL_PACK_BUFFER",
	0x88ed:     "PIXEL_PACK_BUFFER_BINDING",
	0x88ec:     "PIXEL_UNPACK_BUFFER",
	0x88ef:     "PIXEL_UNPACK_BUFFER_BINDING",
	0x8d69:     "PRIMITIVE_RESTART_FIXED_INDEX",
	0x87ff:     "PROGRAM_BINARY_FORMATS",
	0x8741:     "PROGRAM_BINARY_LENGTH",
	0x8257:     "PROGRAM_BINARY_RETRIEVABLE_HINT",
	0x8866:     "QUERY_RESULT",
	0x8867:     "QUERY_RESULT_AVAILABLE",
	0x8c3a:     "R11F_G11F_B10F",
	0x822d:     "R16F",
	0x8233:     "R16I",
	0x8234:     "R16UI",
	0x822e:     "R32F",
	0x8235:     "R32I",
	0x8236:     "R32UI",
	0x8229:     "R8",
	0x8231:     "R8I",
	0x8f94:     "R8_SNORM",
	0x8232:     "R8UI",
	0x8c89:     "RASTERIZER_DISCARD",
	0xc02:      "READ_BUFFER",
	0x8ca8:     "READ_FRAMEBUFFER",
	0x8caa:     "READ_FRAMEBUFFER_BINDING",
	0x1903:     "RED",
	0x8d94:     "RED_INTEGER",
	0x8cab:     "RENDERBUFFER_SAMPLES",
	0x8227:     "RG",
	0x822f:     "RG16F",
	0x8239:     "RG16I",
	0x823a:     "RG16UI",
	0x8230:     "RG32F",
	0x823b:     "RG32I",
	0x823c:     "RG32UI",
	0x822b:     "RG8",
	0x8237:     "RG8I",
	0x8f95:     "RG8_SNORM",
	0x8238:     "RG8UI",
	0x8059:     "RGB10_A2",
	0x906f:     "RGB10_A2UI",
	0x881b:     "RGB16F",
	0x8d89:     "RGB16I",
	0x8d77:     "RGB16UI",
	0x8815:     "RGB32F",
	0x8d83:     "RGB32I",
	0x8d71:     "RGB32UI",
	0x8051:     "RGB8",
	0x8d8f:     "RGB8I",
	0x8f96:     "RGB8_SNORM",
	0x8d7d:     "RGB8UI",
	0x8c3d:     "RGB9_E5",
	0x881a:     "RGBA16F",
	0x8d88:     "RGBA16I",
	0x8d76:     "RGBA16UI",
	0x8814:     "RGBA32F",
	0x8d82:     "RGBA32I",
	0x8d70:     "RGBA32UI",
	0x8058:     "RGBA8",
	0x8d8e:     "RGBA8I",
	0x8f97:     "RGBA8_SNORM",
	0x8d7c:     "RGBA8UI",
	0x8d99:     "RGBA_INTEGER",
	0x8d98:     "RGB_INTEGER",
	0x8228:     "RG_INTEGER",
	0x8dc1:     "SAMPLER_2D_ARRAY",
	0x8dc4:     "SAMPLER_2D_ARRAY_SHADOW",
	0x8b62:     "SAMPLER_2D_SHADOW",
	0x8b5f:     "SAMPLER_3D",
	0x8919:     "SAMPLER_BINDING",
	0x8dc5:     "SAMPLER_CUBE_SHADOW",
	0x8c8d:     "SEPARATE_ATTRIBS",
	0x9119:     "SIGNALED",
	0x8f9c:     "SIGNED_NORMALIZED",
	0x8c40:     "SRGB",
	0x8c41:     "SRGB8",
	0x8c43:     "SRGB8_ALPHA8",
	0x88e6:     "STATIC_COPY",
	0x88e5:     "STATIC_READ",
	0x1802:     "STENCIL",
	0x88e2:     "STREAM_COPY",
	0x88e1:     "STREAM_READ",
	0x9113:     "SYNC_CONDITION",
	0x9116:     "SYNC_FENCE",
	0x9115:     "SYNC_FLAGS",
	0x9117:     "SYNC_GPU_COMMANDS_COMPLETE",
	0x9114:     "SYNC_STATUS",
	0x8c1a:     "TEXTURE_2D_ARRAY",
	0x806f:     "TEXTURE_3D",
	0x813c:     "TEXTURE_BASE_LEVEL",
	0x8c1d:     "TEXTURE_BINDING_2D_ARRAY",
	0x806a:     "TEXTURE_BINDING_3D",
	0x884d:     "TEXTURE_COMPARE_FUNC",
	0x884c:     "TEXTURE_COMPARE_MODE",
	0x912f:     "TEXTURE_IMMUTABLE_FORMAT",
	0x82df:     "TEXTURE_IMMUTABLE_LEVELS",
	0x813d:     "TEXTURE_MAX_LEVEL",
	0x813b:     "TEXTURE_MAX_LOD",
	0x813a:     "TEXTURE_MIN_LOD",
	0x8e45:     "TEXTURE_SWIZZLE_A",
	0x8e44:     "TEXTURE_SWIZZLE_B",
	0x8e43:     "TEXTURE_SWIZZLE_G",
	0x8e42:     "TEXTURE_SWIZZLE_R",
	0x8072:     "TEXTURE_WRAP_R",
	0x911b:     "TIMEOUT_EXPIRED",
	0x8e22:     "TRANSFORM_FEEDBACK",
	0x8e24:     "TRANSFORM_FEEDBACK_ACTIVE",
	0x8e25:     "TRANSFORM_FEEDBACK_BINDING",
	0x8c8e:     "TRANSFORM_FEEDBACK_BUFFER",
	0x8c8f:     "TRANSFORM_FEEDBACK_BUFFER_BINDING",
	0x8c7f:     "TRANSFORM_FEEDBACK_BUFFER_MODE",
	0x8c85:     "TRANSFORM_FEEDBACK_BUFFER_SIZE",
	0x8c84:     "TRANSFORM_FEEDBACK_BUFFER_START",
	0x8e23:     "TRANSFORM_FEEDBACK_PAUSED",
	0x8c88:     "TRANSFORM_FEEDBACK_PRIMITIVES_WRITTEN",
	0x8c76:     "TRANSFORM_FEEDBACK_VARYING_MAX_LENGTH",
	0x8c83:     "TRANSFORM_FEEDBACK_VARYINGS",
	0x8a3c:     "UNIFORM_ARRAY_STRIDE",
	0x8a43:     "UNIFORM_BLOCK_ACTIVE_UNIFORM_INDICES",
	0x8a42:     "UNIFORM_BLOCK_ACTIVE_UNIFORMS",
	0x8a3f:     "UNIFORM_BLOCK_BINDING",
	0x8a40:     "UNIFORM_BLOCK_DATA_SIZE",
	0x8a3a:     "UNIFORM_BLOCK_INDEX",
	0x8a41:     "UNIFORM_BLOCK_NAME_LENGTH",
	0x8a46:     "UNIFORM_BLOCK_REFERENCED_BY_FRAGMENT_SHADER",
	0x8a44:     "UNIFORM_BLOCK_REFERENCED_BY_VERTEX_SHADER",
	0x8a11:     "UNIFORM_BUFFER",
	0x8a28:     "UNIFORM_BUFFER_BINDING",
	0x8a34:     "UNIFORM_BUFFER_OFFSET_ALIGNMENT",
	0x8a2a:     "UNIFORM_BUFFER_SIZE",
	0x8a29:     "UNIFORM_BUFFER_START",
	0x8a3e:     "UNIFORM_IS_ROW_MAJOR",
	0x8a3d:     "UNIFORM_MATRIX_STRIDE",
	0x8a39:     "UNIFORM_NAME_LENGTH",
	0x8a3b:     "UNIFORM_OFFSET",
	0x8a38:     "UNIFORM_SIZE",
	0x8a37:     "UNIFORM_TYPE",
	0x806e:     "UNPACK_IMAGE_HEIGHT",
	0xcf2:      "UNPACK_ROW_LENGTH",
	0x806d:     "UNPACK_SKIP_IMAGES",
	0xcf4:      "UNPACK_SKIP_PIXELS",
	0xcf3:      "UNPACK_SKIP_ROWS",
	0x9118:     "UNSIGNALED",
	0x8c3b:     "UNSIGNED_INT_10F_11F_11F_REV",
	0x8368:     "UNSIGNED_INT_2_10_10_10_REV",
	0x84fa:     "UNSIGNED_INT_24_8",
	0x8c3e:     "UNSIGNED_INT_5_9_9_9_REV",
	0x8dd2:     "UNSIGNED_INT_SAMPLER_2D",
	0x8dd7:     "UNSIGNED_INT_SAMPLER_2D_ARRAY",
	0x8dd3:     "UNSIGNED_INT_SAMPLER_3D",
	0x8dd4:     "UNSIGNED_INT_SAMPLER_CUBE",
	0x8dc6:     "UNSIGNED_INT_VEC2",
	0x8dc7:     "UNSIGNED_INT_VEC3",
	0x8dc8:     "UNSIGNED_INT_VEC4",
	0x8c17:     "UNSIGNED_NORMALIZED",
	0x85b5:     "VERTEX_ARRAY_BINDING",
	0x88fe:     "VERTEX_ATTRIB_ARRAY_DIVISOR",
	0x88fd:     "VERTEX_ATTRIB_ARRAY_INTEGER",
	0x911d:     "WAIT_FAILED",
}

func (v Enum) String() string {
	if s, ok := enumString[v]; ok {
		return s
	}
	return fmt.Sprintf("gl.Enum(0x%x)", uint32(v))
}

func (ctx *context) ActiveTexture(texture Enum) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.ActiveTexture(%v) %v", texture, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnActiveTexture,
			a0: texture.c(),
		},
		blocking: true})
}

func (ctx *context) AttachShader(p Program, s Shader) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.AttachShader(%v, %v) %v", p, s, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnAttachShader,
			a0: p.c(),
			a1: s.c(),
		},
		blocking: true})
}

func (ctx *context) BindAttribLocation(p Program, a Attrib, name string) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.BindAttribLocation(%v, %v, %v) %v", p, a, name, errstr)
	}()
	s, free := ctx.cString(name)
	defer free()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnBindAttribLocation,
			a0: p.c(),
			a1: a.c(),
			a2: s,
		},
		blocking: true,
	})
}

func (ctx *context) BindBuffer(target Enum, b Buffer) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.BindBuffer(%v, %v) %v", target, b, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnBindBuffer,
			a0: target.c(),
			a1: b.c(),
		},
		blocking: true})
}

func (ctx *context) BindFramebuffer(target Enum, fb Framebuffer) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.BindFramebuffer(%v, %v) %v", target, fb, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnBindFramebuffer,
			a0: target.c(),
			a1: fb.c(),
		},
		blocking: true})
}

func (ctx *context) BindRenderbuffer(target Enum, rb Renderbuffer) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.BindRenderbuffer(%v, %v) %v", target, rb, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnBindRenderbuffer,
			a0: target.c(),
			a1: rb.c(),
		},
		blocking: true})
}

func (ctx *context) BindTexture(target Enum, t Texture) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.BindTexture(%v, %v) %v", target, t, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnBindTexture,
			a0: target.c(),
			a1: t.c(),
		},
		blocking: true})
}

func (ctx *context) BindVertexArray(va VertexArray) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.BindVertexArray(%v) %v", va, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnBindVertexArray,
			a0: va.c(),
		},
		blocking: true})
}

func (ctx *context) BlendColor(red, green, blue, alpha float32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.BlendColor(%v, %v, %v, %v) %v", red, green, blue, alpha, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnBlendColor,
			a0: uintptr(math.Float32bits(red)),
			a1: uintptr(math.Float32bits(green)),
			a2: uintptr(math.Float32bits(blue)),
			a3: uintptr(math.Float32bits(alpha)),
		},
		blocking: true})
}

func (ctx *context) BlendEquation(mode Enum) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.BlendEquation(%v) %v", mode, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnBlendEquation,
			a0: mode.c(),
		},
		blocking: true})
}

func (ctx *context) BlendEquationSeparate(modeRGB, modeAlpha Enum) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.BlendEquationSeparate(%v, %v) %v", modeRGB, modeAlpha, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnBlendEquationSeparate,
			a0: modeRGB.c(),
			a1: modeAlpha.c(),
		},
		blocking: true})
}

func (ctx *context) BlendFunc(sfactor, dfactor Enum) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.BlendFunc(%v, %v) %v", sfactor, dfactor, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnBlendFunc,
			a0: sfactor.c(),
			a1: dfactor.c(),
		},
		blocking: true})
}

func (ctx *context) BlendFuncSeparate(sfactorRGB, dfactorRGB, sfactorAlpha, dfactorAlpha Enum) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.BlendFuncSeparate(%v, %v, %v, %v) %v", sfactorRGB, dfactorRGB, sfactorAlpha, dfactorAlpha, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnBlendFuncSeparate,
			a0: sfactorRGB.c(),
			a1: dfactorRGB.c(),
			a2: sfactorAlpha.c(),
			a3: dfactorAlpha.c(),
		},
		blocking: true})
}

func (ctx *context) BufferData(target Enum, src []byte, usage Enum) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.BufferData(%v, len(%d), %v) %v", target, len(src), usage, errstr)
	}()
	parg := unsafe.Pointer(nil)
	if len(src) > 0 {
		parg = unsafe.Pointer(&src[0])
	}
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnBufferData,
			a0: target.c(),
			a1: uintptr(len(src)),
			a2: usage.c(),
		},
		parg:     parg,
		blocking: true,
	})
}

func (ctx *context) BufferInit(target Enum, size int, usage Enum) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.BufferInit(%v, %v, %v) %v", target, size, usage, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnBufferData,
			a0: target.c(),
			a1: uintptr(size),
			a2: usage.c(),
		},
		parg:     unsafe.Pointer(nil),
		blocking: true})
}

func (ctx *context) BufferSubData(target Enum, offset int, data []byte) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.BufferSubData(%v, %v, len(%d)) %v", target, offset, len(data), errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnBufferSubData,
			a0: target.c(),
			a1: uintptr(offset),
			a2: uintptr(len(data)),
		},
		parg:     unsafe.Pointer(&data[0]),
		blocking: true,
	})
}

func (ctx *context) CheckFramebufferStatus(target Enum) (r0 Enum) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.CheckFramebufferStatus(%v) %v%v", target, r0, errstr)
	}()
	return Enum(ctx.enqueue(call{
		args: fnargs{
			fn: glfnCheckFramebufferStatus,
			a0: target.c(),
		},
		blocking: true,
	}))
}

func (ctx *context) Clear(mask Enum) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.Clear(%v) %v", mask, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnClear,
			a0: uintptr(mask),
		},
		blocking: true})
}

func (ctx *context) ClearColor(red, green, blue, alpha float32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.ClearColor(%v, %v, %v, %v) %v", red, green, blue, alpha, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnClearColor,
			a0: uintptr(math.Float32bits(red)),
			a1: uintptr(math.Float32bits(green)),
			a2: uintptr(math.Float32bits(blue)),
			a3: uintptr(math.Float32bits(alpha)),
		},
		blocking: true})
}

func (ctx *context) ClearDepthf(d float32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.ClearDepthf(%v) %v", d, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnClearDepthf,
			a0: uintptr(math.Float32bits(d)),
		},
		blocking: true})
}

func (ctx *context) ClearStencil(s int) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.ClearStencil(%v) %v", s, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnClearStencil,
			a0: uintptr(s),
		},
		blocking: true})
}

func (ctx *context) ColorMask(red, green, blue, alpha bool) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.ColorMask(%v, %v, %v, %v) %v", red, green, blue, alpha, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnColorMask,
			a0: glBoolean(red),
			a1: glBoolean(green),
			a2: glBoolean(blue),
			a3: glBoolean(alpha),
		},
		blocking: true})
}

func (ctx *context) CompileShader(s Shader) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.CompileShader(%v) %v", s, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnCompileShader,
			a0: s.c(),
		},
		blocking: true})
}

func (ctx *context) CompressedTexImage2D(target Enum, level int, internalformat Enum, width, height, border int, data []byte) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.CompressedTexImage2D(%v, %v, %v, %v, %v, %v, len(%d)) %v", target, level, internalformat, width, height, border, len(data), errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnCompressedTexImage2D,
			a0: target.c(),
			a1: uintptr(level),
			a2: internalformat.c(),
			a3: uintptr(width),
			a4: uintptr(height),
			a5: uintptr(border),
			a6: uintptr(len(data)),
		},
		parg:     unsafe.Pointer(&data[0]),
		blocking: true,
	})
}

func (ctx *context) CompressedTexSubImage2D(target Enum, level, xoffset, yoffset, width, height int, format Enum, data []byte) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.CompressedTexSubImage2D(%v, %v, %v, %v, %v, %v, %v, len(%d)) %v", target, level, xoffset, yoffset, width, height, format, len(data), errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnCompressedTexSubImage2D,
			a0: target.c(),
			a1: uintptr(level),
			a2: uintptr(xoffset),
			a3: uintptr(yoffset),
			a4: uintptr(width),
			a5: uintptr(height),
			a6: format.c(),
			a7: uintptr(len(data)),
		},
		parg:     unsafe.Pointer(&data[0]),
		blocking: true,
	})
}

func (ctx *context) CopyTexImage2D(target Enum, level int, internalformat Enum, x, y, width, height, border int) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.CopyTexImage2D(%v, %v, %v, %v, %v, %v, %v, %v) %v", target, level, internalformat, x, y, width, height, border, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnCopyTexImage2D,
			a0: target.c(),
			a1: uintptr(level),
			a2: internalformat.c(),
			a3: uintptr(x),
			a4: uintptr(y),
			a5: uintptr(width),
			a6: uintptr(height),
			a7: uintptr(border),
		},
		blocking: true})
}

func (ctx *context) CopyTexSubImage2D(target Enum, level, xoffset, yoffset, x, y, width, height int) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.CopyTexSubImage2D(%v, %v, %v, %v, %v, %v, %v, %v) %v", target, level, xoffset, yoffset, x, y, width, height, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnCopyTexSubImage2D,
			a0: target.c(),
			a1: uintptr(level),
			a2: uintptr(xoffset),
			a3: uintptr(yoffset),
			a4: uintptr(x),
			a5: uintptr(y),
			a6: uintptr(width),
			a7: uintptr(height),
		},
		blocking: true})
}

func (ctx *context) CreateBuffer() (r0 Buffer) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.CreateBuffer() %v%v", r0, errstr)
	}()
	return Buffer{Value: uint32(ctx.enqueue(call{
		args: fnargs{
			fn: glfnGenBuffer,
		},
		blocking: true,
	}))}
}

func (ctx *context) CreateFramebuffer() (r0 Framebuffer) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.CreateFramebuffer() %v%v", r0, errstr)
	}()
	return Framebuffer{Value: uint32(ctx.enqueue(call{
		args: fnargs{
			fn: glfnGenFramebuffer,
		},
		blocking: true,
	}))}
}

func (ctx *context) CreateProgram() (r0 Program) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.CreateProgram() %v%v", r0, errstr)
	}()
	return Program{
		Init: true,
		Value: uint32(ctx.enqueue(call{
			args: fnargs{
				fn: glfnCreateProgram,
			},
			blocking: true,
		},
		))}
}

func (ctx *context) CreateRenderbuffer() (r0 Renderbuffer) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.CreateRenderbuffer() %v%v", r0, errstr)
	}()
	return Renderbuffer{Value: uint32(ctx.enqueue(call{
		args: fnargs{
			fn: glfnGenRenderbuffer,
		},
		blocking: true,
	}))}
}

func (ctx *context) CreateShader(ty Enum) (r0 Shader) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.CreateShader(%v) %v%v", ty, r0, errstr)
	}()
	return Shader{Value: uint32(ctx.enqueue(call{
		args: fnargs{
			fn: glfnCreateShader,
			a0: uintptr(ty),
		},
		blocking: true,
	}))}
}

func (ctx *context) CreateTexture() (r0 Texture) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.CreateTexture() %v%v", r0, errstr)
	}()
	return Texture{Value: uint32(ctx.enqueue(call{
		args: fnargs{
			fn: glfnGenTexture,
		},
		blocking: true,
	}))}
}

func (ctx *context) CreateVertexArray() (r0 VertexArray) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.CreateVertexArray() %v%v", r0, errstr)
	}()
	return VertexArray{Value: uint32(ctx.enqueue(call{
		args: fnargs{
			fn: glfnGenVertexArray,
		},
		blocking: true,
	}))}
}

func (ctx *context) CullFace(mode Enum) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.CullFace(%v) %v", mode, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnCullFace,
			a0: mode.c(),
		},
		blocking: true})
}

func (ctx *context) DeleteBuffer(v Buffer) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.DeleteBuffer(%v) %v", v, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnDeleteBuffer,
			a0: v.c(),
		},
		blocking: true})
}

func (ctx *context) DeleteFramebuffer(v Framebuffer) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.DeleteFramebuffer(%v) %v", v, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnDeleteFramebuffer,
			a0: v.c(),
		},
		blocking: true})
}

func (ctx *context) DeleteProgram(p Program) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.DeleteProgram(%v) %v", p, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnDeleteProgram,
			a0: p.c(),
		},
		blocking: true})
}

func (ctx *context) DeleteRenderbuffer(v Renderbuffer) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.DeleteRenderbuffer(%v) %v", v, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnDeleteRenderbuffer,
			a0: v.c(),
		},
		blocking: true})
}

func (ctx *context) DeleteShader(s Shader) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.DeleteShader(%v) %v", s, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnDeleteShader,
			a0: s.c(),
		},
		blocking: true})
}

func (ctx *context) DeleteTexture(v Texture) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.DeleteTexture(%v) %v", v, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnDeleteTexture,
			a0: v.c(),
		},
		blocking: true})
}

func (ctx *context) DeleteVertexArray(v VertexArray) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.DeleteVertexArray(%v) %v", v, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnDeleteVertexArray,
			a0: v.c(),
		},
		blocking: true})
}

func (ctx *context) DepthFunc(fn Enum) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.DepthFunc(%v) %v", fn, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnDepthFunc,
			a0: fn.c(),
		},
		blocking: true})
}

func (ctx *context) DepthMask(flag bool) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.DepthMask(%v) %v", flag, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnDepthMask,
			a0: glBoolean(flag),
		},
		blocking: true})
}

func (ctx *context) DepthRangef(n, f float32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.DepthRangef(%v, %v) %v", n, f, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnDepthRangef,
			a0: uintptr(math.Float32bits(n)),
			a1: uintptr(math.Float32bits(f)),
		},
		blocking: true})
}

func (ctx *context) DetachShader(p Program, s Shader) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.DetachShader(%v, %v) %v", p, s, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnDetachShader,
			a0: p.c(),
			a1: s.c(),
		},
		blocking: true})
}

func (ctx *context) Disable(cap Enum) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.Disable(%v) %v", cap, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnDisable,
			a0: cap.c(),
		},
		blocking: true})
}

func (ctx *context) DisableVertexAttribArray(a Attrib) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.DisableVertexAttribArray(%v) %v", a, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnDisableVertexAttribArray,
			a0: a.c(),
		},
		blocking: true})
}

func (ctx *context) DrawArrays(mode Enum, first, count int) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.DrawArrays(%v, %v, %v) %v", mode, first, count, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnDrawArrays,
			a0: mode.c(),
			a1: uintptr(first),
			a2: uintptr(count),
		},
		blocking: true})
}

func (ctx *context) DrawElements(mode Enum, count int, ty Enum, offset int) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.DrawElements(%v, %v, %v, %v) %v", mode, count, ty, offset, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnDrawElements,
			a0: mode.c(),
			a1: uintptr(count),
			a2: ty.c(),
			a3: uintptr(offset),
		},
		blocking: true})
}

func (ctx *context) Enable(cap Enum) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.Enable(%v) %v", cap, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnEnable,
			a0: cap.c(),
		},
		blocking: true})
}

func (ctx *context) EnableVertexAttribArray(a Attrib) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.EnableVertexAttribArray(%v) %v", a, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnEnableVertexAttribArray,
			a0: a.c(),
		},
		blocking: true})
}

func (ctx *context) Finish() {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.Finish() %v", errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnFinish,
		},
		blocking: true,
	})
}

func (ctx *context) Flush() {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.Flush() %v", errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnFlush,
		},
		blocking: true,
	})
}

func (ctx *context) FramebufferRenderbuffer(target, attachment, rbTarget Enum, rb Renderbuffer) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.FramebufferRenderbuffer(%v, %v, %v, %v) %v", target, attachment, rbTarget, rb, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnFramebufferRenderbuffer,
			a0: target.c(),
			a1: attachment.c(),
			a2: rbTarget.c(),
			a3: rb.c(),
		},
		blocking: true})
}

func (ctx *context) FramebufferTexture2D(target, attachment, texTarget Enum, t Texture, level int) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.FramebufferTexture2D(%v, %v, %v, %v, %v) %v", target, attachment, texTarget, t, level, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnFramebufferTexture2D,
			a0: target.c(),
			a1: attachment.c(),
			a2: texTarget.c(),
			a3: t.c(),
			a4: uintptr(level),
		},
		blocking: true})
}

func (ctx *context) FrontFace(mode Enum) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.FrontFace(%v) %v", mode, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnFrontFace,
			a0: mode.c(),
		},
		blocking: true})
}

func (ctx *context) GenerateMipmap(target Enum) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.GenerateMipmap(%v) %v", target, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnGenerateMipmap,
			a0: target.c(),
		},
		blocking: true})
}

func (ctx *context) GetActiveAttrib(p Program, index uint32) (name string, size int, ty Enum) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.GetActiveAttrib(%v, %v) (%v, %v, %v) %v", p, index, name, size, ty, errstr)
	}()
	bufSize := ctx.GetProgrami(p, ACTIVE_ATTRIBUTE_MAX_LENGTH)
	buf := make([]byte, bufSize)
	var cType int
	cSize := ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetActiveAttrib,
			a0: p.c(),
			a1: uintptr(index),
			a2: uintptr(bufSize),
			a3: uintptr(unsafe.Pointer(&cType)),
		},
		parg:     unsafe.Pointer(&buf[0]),
		blocking: true,
	})
	return goString(buf), int(cSize), Enum(cType)
}

func (ctx *context) GetActiveUniform(p Program, index uint32) (name string, size int, ty Enum) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.GetActiveUniform(%v, %v) (%v, %v, %v) %v", p, index, name, size, ty, errstr)
	}()
	bufSize := ctx.GetProgrami(p, ACTIVE_UNIFORM_MAX_LENGTH)
	buf := make([]byte, bufSize+8)
	var cType int
	cSize := ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetActiveUniform,
			a0: p.c(),
			a1: uintptr(index),
			a2: uintptr(bufSize),
			a3: uintptr(unsafe.Pointer(&cType)),
		},
		parg:     unsafe.Pointer(&buf[0]),
		blocking: true,
	})
	return goString(buf), int(cSize), Enum(cType)
}

func (ctx *context) GetAttachedShaders(p Program) (r0 []Shader) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.GetAttachedShaders(%v) %v%v", p, r0, errstr)
	}()
	shadersLen := ctx.GetProgrami(p, ATTACHED_SHADERS)
	if shadersLen == 0 {
		return nil
	}
	buf := make([]uint32, shadersLen)
	n := int(ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetAttachedShaders,
			a0: p.c(),
			a1: uintptr(shadersLen),
		},
		parg:     unsafe.Pointer(&buf[0]),
		blocking: true,
	}))
	buf = buf[:int(n)]
	shaders := make([]Shader, len(buf))
	for i, s := range buf {
		shaders[i] = Shader{Value: uint32(s)}
	}
	return shaders
}

func (ctx *context) GetAttribLocation(p Program, name string) (r0 Attrib) {
	defer func() {
		errstr := ctx.errDrain()
		r0.name = name
		log.Printf("gl.GetAttribLocation(%v, %v) %v%v", p, name, r0, errstr)
	}()
	s, free := ctx.cString(name)
	defer free()
	return Attrib{Value: uint(ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetAttribLocation,
			a0: p.c(),
			a1: s,
		},
		blocking: true,
	}))}
}

func (ctx *context) GetBooleanv(dst []bool, pname Enum) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.GetBooleanv(%v, %v) %v", dst, pname, errstr)
	}()
	buf := make([]int32, len(dst))
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnGetBooleanv,
			a0: pname.c(),
		},
		parg:     unsafe.Pointer(&buf[0]),
		blocking: true,
	})
	for i, v := range buf {
		dst[i] = v != 0
	}
}

func (ctx *context) GetFloatv(dst []float32, pname Enum) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.GetFloatv(len(%d), %v) %v", len(dst), pname, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnGetFloatv,
			a0: pname.c(),
		},
		parg:     unsafe.Pointer(&dst[0]),
		blocking: true,
	})
}

func (ctx *context) GetIntegerv(dst []int32, pname Enum) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.GetIntegerv(%v, %v) %v", dst, pname, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnGetIntegerv,
			a0: pname.c(),
		},
		parg:     unsafe.Pointer(&dst[0]),
		blocking: true,
	})
}

func (ctx *context) GetInteger(pname Enum) (r0 int) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.GetInteger(%v) %v%v", pname, r0, errstr)
	}()
	var v [1]int32
	ctx.GetIntegerv(v[:], pname)
	return int(v[0])
}

func (ctx *context) GetBufferParameteri(target, value Enum) (r0 int) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.GetBufferParameteri(%v, %v) %v%v", target, value, r0, errstr)
	}()
	return int(ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetBufferParameteri,
			a0: target.c(),
			a1: value.c(),
		},
		blocking: true,
	}))
}

func (ctx *context) GetError() (r0 Enum) {
	return Enum(ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetError,
		},
		blocking: true,
	}))
}

func (ctx *context) GetFramebufferAttachmentParameteri(target, attachment, pname Enum) (r0 int) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.GetFramebufferAttachmentParameteri(%v, %v, %v) %v%v", target, attachment, pname, r0, errstr)
	}()
	return int(ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetFramebufferAttachmentParameteriv,
			a0: target.c(),
			a1: attachment.c(),
			a2: pname.c(),
		},
		blocking: true,
	}))
}

func (ctx *context) GetProgrami(p Program, pname Enum) (r0 int) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.GetProgrami(%v, %v) %v%v", p, pname, r0, errstr)
	}()
	return int(ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetProgramiv,
			a0: p.c(),
			a1: pname.c(),
		},
		blocking: true,
	}))
}

func (ctx *context) GetProgramInfoLog(p Program) (r0 string) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.GetProgramInfoLog(%v) %v%v", p, r0, errstr)
	}()
	infoLen := ctx.GetProgrami(p, INFO_LOG_LENGTH)
	if infoLen == 0 {
		return ""
	}
	buf := make([]byte, infoLen)
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnGetProgramInfoLog,
			a0: p.c(),
			a1: uintptr(infoLen),
		},
		parg:     unsafe.Pointer(&buf[0]),
		blocking: true,
	})
	return goString(buf)
}

func (ctx *context) GetRenderbufferParameteri(target, pname Enum) (r0 int) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.GetRenderbufferParameteri(%v, %v) %v%v", target, pname, r0, errstr)
	}()
	return int(ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetRenderbufferParameteriv,
			a0: target.c(),
			a1: pname.c(),
		},
		blocking: true,
	}))
}

func (ctx *context) GetShaderi(s Shader, pname Enum) (r0 int) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.GetShaderi(%v, %v) %v%v", s, pname, r0, errstr)
	}()
	return int(ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetShaderiv,
			a0: s.c(),
			a1: pname.c(),
		},
		blocking: true,
	}))
}

func (ctx *context) GetShaderInfoLog(s Shader) (r0 string) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.GetShaderInfoLog(%v) %v%v", s, r0, errstr)
	}()
	infoLen := ctx.GetShaderi(s, INFO_LOG_LENGTH)
	if infoLen == 0 {
		return ""
	}
	buf := make([]byte, infoLen)
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnGetShaderInfoLog,
			a0: s.c(),
			a1: uintptr(infoLen),
		},
		parg:     unsafe.Pointer(&buf[0]),
		blocking: true,
	})
	return goString(buf)
}

func (ctx *context) GetShaderPrecisionFormat(shadertype, precisiontype Enum) (rangeLow, rangeHigh, precision int) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.GetShaderPrecisionFormat(%v, %v) (%v, %v, %v) %v", shadertype, precisiontype, rangeLow, rangeHigh, precision, errstr)
	}()
	var rangeAndPrec [3]int32
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnGetShaderPrecisionFormat,
			a0: shadertype.c(),
			a1: precisiontype.c(),
		},
		parg:     unsafe.Pointer(&rangeAndPrec[0]),
		blocking: true,
	})
	return int(rangeAndPrec[0]), int(rangeAndPrec[1]), int(rangeAndPrec[2])
}

func (ctx *context) GetShaderSource(s Shader) (r0 string) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.GetShaderSource(%v) %v%v", s, r0, errstr)
	}()
	sourceLen := ctx.GetShaderi(s, SHADER_SOURCE_LENGTH)
	if sourceLen == 0 {
		return ""
	}
	buf := make([]byte, sourceLen)
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnGetShaderSource,
			a0: s.c(),
			a1: uintptr(sourceLen),
		},
		parg:     unsafe.Pointer(&buf[0]),
		blocking: true,
	})
	return goString(buf)
}

func (ctx *context) GetTexParameterfv(dst []float32, target, pname Enum) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.GetTexParameterfv(len(%d), %v, %v) %v", len(dst), target, pname, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnGetTexParameterfv,
			a0: target.c(),
			a1: pname.c(),
		},
		parg:     unsafe.Pointer(&dst[0]),
		blocking: true,
	})
}

func (ctx *context) GetTexParameteriv(dst []int32, target, pname Enum) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.GetTexParameteriv(%v, %v, %v) %v", dst, target, pname, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnGetTexParameteriv,
			a0: target.c(),
			a1: pname.c(),
		},
		blocking: true,
	})
}

func (ctx *context) GetUniformfv(dst []float32, src Uniform, p Program) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.GetUniformfv(len(%d), %v, %v) %v", len(dst), src, p, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnGetUniformfv,
			a0: p.c(),
			a1: src.c(),
		},
		parg:     unsafe.Pointer(&dst[0]),
		blocking: true,
	})
}

func (ctx *context) GetUniformiv(dst []int32, src Uniform, p Program) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.GetUniformiv(%v, %v, %v) %v", dst, src, p, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnGetUniformiv,
			a0: p.c(),
			a1: src.c(),
		},
		parg:     unsafe.Pointer(&dst[0]),
		blocking: true,
	})
}

func (ctx *context) GetUniformLocation(p Program, name string) (r0 Uniform) {
	defer func() {
		errstr := ctx.errDrain()
		r0.name = name
		log.Printf("gl.GetUniformLocation(%v, %v) %v%v", p, name, r0, errstr)
	}()
	s, free := ctx.cString(name)
	defer free()
	return Uniform{Value: int32(ctx.enqueue(call{
		args: fnargs{
			fn: glfnGetUniformLocation,
			a0: p.c(),
			a1: s,
		},
		blocking: true,
	}))}
}

func (ctx *context) GetVertexAttribf(src Attrib, pname Enum) (r0 float32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.GetVertexAttribf(%v, %v) %v%v", src, pname, r0, errstr)
	}()
	var params [1]float32
	ctx.GetVertexAttribfv(params[:], src, pname)
	return params[0]
}

func (ctx *context) GetVertexAttribfv(dst []float32, src Attrib, pname Enum) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.GetVertexAttribfv(len(%d), %v, %v) %v", len(dst), src, pname, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnGetVertexAttribfv,
			a0: src.c(),
			a1: pname.c(),
		},
		parg:     unsafe.Pointer(&dst[0]),
		blocking: true,
	})
}

func (ctx *context) GetVertexAttribi(src Attrib, pname Enum) (r0 int32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.GetVertexAttribi(%v, %v) %v%v", src, pname, r0, errstr)
	}()
	var params [1]int32
	ctx.GetVertexAttribiv(params[:], src, pname)
	return params[0]
}

func (ctx *context) GetVertexAttribiv(dst []int32, src Attrib, pname Enum) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.GetVertexAttribiv(%v, %v, %v) %v", dst, src, pname, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnGetVertexAttribiv,
			a0: src.c(),
			a1: pname.c(),
		},
		parg:     unsafe.Pointer(&dst[0]),
		blocking: true,
	})
}

func (ctx *context) Hint(target, mode Enum) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.Hint(%v, %v) %v", target, mode, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnHint,
			a0: target.c(),
			a1: mode.c(),
		},
		blocking: true})
}

func (ctx *context) IsBuffer(b Buffer) (r0 bool) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.IsBuffer(%v) %v%v", b, r0, errstr)
	}()
	return 0 != ctx.enqueue(call{
		args: fnargs{
			fn: glfnIsBuffer,
			a0: b.c(),
		},
		blocking: true,
	})
}

func (ctx *context) IsEnabled(cap Enum) (r0 bool) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.IsEnabled(%v) %v%v", cap, r0, errstr)
	}()
	return 0 != ctx.enqueue(call{
		args: fnargs{
			fn: glfnIsEnabled,
			a0: cap.c(),
		},
		blocking: true,
	})
}

func (ctx *context) IsFramebuffer(fb Framebuffer) (r0 bool) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.IsFramebuffer(%v) %v%v", fb, r0, errstr)
	}()
	return 0 != ctx.enqueue(call{
		args: fnargs{
			fn: glfnIsFramebuffer,
			a0: fb.c(),
		},
		blocking: true,
	})
}

func (ctx *context) IsProgram(p Program) (r0 bool) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.IsProgram(%v) %v%v", p, r0, errstr)
	}()
	return 0 != ctx.enqueue(call{
		args: fnargs{
			fn: glfnIsProgram,
			a0: p.c(),
		},
		blocking: true,
	})
}

func (ctx *context) IsRenderbuffer(rb Renderbuffer) (r0 bool) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.IsRenderbuffer(%v) %v%v", rb, r0, errstr)
	}()
	return 0 != ctx.enqueue(call{
		args: fnargs{
			fn: glfnIsRenderbuffer,
			a0: rb.c(),
		},
		blocking: true,
	})
}

func (ctx *context) IsShader(s Shader) (r0 bool) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.IsShader(%v) %v%v", s, r0, errstr)
	}()
	return 0 != ctx.enqueue(call{
		args: fnargs{
			fn: glfnIsShader,
			a0: s.c(),
		},
		blocking: true,
	})
}

func (ctx *context) IsTexture(t Texture) (r0 bool) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.IsTexture(%v) %v%v", t, r0, errstr)
	}()
	return 0 != ctx.enqueue(call{
		args: fnargs{
			fn: glfnIsTexture,
			a0: t.c(),
		},
		blocking: true,
	})
}

func (ctx *context) LineWidth(width float32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.LineWidth(%v) %v", width, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnLineWidth,
			a0: uintptr(math.Float32bits(width)),
		},
		blocking: true})
}

func (ctx *context) LinkProgram(p Program) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.LinkProgram(%v) %v", p, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnLinkProgram,
			a0: p.c(),
		},
		blocking: true})
}

func (ctx *context) PixelStorei(pname Enum, param int32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.PixelStorei(%v, %v) %v", pname, param, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnPixelStorei,
			a0: pname.c(),
			a1: uintptr(param),
		},
		blocking: true})
}

func (ctx *context) PolygonOffset(factor, units float32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.PolygonOffset(%v, %v) %v", factor, units, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnPolygonOffset,
			a0: uintptr(math.Float32bits(factor)),
			a1: uintptr(math.Float32bits(units)),
		},
		blocking: true})
}

func (ctx *context) ReadPixels(dst []byte, x, y, width, height int, format, ty Enum) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.ReadPixels(len(%d), %v, %v, %v, %v, %v, %v) %v", len(dst), x, y, width, height, format, ty, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnReadPixels,

			a0: uintptr(x),
			a1: uintptr(y),
			a2: uintptr(width),
			a3: uintptr(height),
			a4: format.c(),
			a5: ty.c(),
		},
		parg:     unsafe.Pointer(&dst[0]),
		blocking: true,
	})
}

func (ctx *context) ReleaseShaderCompiler() {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.ReleaseShaderCompiler() %v", errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnReleaseShaderCompiler,
		},
		blocking: true})
}

func (ctx *context) RenderbufferStorage(target, internalFormat Enum, width, height int) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.RenderbufferStorage(%v, %v, %v, %v) %v", target, internalFormat, width, height, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnRenderbufferStorage,
			a0: target.c(),
			a1: internalFormat.c(),
			a2: uintptr(width),
			a3: uintptr(height),
		},
		blocking: true})
}

func (ctx *context) SampleCoverage(value float32, invert bool) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.SampleCoverage(%v, %v) %v", value, invert, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnSampleCoverage,
			a0: uintptr(math.Float32bits(value)),
			a1: glBoolean(invert),
		},
		blocking: true})
}

func (ctx *context) Scissor(x, y, width, height int32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.Scissor(%v, %v, %v, %v) %v", x, y, width, height, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnScissor,
			a0: uintptr(x),
			a1: uintptr(y),
			a2: uintptr(width),
			a3: uintptr(height),
		},
		blocking: true})
}

func (ctx *context) ShaderSource(s Shader, src string) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.ShaderSource(%v, %v) %v", s, src, errstr)
	}()
	strp, free := ctx.cStringPtr(src)
	defer free()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnShaderSource,
			a0: s.c(),
			a1: 1,
			a2: strp,
		},
		blocking: true,
	})
}

func (ctx *context) StencilFunc(fn Enum, ref int, mask uint32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.StencilFunc(%v, %v, %v) %v", fn, ref, mask, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnStencilFunc,
			a0: fn.c(),
			a1: uintptr(ref),
			a2: uintptr(mask),
		},
		blocking: true})
}

func (ctx *context) StencilFuncSeparate(face, fn Enum, ref int, mask uint32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.StencilFuncSeparate(%v, %v, %v, %v) %v", face, fn, ref, mask, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnStencilFuncSeparate,
			a0: face.c(),
			a1: fn.c(),
			a2: uintptr(ref),
			a3: uintptr(mask),
		},
		blocking: true})
}

func (ctx *context) StencilMask(mask uint32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.StencilMask(%v) %v", mask, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnStencilMask,
			a0: uintptr(mask),
		},
		blocking: true})
}

func (ctx *context) StencilMaskSeparate(face Enum, mask uint32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.StencilMaskSeparate(%v, %v) %v", face, mask, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnStencilMaskSeparate,
			a0: face.c(),
			a1: uintptr(mask),
		},
		blocking: true})
}

func (ctx *context) StencilOp(fail, zfail, zpass Enum) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.StencilOp(%v, %v, %v) %v", fail, zfail, zpass, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnStencilOp,
			a0: fail.c(),
			a1: zfail.c(),
			a2: zpass.c(),
		},
		blocking: true})
}

func (ctx *context) StencilOpSeparate(face, sfail, dpfail, dppass Enum) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.StencilOpSeparate(%v, %v, %v, %v) %v", face, sfail, dpfail, dppass, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnStencilOpSeparate,
			a0: face.c(),
			a1: sfail.c(),
			a2: dpfail.c(),
			a3: dppass.c(),
		},
		blocking: true})
}

func (ctx *context) TexImage2D(target Enum, level int, internalFormat int, width, height int, format Enum, ty Enum, data []byte) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.TexImage2D(%v, %v, %v, %v, %v, %v, %v, len(%d)) %v", target, level, internalFormat, width, height, format, ty, len(data), errstr)
	}()
	parg := unsafe.Pointer(nil)
	if len(data) > 0 {
		parg = unsafe.Pointer(&data[0])
	}
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnTexImage2D,

			a0: target.c(),
			a1: uintptr(level),
			a2: uintptr(internalFormat),
			a3: uintptr(width),
			a4: uintptr(height),
			a5: format.c(),
			a6: ty.c(),
		},
		parg:     parg,
		blocking: true,
	})
}

func (ctx *context) TexSubImage2D(target Enum, level int, x, y, width, height int, format, ty Enum, data []byte) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.TexSubImage2D(%v, %v, %v, %v, %v, %v, %v, %v, len(%d)) %v", target, level, x, y, width, height, format, ty, len(data), errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnTexSubImage2D,

			a0: target.c(),
			a1: uintptr(level),
			a2: uintptr(x),
			a3: uintptr(y),
			a4: uintptr(width),
			a5: uintptr(height),
			a6: format.c(),
			a7: ty.c(),
		},
		parg:     unsafe.Pointer(&data[0]),
		blocking: true,
	})
}

func (ctx *context) TexParameterf(target, pname Enum, param float32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.TexParameterf(%v, %v, %v) %v", target, pname, param, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnTexParameterf,
			a0: target.c(),
			a1: pname.c(),
			a2: uintptr(math.Float32bits(param)),
		},
		blocking: true})
}

func (ctx *context) TexParameterfv(target, pname Enum, params []float32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.TexParameterfv(%v, %v, len(%d)) %v", target, pname, len(params), errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnTexParameterfv,
			a0: target.c(),
			a1: pname.c(),
		},
		parg:     unsafe.Pointer(&params[0]),
		blocking: true,
	})
}

func (ctx *context) TexParameteri(target, pname Enum, param int) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.TexParameteri(%v, %v, %v) %v", target, pname, param, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnTexParameteri,
			a0: target.c(),
			a1: pname.c(),
			a2: uintptr(param),
		},
		blocking: true})
}

func (ctx *context) TexParameteriv(target, pname Enum, params []int32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.TexParameteriv(%v, %v, %v) %v", target, pname, params, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnTexParameteriv,
			a0: target.c(),
			a1: pname.c(),
		},
		parg:     unsafe.Pointer(&params[0]),
		blocking: true,
	})
}

func (ctx *context) Uniform1f(dst Uniform, v float32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.Uniform1f(%v, %v) %v", dst, v, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnUniform1f,
			a0: dst.c(),
			a1: uintptr(math.Float32bits(v)),
		},
		blocking: true})
}

func (ctx *context) Uniform1fv(dst Uniform, src []float32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.Uniform1fv(%v, len(%d)) %v", dst, len(src), errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnUniform1fv,
			a0: dst.c(),
			a1: uintptr(len(src)),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx *context) Uniform1i(dst Uniform, v int) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.Uniform1i(%v, %v) %v", dst, v, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnUniform1i,
			a0: dst.c(),
			a1: uintptr(v),
		},
		blocking: true})
}

func (ctx *context) Uniform1iv(dst Uniform, src []int32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.Uniform1iv(%v, %v) %v", dst, src, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnUniform1iv,
			a0: dst.c(),
			a1: uintptr(len(src)),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx *context) Uniform2f(dst Uniform, v0, v1 float32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.Uniform2f(%v, %v, %v) %v", dst, v0, v1, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnUniform2f,
			a0: dst.c(),
			a1: uintptr(math.Float32bits(v0)),
			a2: uintptr(math.Float32bits(v1)),
		},
		blocking: true})
}

func (ctx *context) Uniform2fv(dst Uniform, src []float32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.Uniform2fv(%v, len(%d)) %v", dst, len(src), errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnUniform2fv,
			a0: dst.c(),
			a1: uintptr(len(src) / 2),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx *context) Uniform2i(dst Uniform, v0, v1 int) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.Uniform2i(%v, %v, %v) %v", dst, v0, v1, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnUniform2i,
			a0: dst.c(),
			a1: uintptr(v0),
			a2: uintptr(v1),
		},
		blocking: true})
}

func (ctx *context) Uniform2iv(dst Uniform, src []int32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.Uniform2iv(%v, %v) %v", dst, src, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnUniform2iv,
			a0: dst.c(),
			a1: uintptr(len(src) / 2),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx *context) Uniform3f(dst Uniform, v0, v1, v2 float32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.Uniform3f(%v, %v, %v, %v) %v", dst, v0, v1, v2, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnUniform3f,
			a0: dst.c(),
			a1: uintptr(math.Float32bits(v0)),
			a2: uintptr(math.Float32bits(v1)),
			a3: uintptr(math.Float32bits(v2)),
		},
		blocking: true})
}

func (ctx *context) Uniform3fv(dst Uniform, src []float32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.Uniform3fv(%v, len(%d)) %v", dst, len(src), errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnUniform3fv,
			a0: dst.c(),
			a1: uintptr(len(src) / 3),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx *context) Uniform3i(dst Uniform, v0, v1, v2 int32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.Uniform3i(%v, %v, %v, %v) %v", dst, v0, v1, v2, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnUniform3i,
			a0: dst.c(),
			a1: uintptr(v0),
			a2: uintptr(v1),
			a3: uintptr(v2),
		},
		blocking: true})
}

func (ctx *context) Uniform3iv(dst Uniform, src []int32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.Uniform3iv(%v, %v) %v", dst, src, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnUniform3iv,
			a0: dst.c(),
			a1: uintptr(len(src) / 3),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx *context) Uniform4f(dst Uniform, v0, v1, v2, v3 float32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.Uniform4f(%v, %v, %v, %v, %v) %v", dst, v0, v1, v2, v3, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnUniform4f,
			a0: dst.c(),
			a1: uintptr(math.Float32bits(v0)),
			a2: uintptr(math.Float32bits(v1)),
			a3: uintptr(math.Float32bits(v2)),
			a4: uintptr(math.Float32bits(v3)),
		},
		blocking: true})
}

func (ctx *context) Uniform4fv(dst Uniform, src []float32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.Uniform4fv(%v, len(%d)) %v", dst, len(src), errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnUniform4fv,
			a0: dst.c(),
			a1: uintptr(len(src) / 4),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx *context) Uniform4i(dst Uniform, v0, v1, v2, v3 int32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.Uniform4i(%v, %v, %v, %v, %v) %v", dst, v0, v1, v2, v3, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnUniform4i,
			a0: dst.c(),
			a1: uintptr(v0),
			a2: uintptr(v1),
			a3: uintptr(v2),
			a4: uintptr(v3),
		},
		blocking: true})
}

func (ctx *context) Uniform4iv(dst Uniform, src []int32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.Uniform4iv(%v, %v) %v", dst, src, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnUniform4iv,
			a0: dst.c(),
			a1: uintptr(len(src) / 4),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx *context) UniformMatrix2fv(dst Uniform, src []float32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.UniformMatrix2fv(%v, len(%d)) %v", dst, len(src), errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnUniformMatrix2fv,

			a0: dst.c(),
			a1: uintptr(len(src) / 4),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx *context) UniformMatrix3fv(dst Uniform, src []float32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.UniformMatrix3fv(%v, len(%d)) %v", dst, len(src), errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnUniformMatrix3fv,
			a0: dst.c(),
			a1: uintptr(len(src) / 9),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx *context) UniformMatrix4fv(dst Uniform, src []float32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.UniformMatrix4fv(%v, len(%d)) %v", dst, len(src), errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnUniformMatrix4fv,
			a0: dst.c(),
			a1: uintptr(len(src) / 16),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx *context) UseProgram(p Program) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.UseProgram(%v) %v", p, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnUseProgram,
			a0: p.c(),
		},
		blocking: true})
}

func (ctx *context) ValidateProgram(p Program) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.ValidateProgram(%v) %v", p, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnValidateProgram,
			a0: p.c(),
		},
		blocking: true})
}

func (ctx *context) VertexAttrib1f(dst Attrib, x float32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.VertexAttrib1f(%v, %v) %v", dst, x, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnVertexAttrib1f,
			a0: dst.c(),
			a1: uintptr(math.Float32bits(x)),
		},
		blocking: true})
}

func (ctx *context) VertexAttrib1fv(dst Attrib, src []float32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.VertexAttrib1fv(%v, len(%d)) %v", dst, len(src), errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnVertexAttrib1fv,
			a0: dst.c(),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx *context) VertexAttrib2f(dst Attrib, x, y float32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.VertexAttrib2f(%v, %v, %v) %v", dst, x, y, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnVertexAttrib2f,
			a0: dst.c(),
			a1: uintptr(math.Float32bits(x)),
			a2: uintptr(math.Float32bits(y)),
		},
		blocking: true})
}

func (ctx *context) VertexAttrib2fv(dst Attrib, src []float32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.VertexAttrib2fv(%v, len(%d)) %v", dst, len(src), errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnVertexAttrib2fv,
			a0: dst.c(),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx *context) VertexAttrib3f(dst Attrib, x, y, z float32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.VertexAttrib3f(%v, %v, %v, %v) %v", dst, x, y, z, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnVertexAttrib3f,
			a0: dst.c(),
			a1: uintptr(math.Float32bits(x)),
			a2: uintptr(math.Float32bits(y)),
			a3: uintptr(math.Float32bits(z)),
		},
		blocking: true})
}

func (ctx *context) VertexAttrib3fv(dst Attrib, src []float32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.VertexAttrib3fv(%v, len(%d)) %v", dst, len(src), errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnVertexAttrib3fv,
			a0: dst.c(),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx *context) VertexAttrib4f(dst Attrib, x, y, z, w float32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.VertexAttrib4f(%v, %v, %v, %v, %v) %v", dst, x, y, z, w, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnVertexAttrib4f,
			a0: dst.c(),
			a1: uintptr(math.Float32bits(x)),
			a2: uintptr(math.Float32bits(y)),
			a3: uintptr(math.Float32bits(z)),
			a4: uintptr(math.Float32bits(w)),
		},
		blocking: true})
}

func (ctx *context) VertexAttrib4fv(dst Attrib, src []float32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.VertexAttrib4fv(%v, len(%d)) %v", dst, len(src), errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnVertexAttrib4fv,
			a0: dst.c(),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx *context) VertexAttribPointer(dst Attrib, size int, ty Enum, normalized bool, stride, offset int) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.VertexAttribPointer(%v, %v, %v, %v, %v, %v) %v", dst, size, ty, normalized, stride, offset, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnVertexAttribPointer,
			a0: dst.c(),
			a1: uintptr(size),
			a2: ty.c(),
			a3: glBoolean(normalized),
			a4: uintptr(stride),
			a5: uintptr(offset),
		},
		blocking: true})
}

func (ctx *context) Viewport(x, y, width, height int) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.Viewport(%v, %v, %v, %v) %v", x, y, width, height, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnViewport,
			a0: uintptr(x),
			a1: uintptr(y),
			a2: uintptr(width),
			a3: uintptr(height),
		},
		blocking: true})
}

func (ctx context3) UniformMatrix2x3fv(dst Uniform, src []float32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.UniformMatrix2x3fv(%v, len(%d)) %v", dst, len(src), errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnUniformMatrix2x3fv,
			a0: dst.c(),
			a1: uintptr(len(src) / 6),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx context3) UniformMatrix3x2fv(dst Uniform, src []float32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.UniformMatrix3x2fv(%v, len(%d)) %v", dst, len(src), errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnUniformMatrix3x2fv,
			a0: dst.c(),
			a1: uintptr(len(src) / 6),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx context3) UniformMatrix2x4fv(dst Uniform, src []float32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.UniformMatrix2x4fv(%v, len(%d)) %v", dst, len(src), errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnUniformMatrix2x4fv,
			a0: dst.c(),
			a1: uintptr(len(src) / 8),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx context3) UniformMatrix4x2fv(dst Uniform, src []float32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.UniformMatrix4x2fv(%v, len(%d)) %v", dst, len(src), errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnUniformMatrix4x2fv,
			a0: dst.c(),
			a1: uintptr(len(src) / 8),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx context3) UniformMatrix3x4fv(dst Uniform, src []float32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.UniformMatrix3x4fv(%v, len(%d)) %v", dst, len(src), errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnUniformMatrix3x4fv,
			a0: dst.c(),
			a1: uintptr(len(src) / 12),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx context3) UniformMatrix4x3fv(dst Uniform, src []float32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.UniformMatrix4x3fv(%v, len(%d)) %v", dst, len(src), errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnUniformMatrix4x3fv,
			a0: dst.c(),
			a1: uintptr(len(src) / 12),
		},
		parg:     unsafe.Pointer(&src[0]),
		blocking: true,
	})
}

func (ctx context3) BlitFramebuffer(srcX0, srcY0, srcX1, srcY1, dstX0, dstY0, dstX1, dstY1 int, mask uint, filter Enum) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.BlitFramebuffer(%v, %v, %v, %v, %v, %v, %v, %v, %v, %v) %v", srcX0, srcY0, srcX1, srcY1, dstX0, dstY0, dstX1, dstY1, mask, filter, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnBlitFramebuffer,
			a0: uintptr(srcX0),
			a1: uintptr(srcY0),
			a2: uintptr(srcX1),
			a3: uintptr(srcY1),
			a4: uintptr(dstX0),
			a5: uintptr(dstY0),
			a6: uintptr(dstX1),
			a7: uintptr(dstY1),
			a8: uintptr(mask),
			a9: filter.c(),
		},
		blocking: true})
}

func (ctx context3) Uniform1ui(dst Uniform, v uint32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.Uniform1ui(%v, %v) %v", dst, v, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnUniform1ui,
			a0: dst.c(),
			a1: uintptr(v),
		},
		blocking: true})
}

func (ctx context3) Uniform2ui(dst Uniform, v0, v1 uint32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.Uniform2ui(%v, %v, %v) %v", dst, v0, v1, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnUniform2ui,
			a0: dst.c(),
			a1: uintptr(v0),
			a2: uintptr(v1),
		},
		blocking: true})
}

func (ctx context3) Uniform3ui(dst Uniform, v0, v1, v2 uint) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.Uniform3ui(%v, %v, %v, %v) %v", dst, v0, v1, v2, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnUniform3ui,
			a0: dst.c(),
			a1: uintptr(v0),
			a2: uintptr(v1),
			a3: uintptr(v2),
		},
		blocking: true})
}

func (ctx context3) Uniform4ui(dst Uniform, v0, v1, v2, v3 uint32) {
	defer func() {
		errstr := ctx.errDrain()
		log.Printf("gl.Uniform4ui(%v, %v, %v, %v, %v) %v", dst, v0, v1, v2, v3, errstr)
	}()
	ctx.enqueueDebug(call{
		args: fnargs{
			fn: glfnUniform4ui,
			a0: dst.c(),
			a1: uintptr(v0),
			a2: uintptr(v1),
			a3: uintptr(v2),
			a4: uintptr(v3),
		},
		blocking: true})
}
