//go:build android || ios || mobile
// +build android ios mobile

package gl

// Contains Khronos OpenGL API specification constants.
const (
	FALSE               = 0
	TRUE                = 1
	TRIANGLES           = 0x0004
	TRIANGLE_STRIP      = 0x0005
	SRC_ALPHA           = 0x0302
	ONE_MINUS_SRC_ALPHA = 0x0303
	DEPTH_TEST          = 0x0B71
	BLEND               = 0x0BE2
	SCISSOR_TEST        = 0x0C11
	TEXTURE_2D          = 0x0DE1

	UNSIGNED_BYTE = 0x1401
	FLOAT         = 0x1406
	RGBA          = 0x1908

	NEAREST            = 0x2600
	LINEAR             = 0x2601
	TEXTURE_MAG_FILTER = 0x2800
	TEXTURE_MIN_FILTER = 0x2801
	TEXTURE_WRAP_S     = 0x2802
	TEXTURE_WRAP_T     = 0x2803

	CONSTANT_ALPHA           = 0x8003
	ONE_MINUS_CONSTANT_ALPHA = 0x8004
	CLAMP_TO_EDGE            = 0x812F
	TEXTURE0                 = 0x84C0
	DYNAMIC_DRAW             = 0x88E8
	FRAGMENT_SHADER          = 0x8B30
	VERTEX_SHADER            = 0x8B31
	ARRAY_BUFFER             = 0x8892
	COMPILE_STATUS           = 0x8B81
	INFO_LOG_LENGTH          = 0x8B84
	SHADER_SOURCE_LENGTH     = 0x8B88

	DEPTH_BUFFER_BIT = 0x00000100
	COLOR_BUFFER_BIT = 0x00004000
)
