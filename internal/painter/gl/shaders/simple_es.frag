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

uniform sampler2D tex;
uniform float cornerRadius;  // in pixels
uniform vec2 size;           // in pixels: size of the rendered image quad

varying vec2 fragTexCoord;
varying float fragAlpha;

void main() {
    if (cornerRadius > 0.5) {
        float r = cornerRadius;
        vec2 pos = fragTexCoord * size;

        if (pos.x < r && pos.y < r && distance(pos, vec2(r, r)) > r) discard;
        if (pos.x > size.x - r && pos.y < r && distance(pos, vec2(size.x - r, r)) > r) discard;
        if (pos.x < r && pos.y > size.y - r && distance(pos, vec2(r, size.y - r)) > r) discard;
        if (pos.x > size.x - r && pos.y > size.y - r && distance(pos, vec2(size.x - r, size.y - r)) > r) discard;
    }

    vec4 texColor = texture2D(tex, fragTexCoord);
    texColor.a *= fragAlpha;
    texColor.r *= fragAlpha;
    texColor.g *= fragAlpha;
    texColor.b *= fragAlpha;

    if(texColor.a < 0.01)
        discard;
    gl_FragColor = texColor;
}
