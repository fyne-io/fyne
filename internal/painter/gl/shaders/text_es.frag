#version 100

#ifdef GL_ES
# ifdef GL_FRAGMENT_PRECISION_HIGH
precision highp float;
# else
precision mediump float;
#endif
precision mediump int;
precision lowp sampler2D;
#endif

// The atlas holds one byte of glyph coverage per pixel rather than colour, so a
// single bitmap serves every colour a glyph is drawn in and the tint is applied
// here. color is alpha premultiplied, matching the blend the painter sets up.

uniform sampler2D tex;
uniform vec4 color;
// 0 for glyphs held as coverage and tinted by color, 1 for glyphs that carry
// their own colours, such as emoji, which are drawn as they were rasterised.
uniform float ownColor;

varying vec2 fragTexCoord;

void main() {
    vec4 texel = texture2D(tex, fragTexCoord);
    vec4 texColor = mix(color * texel.a, texel, ownColor);
    if (texColor.a < 0.01)
        discard;

    gl_FragColor = texColor;
}
