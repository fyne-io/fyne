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

varying vec2 fragTexCoord;

void main() {
    vec4 texColor = texture2D(tex, fragTexCoord);
    if(texColor.a < 0.01)
        discard;
    gl_FragColor = texColor;
}
