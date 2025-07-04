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
uniform float stroke_width_half;
uniform vec2 rect_size_half;
uniform float radius;
uniform float edge_softness;
/* colors params*/
uniform vec4 fill_color;
uniform vec4 stroke_color;

float calc_distance(vec2 p, vec2 b, float r)
{
    vec2 d = abs(p) - b + vec2(r);
	return min(max(d.x, d.y), 0.0) + length(max(d, 0.0)) - r;   
}

void main() {

    vec4 frag_rect_coords = vec4(rect_coords[0], rect_coords[1], frame_size.y - rect_coords[3], frame_size.y - rect_coords[2]);
    vec2 vec_centered_pos = (gl_FragCoord.xy - vec2(frag_rect_coords[0] + frag_rect_coords[1], frag_rect_coords[2] + frag_rect_coords[3]) * 0.5);

    float distance = calc_distance(vec_centered_pos, rect_size_half, radius - stroke_width_half);

    vec4 from_color = stroke_color; //Always the border color. If no border, this still should be set
    vec4 to_color = stroke_color; //Outside color

    if (stroke_width_half == 0.0)
    {
        from_color = fill_color;
        to_color = fill_color;
    }
    to_color[3] = 0.0; // blend the fill colour to alpha

    if (distance < 0.0)
    {
        to_color = fill_color;
    }

    distance = abs(distance) - stroke_width_half;

    float blend_amount = smoothstep(edge_softness - 1.0, edge_softness + 1.0, distance);

    // final color
    gl_FragColor = mix(from_color, to_color, blend_amount);
}
