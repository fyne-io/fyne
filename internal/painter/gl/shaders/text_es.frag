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

// Glyphs are rasterised into the atlas in their final colour, so this samples
// and blends without the corner rounding and inset handling that images need.

uniform sampler2D tex;

varying vec2 fragTexCoord;

void main() {
    vec4 texColor = texture2D(tex, fragTexCoord);
    if (texColor.a < 0.01)
        discard;

    gl_FragColor = texColor;
}
