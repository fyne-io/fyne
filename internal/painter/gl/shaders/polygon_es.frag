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

uniform vec2 frame_size;
uniform vec4 rect_coords;
uniform float edge_softness;

uniform float outer_radius;
uniform float angle;
uniform float sides;

uniform vec4 fill_color;
uniform float corner_radius;
uniform float stroke_width;
uniform vec4 stroke_color;

const float PI = 3.141592653589793;

mat2 rotate(float angle) {
    float s = sin(-angle);
    float c = cos(-angle);
    return mat2(c, -s, s, c);
}

// The signed distance (float) from the point to the regular polygon's edge
float regular_distance(vec2 p, float r, int s)
{
    float angle = PI / float(s);
    float angle_cos = cos(angle);
    float angle_sin = sin(angle);
    float angular_offset = mod(atan(p.x, p.y), 2.0*angle) - angle;
    vec2 distance = length(p) * vec2(cos(angular_offset), abs(sin(angular_offset))) - r*vec2(angle_cos, angle_sin);
    distance.y += clamp(-distance.y, 0.0, r*angle_sin);
    return length(distance) * sign(distance.x);
}

void main()
{
    vec4 frag_rect_coords = vec4(rect_coords[0], rect_coords[1], frame_size.y - rect_coords[3], frame_size.y - rect_coords[2]);
    vec2 vec_centered_pos = (gl_FragCoord.xy - vec2(frag_rect_coords[0] + frag_rect_coords[1], frag_rect_coords[2] + frag_rect_coords[3]) * 0.5);

    vec_centered_pos = rotate(radians(angle)) * vec_centered_pos;
    float dist = regular_distance(vec_centered_pos, outer_radius - corner_radius, int(sides)) - corner_radius;
    vec4 final_color = fill_color;

    if (stroke_width > 0.0)
    {
        // create a mask for the fill area (inside, shrunk by stroke width)
        float fill_mask = smoothstep(-stroke_width + edge_softness, -stroke_width - edge_softness, dist);

        // combine fill mask and colors (fill + stroke)
        final_color = mix(stroke_color, fill_color, fill_mask);
    }

    // smooth edges
    float final_alpha = smoothstep(edge_softness, -edge_softness, dist);

    // apply the final alpha to the combined color
    gl_FragColor = vec4(final_color.rgb, final_color.a * final_alpha);
}
