#version 110

// Coverage glyphs are tinted by color; colour glyphs are drawn as rasterised.

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
