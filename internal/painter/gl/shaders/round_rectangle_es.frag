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
uniform vec4 radius;
uniform float edge_softness;
/* colors params*/
uniform vec4 fill_color;
uniform vec4 stroke_color;
/* shadow params*/
uniform float add_shadow;
uniform float shadow_softness;
uniform vec2 shadow_offset;
uniform vec4 shadow_color;
uniform float shadow_type;

// distance is calculated for a single quadrant
// returns invalid output if corner radius exceed half of the shorter edge
float calc_distance(vec2 p, vec2 b, vec4 r)
{
    r.xy = (p.x > 0.0) ? r.xy : r.zw;
    r.x  = (p.y > 0.0) ? r.x  : r.y;

    vec2 d = abs(p) - b + r.x;
    return min(max(d.x, d.y), 0.0) + length(max(d, 0.0)) - r.x;
}

// distance is calculated for all necessary quadrants
// corner radius may exceed half of the shorter edge
float calc_distance_all_quadrants(vec2 p, vec2 size, vec4 radius)
{
    vec2 d = abs(p) - size;
    float dist = length(max(d, 0.0)) + min(max(d.x, d.y), 0.0);

    // top-left corner
    vec2 p_tl = p - vec2(radius.z - size.x, size.y - radius.z);
    if (p_tl.x < 0.0 && p_tl.y > 0.0) dist = max(dist, length(p_tl) - radius.z);

    // top-right corner
    vec2 p_tr = p - vec2(size.x - radius.x, size.y - radius.x);
    if (p_tr.x > 0.0 && p_tr.y > 0.0) dist = max(dist, length(p_tr) - radius.x);

    // bottom-right corner
    vec2 p_br = p - vec2(size.x - radius.y, radius.y - size.y);
    if (p_br.x > 0.0 && p_br.y < 0.0) dist = max(dist, length(p_br) - radius.y);

    // bottom-left corner
    vec2 p_bl = p - vec2(radius.w - size.x, radius.w - size.y);
    if (p_bl.x < 0.0 && p_bl.y < 0.0) dist = max(dist, length(p_bl) - radius.w);

    return dist;
}

vec4 blend_shadow(vec4 color, vec4 shadow)
{
    float alpha = color.a + shadow.a * (1.0 - color.a);
    return vec4(
        (color.rgb * color.a + shadow.rgb * shadow.a * (1.0 - color.a)) / alpha,
        alpha
    );
}

void main() {
    vec4 frag_rect_coords = vec4(rect_coords[0], rect_coords[1], frame_size.y - rect_coords[3], frame_size.y - rect_coords[2]);
    vec2 vec_centered_pos = (gl_FragCoord.xy - vec2(frag_rect_coords[0] + frag_rect_coords[1], frag_rect_coords[2] + frag_rect_coords[3]) * 0.5);

    float distance;
    float max_radius = max(max(radius.x, radius.y), max(radius.z, radius.w));
    vec4 final_color = fill_color;
    float final_alpha;

    // subtract a small threshold value to avoid calling calc_distance_all_quadrants when the largest corner radius is very close to half the length of the rectangle's shortest edge
    bool calc_all_quadrants = max_radius - 0.9 > min(rect_size_half.x, rect_size_half.y) + stroke_width_half;
    if (calc_all_quadrants)
    {
        // at least one corner radius is larger than half of the shorter edge
        distance = calc_distance_all_quadrants(vec_centered_pos, rect_size_half + stroke_width_half, radius);
        final_alpha = 1.0 - smoothstep(-edge_softness, edge_softness, distance);

        if (stroke_width_half > 0.0)
        {
            float color_blend = 1.0 - smoothstep(stroke_width_half * 2.0 - edge_softness, stroke_width_half * 2.0 + edge_softness, abs(distance));
            final_color = mix(fill_color, stroke_color, color_blend);
        }
    }
    else
    {
        distance = calc_distance(vec_centered_pos, rect_size_half, radius - stroke_width_half);
        final_alpha = 1.0 - smoothstep(stroke_width_half - edge_softness, stroke_width_half + edge_softness, distance);

        if (stroke_width_half > 0.0)
        {
            float color_blend = smoothstep(-stroke_width_half - edge_softness, -stroke_width_half + edge_softness, distance);
            final_color = mix(fill_color, stroke_color, color_blend);
        }
    }

    // final color
    final_color = vec4(final_color.rgb, final_color.a * final_alpha);

    if (add_shadow == 1.0)
    {
        // Apply shadow effect
        float distance_shadow;
        if (calc_all_quadrants)
        {
            distance_shadow = smoothstep(0.0, shadow_softness, calc_distance_all_quadrants(vec_centered_pos + shadow_offset, rect_size_half + stroke_width_half, radius));
        }
        else
        {
            distance_shadow = smoothstep(0.0, shadow_softness, calc_distance(vec_centered_pos + shadow_offset, rect_size_half + stroke_width_half, radius));
        }
        float shadow_alpha = shadow_color.a * (1.0 - distance_shadow);

        if (shadow_type == 0.0)
        {
            // remove shadow inside rectangle
            float d_shape;
            if (calc_all_quadrants)
            {
                d_shape = calc_distance_all_quadrants(vec_centered_pos, rect_size_half + stroke_width_half, radius);
            }
            else
            {
                d_shape = calc_distance(vec_centered_pos, rect_size_half + stroke_width_half, radius);
            }
            float mask = smoothstep(-2.0 * edge_softness, 0.0, d_shape);
            shadow_alpha *= mask;
        }

        final_color = blend_shadow(final_color, vec4(shadow_color.rgb, shadow_alpha));
    }

    gl_FragColor = final_color;
}
