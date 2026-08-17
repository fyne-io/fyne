#version 110

// The atlas holds one byte of glyph coverage per pixel rather than colour, so a
// single bitmap serves every colour a glyph is drawn in and the tint is applied
// here. color is alpha premultiplied, matching the blend the painter sets up.

uniform sampler2D tex;
uniform vec4 color;

varying vec2 fragTexCoord;

void main() {
    vec4 texColor = color * texture2D(tex, fragTexCoord).a;
    if (texColor.a < 0.01)
        discard;

    gl_FragColor = texColor;
}
