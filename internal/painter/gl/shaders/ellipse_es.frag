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
uniform vec2 radius;
uniform float edge_softness;
uniform float angle;
/* colors params*/
uniform vec4 fill_color;
uniform vec4 stroke_color;

mat2 rotate(float a) {
    float s = sin(-a);
    float c = cos(-a);
    return mat2(c, -s, s, c);
}

float calc_distance(vec2 p, vec2 r)
{
    r = max(r, 1e-6);
    vec2 f = p / r;
    return (dot(f, f) - 1.0) / max(length(2.0 * f / r), 1e-6);
}

void main()
{
    vec4 frag_rect_coords = vec4(rect_coords[0], rect_coords[1], frame_size.y - rect_coords[3], frame_size.y - rect_coords[2]);
    vec2 vec_centered_pos = (gl_FragCoord.xy - vec2(frag_rect_coords[0] + frag_rect_coords[1], frag_rect_coords[2] + frag_rect_coords[3]) * 0.5);

    vec_centered_pos = rotate(radians(angle)) * vec_centered_pos;

    float dist = calc_distance(vec_centered_pos, radius);
    vec4 final_color = fill_color;

    if (stroke_width > 0.0)
    {
        vec2 inner_radius = radius - stroke_width;
        float fill_mask = 0.0;
        if (inner_radius.x > 1.0 && inner_radius.y > 1.0)
        {
            // create a mask for the fill area (inside, shrunk by stroke width)
            float dist_inner = calc_distance(vec_centered_pos, inner_radius);
            fill_mask = smoothstep(edge_softness, -edge_softness, dist_inner);
        }

        // combine fill mask and colors (fill + stroke)
        final_color = mix(stroke_color, fill_color, fill_mask);
    }

    // smooth edges
    float final_alpha = smoothstep(edge_softness, -edge_softness, dist);

    // apply the final alpha to the combined color
    gl_FragColor = vec4(final_color.rgb, final_color.a * final_alpha);
}
