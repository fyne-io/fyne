#version 110

attribute vec2 vert;
attribute vec2 normal;
attribute vec4 vertColor;
attribute vec2 lineStyle; // half width, feather

varying vec2 delta;
varying vec4 color;
varying vec2 style;

void main() {
    delta = normal * lineStyle.x;
    color = vertColor;
    style = lineStyle;

    gl_Position = vec4(vert + delta, 0, 1);
}
