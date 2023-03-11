#version 110

attribute vec2 vert;
attribute vec2 normal;
    
uniform float lineWidth;

varying vec2 delta;

void main() {
    delta = normal * lineWidth;

    gl_Position = vec4(vert + delta, 0, 1);
}
