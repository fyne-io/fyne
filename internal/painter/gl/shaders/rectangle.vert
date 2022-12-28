#version 110

attribute vec2 vert;
attribute vec2 normal;

uniform vec4 frame_size; //size of view/window = x,y; z = pixScale (w not used); 
varying vec4 frame_resolution;

void main() {
    frame_resolution = frame_size;

    gl_Position = vec4(2.0*vert.x/frame_size.x - 1.0, 1.0 - 2.0*vert.y/frame_size.y, 0, 1);
    // gl_Position = vec4(vert, 0, 1);

}
