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

void main() {
    vert;  //Workaround, because WebGL optimizes away attributes unused
    gl_Position = vec4(normal, 0, 1);
}
