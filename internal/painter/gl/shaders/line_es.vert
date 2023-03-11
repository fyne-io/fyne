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
attribute vec2 normal;
    
uniform float lineWidth;

varying vec2 delta;

void main() {
    delta = normal * lineWidth;

    gl_Position = vec4(vert + delta, 0, 1);
}
