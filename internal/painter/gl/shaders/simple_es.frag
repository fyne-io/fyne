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
    float alpha = 1.0;
    if (cornerRadius > 0.5) {
        vec2 pos = fragTexCoord * size;
        vec2 halfSize = size * 0.5;
        float dist = length(max(abs(pos - halfSize) - halfSize + cornerRadius, 0.0)) - cornerRadius;
        alpha = 1.0 - smoothstep(-1.0, 1.0, dist);
    }

    vec4 texColor = texture2D(tex, fragTexCoord);
    texColor.a *= fragAlpha * alpha;
    texColor.rgb *= fragAlpha * alpha;

    if(texColor.a < 0.01)
        discard;
    gl_FragColor = texColor;
}
