//go:build windows && directx

package dx

import "testing"

func TestParseUserUniforms(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		want    map[string]uint32
		size    uint32
		wantNil bool
	}{{
		name: "scalars pack tightly across registers",
		src: `cbuffer MeshCB : register(b1) {
			float a, b, c;
			float d;
			float e;
		};`,
		want: map[string]uint32{"a": 0, "b": 4, "c": 8, "d": 12, "e": 16},
		size: 32,
	}, {
		name: "a vector is bumped past a register boundary it would straddle",
		src: `cbuffer X : register(b1) {
			float pad0, pad1, pad2;
			float3 v;   // would straddle at offset 12, so lands at 16
			float tail;
		};`,
		want: map[string]uint32{"pad0": 0, "pad1": 4, "pad2": 8, "v": 16, "tail": 28},
		size: 32,
	}, {
		name: "a float4 at a register boundary is not bumped",
		src:  `cbuffer X : register(b1) { float4 v; float t; };`,
		want: map[string]uint32{"v": 0, "t": 16},
		size: 32,
	}, {
		name: "comments are ignored",
		src: `cbuffer X : register(b1) {
			float a; // trailing
			/* block */ float b;
		};`,
		want: map[string]uint32{"a": 0, "b": 4},
		size: 16,
	}, {
		name:    "a cbuffer at another register is not the user one",
		src:     `cbuffer X : register(b2) { float a; };`,
		wantNil: true,
	}, {
		name:    "no cbuffer at all",
		src:     `float4 main(PSIn i) : SV_TARGET { return 0; }`,
		wantNil: true,
	}, {
		// ponytail's stated ceiling: an unhandled type abandons the parse
		// rather than returning offsets that are wrong from that point on.
		name:    "an unsupported type abandons the parse",
		src:     `cbuffer X : register(b1) { float a; float4x4 m; float b; };`,
		wantNil: true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, size := parseUserUniforms(tt.src)
			if tt.wantNil {
				if got != nil || size != 0 {
					t.Fatalf("want no uniforms, got %v size %d", got, size)
				}
				return
			}
			if size != tt.size {
				t.Errorf("size = %d, want %d", size, tt.size)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %d uniforms %v, want %d", len(got), got, len(tt.want))
			}
			for name, off := range tt.want {
				if got[name] != off {
					t.Errorf("%s at offset %d, want %d", name, got[name], off)
				}
			}
		})
	}
}

func TestParseUserUniformsSizeIsAligned(t *testing.T) {
	// D3D11 rejects CreateBuffer for a constant buffer whose size is not a
	// multiple of 16 bytes.
	_, size := parseUserUniforms(`cbuffer X : register(b1) { float a; };`)
	if size == 0 || size%16 != 0 {
		t.Errorf("size %d must be a non-zero multiple of 16", size)
	}
}
