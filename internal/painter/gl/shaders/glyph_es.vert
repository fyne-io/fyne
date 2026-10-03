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

attribute vec2 vert;
attribute vec2 vertTexCoord;
attribute vec4 vertColor;

varying vec2 fragTexCoord;
varying vec4 fragColor;

void main() {
    fragTexCoord = vertTexCoord;
    fragColor = vertColor;

    gl_Position = vec4(vert, 0, 1);
}
