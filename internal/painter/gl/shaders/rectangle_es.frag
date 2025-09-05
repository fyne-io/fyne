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

/* scaled params */
uniform vec2 frame_size;
uniform vec4 rect_coords; //x1 [0], x2 [1], y1 [2], y2 [3]; coords of the rect_frame
uniform float stroke_width;
/* colors params*/
uniform vec4 fill_color;
uniform vec4 stroke_color;


void main() {

    // discard if outside rectangle coords, necessary to draw thin stroke and mitigate inconsistent borders issue
    if (gl_FragCoord.x < rect_coords[0] || gl_FragCoord.x > rect_coords[1] || gl_FragCoord.y < frame_size.y - rect_coords[3] || gl_FragCoord.y > frame_size.y - rect_coords[2]) {
        discard;
    }

    vec4 color = fill_color;
    
    if (gl_FragCoord.x >= rect_coords[1] - stroke_width ){
        color = stroke_color;
    } else if (gl_FragCoord.x <= rect_coords[0] + stroke_width){
        color = stroke_color;
    } else if (gl_FragCoord.y <= frame_size.y - rect_coords[3] + stroke_width ){
        color = stroke_color;
    } else if (gl_FragCoord.y >= frame_size.y - rect_coords[2] - stroke_width ){
        color = stroke_color;
    }

    gl_FragColor = color;
}
