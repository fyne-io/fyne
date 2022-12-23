#version 110

uniform vec4 rect_coords; //x1 [0], x2 [1], y1 [2], y2 [3]; coords of the rect_frame
uniform float stroke;
uniform vec4 fill_color;
uniform vec4 stroke_color;
varying vec2 frame_resolution;

void main() {

    vec4 color = fill_color;
    
    if (gl_FragCoord.x >= rect_coords[1] - stroke ){
        color = stroke_color;
    } else if (gl_FragCoord.x <= rect_coords[0] + stroke){
        color = stroke_color;
    } else if (gl_FragCoord.y <= frame_resolution.y - rect_coords[3] + stroke ){
        color = stroke_color;
    } else if (gl_FragCoord.y >= frame_resolution.y - rect_coords[2] - stroke ){
        color = stroke_color;
    }

    gl_FragColor = color;

}
