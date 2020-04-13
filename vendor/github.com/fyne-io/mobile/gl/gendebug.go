// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// The gendebug program takes gl.go and generates a version of it
// where each function includes tracing code that writes its arguments
// to the standard log.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

var enumWhitelist = []string{
	"POINTS",
	"LINES",
	"LINE_LOOP",
	"LINE_STRIP",
	"TRIANGLES",
	"TRIANGLE_STRIP",
	"TRIANGLE_FAN",
	"SRC_COLOR",
	"ONE_MINUS_SRC_COLOR",
	"SRC_ALPHA",
	"ONE_MINUS_SRC_ALPHA",
	"DST_ALPHA",
	"ONE_MINUS_DST_ALPHA",
	"DST_COLOR",
	"ONE_MINUS_DST_COLOR",
	"SRC_ALPHA_SATURATE",
	"FUNC_ADD",
	"BLEND_EQUATION",
	"BLEND_EQUATION_RGB",
	"BLEND_EQUATION_ALPHA",
	"FUNC_SUBTRACT",
	"FUNC_REVERSE_SUBTRACT",
	"BLEND_DST_RGB",
	"BLEND_SRC_RGB",
	"BLEND_DST_ALPHA",
	"BLEND_SRC_ALPHA",
	"CONSTANT_COLOR",
	"ONE_MINUS_CONSTANT_COLOR",
	"CONSTANT_ALPHA",
	"ONE_MINUS_CONSTANT_ALPHA",
	"BLEND_COLOR",
	"ARRAY_BUFFER",
	"ELEMENT_ARRAY_BUFFER",
	"ARRAY_BUFFER_BINDING",
	"ELEMENT_ARRAY_BUFFER_BINDING",
	"STREAM_DRAW",
	"STATIC_DRAW",
	"DYNAMIC_DRAW",
	"BUFFER_SIZE",
	"BUFFER_USAGE",
	"CURRENT_VERTEX_ATTRIB",
	"FRONT",
	"BACK",
	"FRONT_AND_BACK",
	"TEXTURE_2D",
	"CULL_FACE",
	"BLEND",
	"DITHER",
	"STENCIL_TEST",
	"DEPTH_TEST",
	"SCISSOR_TEST",
	"POLYGON_OFFSET_FILL",
	"SAMPLE_ALPHA_TO_COVERAGE",
	"SAMPLE_COVERAGE",
	"INVALID_ENUM",
	"INVALID_VALUE",
	"INVALID_OPERATION",
	"OUT_OF_MEMORY",
	"CW",
	"CCW",
	"LINE_WIDTH",
	"ALIASED_POINT_SIZE_RANGE",
	"ALIASED_LINE_WIDTH_RANGE",
	"CULL_FACE_MODE",
	"FRONT_FACE",
	"DEPTH_RANGE",
	"DEPTH_WRITEMASK",
	"DEPTH_CLEAR_VALUE",
	"DEPTH_FUNC",
	"STENCIL_CLEAR_VALUE",
	"STENCIL_FUNC",
	"STENCIL_FAIL",
	"STENCIL_PASS_DEPTH_FAIL",
	"STENCIL_PASS_DEPTH_PASS",
	"STENCIL_REF",
	"STENCIL_VALUE_MASK",
	"STENCIL_WRITEMASK",
	"STENCIL_BACK_FUNC",
	"STENCIL_BACK_FAIL",
	"STENCIL_BACK_PASS_DEPTH_FAIL",
	"STENCIL_BACK_PASS_DEPTH_PASS",
	"STENCIL_BACK_REF",
	"STENCIL_BACK_VALUE_MASK",
	"STENCIL_BACK_WRITEMASK",
	"VIEWPORT",
	"SCISSOR_BOX",
	"COLOR_CLEAR_VALUE",
	"COLOR_WRITEMASK",
	"UNPACK_ALIGNMENT",
	"PACK_ALIGNMENT",
	"MAX_TEXTURE_SIZE",
	"MAX_VIEWPORT_DIMS",
	"SUBPIXEL_BITS",
	"RED_BITS",
	"GREEN_BITS",
	"BLUE_BITS",
	"ALPHA_BITS",
	"DEPTH_BITS",
	"STENCIL_BITS",
	"POLYGON_OFFSET_UNITS",
	"POLYGON_OFFSET_FACTOR",
	"TEXTURE_BINDING_2D",
	"SAMPLE_BUFFERS",
	"SAMPLES",
	"SAMPLE_COVERAGE_VALUE",
	"SAMPLE_COVERAGE_INVERT",
	"NUM_COMPRESSED_TEXTURE_FORMATS",
	"COMPRESSED_TEXTURE_FORMATS",
	"DONT_CARE",
	"FASTEST",
	"NICEST",
	"GENERATE_MIPMAP_HINT",
	"BYTE",
	"UNSIGNED_BYTE",
	"SHORT",
	"UNSIGNED_SHORT",
	"INT",
	"UNSIGNED_INT",
	"FLOAT",
	"FIXED",
	"DEPTH_COMPONENT",
	"ALPHA",
	"RGB",
	"RGBA",
	"LUMINANCE",
	"LUMINANCE_ALPHA",
	"UNSIGNED_SHORT_4_4_4_4",
	"UNSIGNED_SHORT_5_5_5_1",
	"UNSIGNED_SHORT_5_6_5",
	"MAX_VERTEX_ATTRIBS",
	"MAX_VERTEX_UNIFORM_VECTORS",
	"MAX_VARYING_VECTORS",
	"MAX_COMBINED_TEXTURE_IMAGE_UNITS",
	"MAX_VERTEX_TEXTURE_IMAGE_UNITS",
	"MAX_TEXTURE_IMAGE_UNITS",
	"MAX_FRAGMENT_UNIFORM_VECTORS",
	"SHADER_TYPE",
	"DELETE_STATUS",
	"LINK_STATUS",
	"VALIDATE_STATUS",
	"ATTACHED_SHADERS",
	"ACTIVE_UNIFORMS",
	"ACTIVE_UNIFORM_MAX_LENGTH",
	"ACTIVE_ATTRIBUTES",
	"ACTIVE_ATTRIBUTE_MAX_LENGTH",
	"SHADING_LANGUAGE_VERSION",
	"CURRENT_PROGRAM",
	"NEVER",
	"LESS",
	"EQUAL",
	"LEQUAL",
	"GREATER",
	"NOTEQUAL",
	"GEQUAL",
	"ALWAYS",
	"KEEP",
	"REPLACE",
	"INCR",
	"DECR",
	"INVERT",
	"INCR_WRAP",
	"DECR_WRAP",
	"VENDOR",
	"RENDERER",
	"VERSION",
	"EXTENSIONS",
	"NEAREST",
	"LINEAR",
	"NEAREST_MIPMAP_NEAREST",
	"LINEAR_MIPMAP_NEAREST",
	"NEAREST_MIPMAP_LINEAR",
	"LINEAR_MIPMAP_LINEAR",
	"TEXTURE_MAG_FILTER",
	"TEXTURE_MIN_FILTER",
	"TEXTURE_WRAP_S",
	"TEXTURE_WRAP_T",
	"TEXTURE",
	"TEXTURE_CUBE_MAP",
	"TEXTURE_BINDING_CUBE_MAP",
	"TEXTURE_CUBE_MAP_POSITIVE_X",
	"TEXTURE_CUBE_MAP_NEGATIVE_X",
	"TEXTURE_CUBE_MAP_POSITIVE_Y",
	"TEXTURE_CUBE_MAP_NEGATIVE_Y",
	"TEXTURE_CUBE_MAP_POSITIVE_Z",
	"TEXTURE_CUBE_MAP_NEGATIVE_Z",
	"MAX_CUBE_MAP_TEXTURE_SIZE",
	"TEXTURE0",
	"TEXTURE1",
	"TEXTURE2",
	"TEXTURE3",
	"TEXTURE4",
	"TEXTURE5",
	"TEXTURE6",
	"TEXTURE7",
	"TEXTURE8",
	"TEXTURE9",
	"TEXTURE10",
	"TEXTURE11",
	"TEXTURE12",
	"TEXTURE13",
	"TEXTURE14",
	"TEXTURE15",
	"TEXTURE16",
	"TEXTURE17",
	"TEXTURE18",
	"TEXTURE19",
	"TEXTURE20",
	"TEXTURE21",
	"TEXTURE22",
	"TEXTURE23",
	"TEXTURE24",
	"TEXTURE25",
	"TEXTURE26",
	"TEXTURE27",
	"TEXTURE28",
	"TEXTURE29",
	"TEXTURE30",
	"TEXTURE31",
	"ACTIVE_TEXTURE",
	"REPEAT",
	"CLAMP_TO_EDGE",
	"MIRRORED_REPEAT",
	"VERTEX_ATTRIB_ARRAY_ENABLED",
	"VERTEX_ATTRIB_ARRAY_SIZE",
	"VERTEX_ATTRIB_ARRAY_STRIDE",
	"VERTEX_ATTRIB_ARRAY_TYPE",
	"VERTEX_ATTRIB_ARRAY_NORMALIZED",
	"VERTEX_ATTRIB_ARRAY_POINTER",
	"VERTEX_ATTRIB_ARRAY_BUFFER_BINDING",
	"IMPLEMENTATION_COLOR_READ_TYPE",
	"IMPLEMENTATION_COLOR_READ_FORMAT",
	"COMPILE_STATUS",
	"INFO_LOG_LENGTH",
	"SHADER_SOURCE_LENGTH",
	"SHADER_COMPILER",
	"SHADER_BINARY_FORMATS",
	"NUM_SHADER_BINARY_FORMATS",
	"LOW_FLOAT",
	"MEDIUM_FLOAT",
	"HIGH_FLOAT",
	"LOW_INT",
	"MEDIUM_INT",
	"HIGH_INT",
	"FRAMEBUFFER",
	"RENDERBUFFER",
	"RGBA4",
	"RGB5_A1",
	"RGB565",
	"DEPTH_COMPONENT16",
	"STENCIL_INDEX8",
	"RENDERBUFFER_WIDTH",
	"RENDERBUFFER_HEIGHT",
	"RENDERBUFFER_INTERNAL_FORMAT",
	"RENDERBUFFER_RED_SIZE",
	"RENDERBUFFER_GREEN_SIZE",
	"RENDERBUFFER_BLUE_SIZE",
	"RENDERBUFFER_ALPHA_SIZE",
	"RENDERBUFFER_DEPTH_SIZE",
	"RENDERBUFFER_STENCIL_SIZE",
	"FRAMEBUFFER_ATTACHMENT_OBJECT_TYPE",
	"FRAMEBUFFER_ATTACHMENT_OBJECT_NAME",
	"FRAMEBUFFER_ATTACHMENT_TEXTURE_LEVEL",
	"FRAMEBUFFER_ATTACHMENT_TEXTURE_CUBE_MAP_FACE",
	"COLOR_ATTACHMENT0",
	"DEPTH_ATTACHMENT",
	"STENCIL_ATTACHMENT",
	"FRAMEBUFFER_COMPLETE",
	"FRAMEBUFFER_INCOMPLETE_ATTACHMENT",
	"FRAMEBUFFER_INCOMPLETE_MISSING_ATTACHMENT",
	"FRAMEBUFFER_INCOMPLETE_DIMENSIONS",
	"FRAMEBUFFER_UNSUPPORTED",
	"FRAMEBUFFER_BINDING",
	"RENDERBUFFER_BINDING",
	"MAX_RENDERBUFFER_SIZE",
	"INVALID_FRAMEBUFFER_OPERATION",
	"DEPTH_BUFFER_BIT",
	"STENCIL_BUFFER_BIT",
	"COLOR_BUFFER_BIT",
	"FLOAT_VEC2",
	"FLOAT_VEC3",
	"FLOAT_VEC4",
	"INT_VEC2",
	"INT_VEC3",
	"INT_VEC4",
	"BOOL",
	"BOOL_VEC2",
	"BOOL_VEC3",
	"BOOL_VEC4",
	"FLOAT_MAT2",
	"FLOAT_MAT3",
	"FLOAT_MAT4",
	"SAMPLER_2D",
	"SAMPLER_CUBE",
	"FRAGMENT_SHADER",
	"VERTEX_SHADER",
	"FALSE",
	"TRUE",
	"ZERO",
	"ONE",
	"NO_ERROR",
	"NONE",
	"ACTIVE_UNIFORM_BLOCK_MAX_NAME_LENGTH",
	"ACTIVE_UNIFORM_BLOCKS",
	"ALREADY_SIGNALED",
	"ANY_SAMPLES_PASSED",
	"ANY_SAMPLES_PASSED_CONSERVATIVE",
	"BLUE",
	"BUFFER_ACCESS_FLAGS",
	"BUFFER_MAP_LENGTH",
	"BUFFER_MAP_OFFSET",
	"BUFFER_MAPPED",
	"BUFFER_MAP_POINTER",
	"COLOR",
	"COLOR_ATTACHMENT10",
	"COLOR_ATTACHMENT1",
	"COLOR_ATTACHMENT11",
	"COLOR_ATTACHMENT12",
	"COLOR_ATTACHMENT13",
	"COLOR_ATTACHMENT14",
	"COLOR_ATTACHMENT15",
	"COLOR_ATTACHMENT2",
	"COLOR_ATTACHMENT3",
	"COLOR_ATTACHMENT4",
	"COLOR_ATTACHMENT5",
	"COLOR_ATTACHMENT6",
	"COLOR_ATTACHMENT7",
	"COLOR_ATTACHMENT8",
	"COLOR_ATTACHMENT9",
	"COMPARE_REF_TO_TEXTURE",
	"COMPRESSED_R11_EAC",
	"COMPRESSED_RG11_EAC",
	"COMPRESSED_RGB8_ETC2",
	"COMPRESSED_RGB8_PUNCHTHROUGH_ALPHA1_ETC2",
	"COMPRESSED_RGBA8_ETC2_EAC",
	"COMPRESSED_SIGNED_R11_EAC",
	"COMPRESSED_SIGNED_RG11_EAC",
	"COMPRESSED_SRGB8_ALPHA8_ETC2_EAC",
	"COMPRESSED_SRGB8_ETC2",
	"COMPRESSED_SRGB8_PUNCHTHROUGH_ALPHA1_ETC2",
	"CONDITION_SATISFIED",
	"COPY_READ_BUFFER",
	"COPY_READ_BUFFER_BINDING",
	"COPY_WRITE_BUFFER",
	"COPY_WRITE_BUFFER_BINDING",
	"CURRENT_QUERY",
	"DEPTH",
	"DEPTH24_STENCIL8",
	"DEPTH32F_STENCIL8",
	"DEPTH_COMPONENT24",
	"DEPTH_COMPONENT32F",
	"DEPTH_STENCIL",
	"DEPTH_STENCIL_ATTACHMENT",
	"DRAW_BUFFER0",
	"DRAW_BUFFER10",
	"DRAW_BUFFER1",
	"DRAW_BUFFER11",
	"DRAW_BUFFER12",
	"DRAW_BUFFER13",
	"DRAW_BUFFER14",
	"DRAW_BUFFER15",
	"DRAW_BUFFER2",
	"DRAW_BUFFER3",
	"DRAW_BUFFER4",
	"DRAW_BUFFER5",
	"DRAW_BUFFER6",
	"DRAW_BUFFER7",
	"DRAW_BUFFER8",
	"DRAW_BUFFER9",
	"DRAW_FRAMEBUFFER",
	"DRAW_FRAMEBUFFER_BINDING",
	"DYNAMIC_COPY",
	"DYNAMIC_READ",
	"FLOAT_32_UNSIGNED_INT_24_8_REV",
	"FLOAT_MAT2x3",
	"FLOAT_MAT2x4",
	"FLOAT_MAT3x2",
	"FLOAT_MAT3x4",
	"FLOAT_MAT4x2",
	"FLOAT_MAT4x3",
	"FRAGMENT_SHADER_DERIVATIVE_HINT",
	"FRAMEBUFFER_ATTACHMENT_ALPHA_SIZE",
	"FRAMEBUFFER_ATTACHMENT_BLUE_SIZE",
	"FRAMEBUFFER_ATTACHMENT_COLOR_ENCODING",
	"FRAMEBUFFER_ATTACHMENT_COMPONENT_TYPE",
	"FRAMEBUFFER_ATTACHMENT_DEPTH_SIZE",
	"FRAMEBUFFER_ATTACHMENT_GREEN_SIZE",
	"FRAMEBUFFER_ATTACHMENT_RED_SIZE",
	"FRAMEBUFFER_ATTACHMENT_STENCIL_SIZE",
	"FRAMEBUFFER_ATTACHMENT_TEXTURE_LAYER",
	"FRAMEBUFFER_DEFAULT",
	"FRAMEBUFFER_INCOMPLETE_MULTISAMPLE",
	"FRAMEBUFFER_UNDEFINED",
	"GREEN",
	"HALF_FLOAT",
	"INT_2_10_10_10_REV",
	"INTERLEAVED_ATTRIBS",
	"INT_SAMPLER_2D",
	"INT_SAMPLER_2D_ARRAY",
	"INT_SAMPLER_3D",
	"INT_SAMPLER_CUBE",
	"INVALID_INDEX",
	"MAJOR_VERSION",
	"MAP_FLUSH_EXPLICIT_BIT",
	"MAP_INVALIDATE_BUFFER_BIT",
	"MAP_INVALIDATE_RANGE_BIT",
	"MAP_READ_BIT",
	"MAP_UNSYNCHRONIZED_BIT",
	"MAP_WRITE_BIT",
	"MAX",
	"MAX_3D_TEXTURE_SIZE",
	"MAX_ARRAY_TEXTURE_LAYERS",
	"MAX_COLOR_ATTACHMENTS",
	"MAX_COMBINED_FRAGMENT_UNIFORM_COMPONENTS",
	"MAX_COMBINED_UNIFORM_BLOCKS",
	"MAX_COMBINED_VERTEX_UNIFORM_COMPONENTS",
	"MAX_DRAW_BUFFERS",
	"MAX_ELEMENT_INDEX",
	"MAX_ELEMENTS_INDICES",
	"MAX_ELEMENTS_VERTICES",
	"MAX_FRAGMENT_INPUT_COMPONENTS",
	"MAX_FRAGMENT_UNIFORM_BLOCKS",
	"MAX_FRAGMENT_UNIFORM_COMPONENTS",
	"MAX_PROGRAM_TEXEL_OFFSET",
	"MAX_SAMPLES",
	"MAX_SERVER_WAIT_TIMEOUT",
	"MAX_TEXTURE_LOD_BIAS",
	"MAX_TRANSFORM_FEEDBACK_INTERLEAVED_COMPONENTS",
	"MAX_TRANSFORM_FEEDBACK_SEPARATE_ATTRIBS",
	"MAX_TRANSFORM_FEEDBACK_SEPARATE_COMPONENTS",
	"MAX_UNIFORM_BLOCK_SIZE",
	"MAX_UNIFORM_BUFFER_BINDINGS",
	"MAX_VARYING_COMPONENTS",
	"MAX_VERTEX_OUTPUT_COMPONENTS",
	"MAX_VERTEX_UNIFORM_BLOCKS",
	"MAX_VERTEX_UNIFORM_COMPONENTS",
	"MIN",
	"MINOR_VERSION",
	"MIN_PROGRAM_TEXEL_OFFSET",
	"NUM_EXTENSIONS",
	"NUM_PROGRAM_BINARY_FORMATS",
	"NUM_SAMPLE_COUNTS",
	"OBJECT_TYPE",
	"PACK_ROW_LENGTH",
	"PACK_SKIP_PIXELS",
	"PACK_SKIP_ROWS",
	"PIXEL_PACK_BUFFER",
	"PIXEL_PACK_BUFFER_BINDING",
	"PIXEL_UNPACK_BUFFER",
	"PIXEL_UNPACK_BUFFER_BINDING",
	"PRIMITIVE_RESTART_FIXED_INDEX",
	"PROGRAM_BINARY_FORMATS",
	"PROGRAM_BINARY_LENGTH",
	"PROGRAM_BINARY_RETRIEVABLE_HINT",
	"QUERY_RESULT",
	"QUERY_RESULT_AVAILABLE",
	"R11F_G11F_B10F",
	"R16F",
	"R16I",
	"R16UI",
	"R32F",
	"R32I",
	"R32UI",
	"R8",
	"R8I",
	"R8_SNORM",
	"R8UI",
	"RASTERIZER_DISCARD",
	"READ_BUFFER",
	"READ_FRAMEBUFFER",
	"READ_FRAMEBUFFER_BINDING",
	"RED",
	"RED_INTEGER",
	"RENDERBUFFER_SAMPLES",
	"RG",
	"RG16F",
	"RG16I",
	"RG16UI",
	"RG32F",
	"RG32I",
	"RG32UI",
	"RG8",
	"RG8I",
	"RG8_SNORM",
	"RG8UI",
	"RGB10_A2",
	"RGB10_A2UI",
	"RGB16F",
	"RGB16I",
	"RGB16UI",
	"RGB32F",
	"RGB32I",
	"RGB32UI",
	"RGB8",
	"RGB8I",
	"RGB8_SNORM",
	"RGB8UI",
	"RGB9_E5",
	"RGBA16F",
	"RGBA16I",
	"RGBA16UI",
	"RGBA32F",
	"RGBA32I",
	"RGBA32UI",
	"RGBA8",
	"RGBA8I",
	"RGBA8_SNORM",
	"RGBA8UI",
	"RGBA_INTEGER",
	"RGB_INTEGER",
	"RG_INTEGER",
	"SAMPLER_2D_ARRAY",
	"SAMPLER_2D_ARRAY_SHADOW",
	"SAMPLER_2D_SHADOW",
	"SAMPLER_3D",
	"SAMPLER_BINDING",
	"SAMPLER_CUBE_SHADOW",
	"SEPARATE_ATTRIBS",
	"SIGNALED",
	"SIGNED_NORMALIZED",
	"SRGB",
	"SRGB8",
	"SRGB8_ALPHA8",
	"STATIC_COPY",
	"STATIC_READ",
	"STENCIL",
	"STREAM_COPY",
	"STREAM_READ",
	"SYNC_CONDITION",
	"SYNC_FENCE",
	"SYNC_FLAGS",
	"SYNC_FLUSH_COMMANDS_BIT",
	"SYNC_GPU_COMMANDS_COMPLETE",
	"SYNC_STATUS",
	"TEXTURE_2D_ARRAY",
	"TEXTURE_3D",
	"TEXTURE_BASE_LEVEL",
	"TEXTURE_BINDING_2D_ARRAY",
	"TEXTURE_BINDING_3D",
	"TEXTURE_COMPARE_FUNC",
	"TEXTURE_COMPARE_MODE",
	"TEXTURE_IMMUTABLE_FORMAT",
	"TEXTURE_IMMUTABLE_LEVELS",
	"TEXTURE_MAX_LEVEL",
	"TEXTURE_MAX_LOD",
	"TEXTURE_MIN_LOD",
	"TEXTURE_SWIZZLE_A",
	"TEXTURE_SWIZZLE_B",
	"TEXTURE_SWIZZLE_G",
	"TEXTURE_SWIZZLE_R",
	"TEXTURE_WRAP_R",
	"TIMEOUT_EXPIRED",
	"TRANSFORM_FEEDBACK",
	"TRANSFORM_FEEDBACK_ACTIVE",
	"TRANSFORM_FEEDBACK_BINDING",
	"TRANSFORM_FEEDBACK_BUFFER",
	"TRANSFORM_FEEDBACK_BUFFER_BINDING",
	"TRANSFORM_FEEDBACK_BUFFER_MODE",
	"TRANSFORM_FEEDBACK_BUFFER_SIZE",
	"TRANSFORM_FEEDBACK_BUFFER_START",
	"TRANSFORM_FEEDBACK_PAUSED",
	"TRANSFORM_FEEDBACK_PRIMITIVES_WRITTEN",
	"TRANSFORM_FEEDBACK_VARYING_MAX_LENGTH",
	"TRANSFORM_FEEDBACK_VARYINGS",
	"UNIFORM_ARRAY_STRIDE",
	"UNIFORM_BLOCK_ACTIVE_UNIFORM_INDICES",
	"UNIFORM_BLOCK_ACTIVE_UNIFORMS",
	"UNIFORM_BLOCK_BINDING",
	"UNIFORM_BLOCK_DATA_SIZE",
	"UNIFORM_BLOCK_INDEX",
	"UNIFORM_BLOCK_NAME_LENGTH",
	"UNIFORM_BLOCK_REFERENCED_BY_FRAGMENT_SHADER",
	"UNIFORM_BLOCK_REFERENCED_BY_VERTEX_SHADER",
	"UNIFORM_BUFFER",
	"UNIFORM_BUFFER_BINDING",
	"UNIFORM_BUFFER_OFFSET_ALIGNMENT",
	"UNIFORM_BUFFER_SIZE",
	"UNIFORM_BUFFER_START",
	"UNIFORM_IS_ROW_MAJOR",
	"UNIFORM_MATRIX_STRIDE",
	"UNIFORM_NAME_LENGTH",
	"UNIFORM_OFFSET",
	"UNIFORM_SIZE",
	"UNIFORM_TYPE",
	"UNPACK_IMAGE_HEIGHT",
	"UNPACK_ROW_LENGTH",
	"UNPACK_SKIP_IMAGES",
	"UNPACK_SKIP_PIXELS",
	"UNPACK_SKIP_ROWS",
	"UNSIGNALED",
	"UNSIGNED_INT_10F_11F_11F_REV",
	"UNSIGNED_INT_2_10_10_10_REV",
	"UNSIGNED_INT_24_8",
	"UNSIGNED_INT_5_9_9_9_REV",
	"UNSIGNED_INT_SAMPLER_2D",
	"UNSIGNED_INT_SAMPLER_2D_ARRAY",
	"UNSIGNED_INT_SAMPLER_3D",
	"UNSIGNED_INT_SAMPLER_CUBE",
	"UNSIGNED_INT_VEC2",
	"UNSIGNED_INT_VEC3",
	"UNSIGNED_INT_VEC4",
	"UNSIGNED_NORMALIZED",
	"VERTEX_ARRAY_BINDING",
	"VERTEX_ATTRIB_ARRAY_DIVISOR",
	"VERTEX_ATTRIB_ARRAY_INTEGER",
	"WAIT_FAILED",
}

var outfile = flag.String("o", "", "result will be written to the file instead of stdout.")

var fset = new(token.FileSet)

func typeString(t ast.Expr) string {
	buf := new(bytes.Buffer)
	printer.Fprint(buf, fset, t)
	return buf.String()
}

func typePrinter(t string) string {
	switch t {
	case "[]float32", "[]byte":
		return "len(%d)"
	}
	return "%v"
}

func typePrinterArg(t, name string) string {
	switch t {
	case "[]float32", "[]byte":
		return "len(" + name + ")"
	}
	return name
}

func die(err error) {
	fmt.Fprintf(os.Stderr, err.Error())
	os.Exit(1)
}

func main() {
	flag.Parse()

	f, err := parser.ParseFile(fset, "consts.go", nil, parser.ParseComments)
	if err != nil {
		die(err)
	}
	entries := enum(f)

	f, err = parser.ParseFile(fset, "gl.go", nil, parser.ParseComments)
	if err != nil {
		die(err)
	}

	buf := new(bytes.Buffer)

	fmt.Fprint(buf, preamble)

	fmt.Fprintf(buf, "func (v Enum) String() string {\n")
	fmt.Fprintf(buf, "\tswitch v {\n")
	for _, e := range dedup(entries) {
		fmt.Fprintf(buf, "\tcase 0x%x: return %q\n", e.value, e.name)
	}
	fmt.Fprintf(buf, "\t%s\n", `default: return fmt.Sprintf("gl.Enum(0x%x)", uint32(v))`)
	fmt.Fprintf(buf, "\t}\n")
	fmt.Fprintf(buf, "}\n\n")

	for _, d := range f.Decls {
		// Before:
		// func (ctx *context) StencilMask(mask uint32) {
		//	C.glStencilMask(C.GLuint(mask))
		// }
		//
		// After:
		// func (ctx *context) StencilMask(mask uint32) {
		// 	defer func() {
		// 		errstr := ctx.errDrain()
		// 		log.Printf("gl.StencilMask(%v) %v", mask, errstr)
		//	}()
		//	C.glStencilMask(C.GLuint(mask))
		// }
		fn, ok := d.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if fn.Recv == nil || fn.Recv.List[0].Names[0].Name != "ctx" {
			continue
		}
		tname := "<unknown>"
		t := fn.Recv.List[0].Type
		if star, ok := t.(*ast.StarExpr); ok {
			tname = "*" + star.X.(*ast.Ident).Name
		} else if t, ok := t.(*ast.Ident); ok {
			tname = t.Name
		}

		var (
			params      []string
			paramTypes  []string
			results     []string
			resultTypes []string
		)

		// Print function signature.
		fmt.Fprintf(buf, "func (ctx %s) %s(", tname, fn.Name.Name)
		for i, p := range fn.Type.Params.List {
			if i > 0 {
				fmt.Fprint(buf, ", ")
			}
			ty := typeString(p.Type)
			for i, n := range p.Names {
				if i > 0 {
					fmt.Fprint(buf, ", ")
				}
				fmt.Fprintf(buf, "%s ", n.Name)
				params = append(params, n.Name)
				paramTypes = append(paramTypes, ty)
			}
			fmt.Fprint(buf, ty)
		}
		fmt.Fprintf(buf, ") (")
		if fn.Type.Results != nil {
			for i, r := range fn.Type.Results.List {
				if i > 0 {
					fmt.Fprint(buf, ", ")
				}
				ty := typeString(r.Type)
				if len(r.Names) == 0 {
					name := fmt.Sprintf("r%d", i)
					fmt.Fprintf(buf, "%s ", name)
					results = append(results, name)
					resultTypes = append(resultTypes, ty)
				}
				for i, n := range r.Names {
					if i > 0 {
						fmt.Fprint(buf, ", ")
					}
					fmt.Fprintf(buf, "%s ", n.Name)
					results = append(results, n.Name)
					resultTypes = append(resultTypes, ty)
				}
				fmt.Fprint(buf, ty)
			}
		}
		fmt.Fprintf(buf, ") {\n")

		// gl.GetError is used by errDrain, which will be made part of
		// all functions. So do not apply it to gl.GetError to avoid
		// infinite recursion.
		skip := fn.Name.Name == "GetError"

		if !skip {
			// Insert a defer block for tracing.
			fmt.Fprintf(buf, "defer func() {\n")
			fmt.Fprintf(buf, "\terrstr := ctx.errDrain()\n")
			switch fn.Name.Name {
			case "GetUniformLocation", "GetAttribLocation":
				fmt.Fprintf(buf, "\tr0.name = name\n")
			}
			fmt.Fprintf(buf, "\tlog.Printf(\"gl.%s(", fn.Name.Name)
			for i, p := range paramTypes {
				if i > 0 {
					fmt.Fprint(buf, ", ")
				}
				fmt.Fprint(buf, typePrinter(p))
			}
			fmt.Fprintf(buf, ") ")
			if len(resultTypes) > 1 {
				fmt.Fprint(buf, "(")
			}
			for i, r := range resultTypes {
				if i > 0 {
					fmt.Fprint(buf, ", ")
				}
				fmt.Fprint(buf, typePrinter(r))
			}
			if len(resultTypes) > 1 {
				fmt.Fprint(buf, ") ")
			}
			fmt.Fprintf(buf, "%%v\"")
			for i, p := range paramTypes {
				fmt.Fprintf(buf, ", %s", typePrinterArg(p, params[i]))
			}
			for i, r := range resultTypes {
				fmt.Fprintf(buf, ", %s", typePrinterArg(r, results[i]))
			}
			fmt.Fprintf(buf, ", errstr)\n")
			fmt.Fprintf(buf, "}()\n")
		}

		// Print original body of function.
		for _, s := range fn.Body.List {
			if c := enqueueCall(s); c != nil {
				c.Fun.(*ast.SelectorExpr).Sel.Name = "enqueueDebug"
				setEnqueueBlocking(c)
			}
			printer.Fprint(buf, fset, s)
			fmt.Fprintf(buf, "\n")
		}
		fmt.Fprintf(buf, "}\n\n")
	}

	b, err := format.Source(buf.Bytes())
	if err != nil {
		os.Stdout.Write(buf.Bytes())
		die(err)
	}

	if *outfile == "" {
		os.Stdout.Write(b)
		return
	}
	if err := ioutil.WriteFile(*outfile, b, 0666); err != nil {
		die(err)
	}
}

func enqueueCall(stmt ast.Stmt) *ast.CallExpr {
	exprStmt, ok := stmt.(*ast.ExprStmt)
	if !ok {
		return nil
	}
	call, ok := exprStmt.X.(*ast.CallExpr)
	if !ok {
		return nil
	}
	fun, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}
	if fun.Sel.Name != "enqueue" {
		return nil
	}
	return call
}

func setEnqueueBlocking(c *ast.CallExpr) {
	lit := c.Args[0].(*ast.CompositeLit)
	for _, elt := range lit.Elts {
		kv := elt.(*ast.KeyValueExpr)
		if kv.Key.(*ast.Ident).Name == "blocking" {
			kv.Value = &ast.Ident{Name: "true"}
			return
		}
	}
	lit.Elts = append(lit.Elts, &ast.KeyValueExpr{
		Key: &ast.Ident{
			NamePos: lit.Rbrace,
			Name:    "blocking",
		},
		Value: &ast.Ident{Name: "true"},
	})
}

const preamble = `// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated from gl.go using go generate. DO NOT EDIT.
// See doc.go for details.

// +build darwin linux openbsd windows
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

`

type entry struct {
	name  string
	value uint64
}

func genWhitelistMap(whitelist []string) map[string]bool {
	m := make(map[string]bool)
	for _, v := range enumWhitelist {
		m[v] = true
	}
	return m
}

// enum builds a list of all GL constants that make up the gl.Enum type.
func enum(f *ast.File) []entry {
	var entries []entry
	whitelist := genWhitelistMap(enumWhitelist)
	for _, d := range f.Decls {
		gendecl, ok := d.(*ast.GenDecl)
		if !ok {
			continue
		}
		if gendecl.Tok != token.CONST {
			continue
		}
		for _, s := range gendecl.Specs {
			v, ok := s.(*ast.ValueSpec)
			if !ok {
				continue
			}
			if len(v.Names) != 1 || len(v.Values) != 1 {
				continue
			}
			if _, ok := whitelist[v.Names[0].Name]; !ok {
				continue
			}
			val, err := strconv.ParseUint(v.Values[0].(*ast.BasicLit).Value, 0, 64)
			if err != nil {
				log.Fatalf("enum %s: %v", v.Names[0].Name, err)
			}
			entries = append(entries, entry{v.Names[0].Name, val})
		}
	}
	return entries
}

func dedup(entries []entry) []entry {
	// Find all duplicates. Use "%d" as the name of any value with duplicates.
	seen := make(map[uint64]int)
	for _, e := range entries {
		seen[e.value]++
	}
	var dedup []entry
	for _, e := range entries {
		switch seen[e.value] {
		case 0: // skip, already here
		case 1:
			dedup = append(dedup, e)
		default:
			// value is duplicated
			dedup = append(dedup, entry{fmt.Sprintf("%d", e.value), e.value})
			seen[e.value] = 0
		}
	}
	return dedup
}
