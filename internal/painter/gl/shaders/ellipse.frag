#version 110

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
/* shadow params*/
uniform float add_shadow;
uniform float shadow_blur_radius;
uniform float shadow_spread;
uniform vec2 shadow_offset;
uniform vec4 shadow_color;
uniform float shadow_type;

mat2 rotate(float a)
{
    float s = sin(-a);
    float c = cos(-a);
    return mat2(c, -s, s, c);
}

float calc_distance(vec2 p, vec2 r)
{
    const float eps = 1e-3;
    r = max(r, eps);
    vec2 f = p / r;
    return (dot(f, f) - 1.0) / max(length(2.0 * f / r), eps);
}

vec4 blend_shadow(vec4 color, vec4 shadow)
{
    float alpha = color.a + shadow.a * (1.0 - color.a);
    return vec4(
        (color.rgb * color.a + shadow.rgb * shadow.a * (1.0 - color.a)) / alpha,
        alpha
    );
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
    final_color = vec4(final_color.rgb, final_color.a * final_alpha);

    if (add_shadow == 1.0)
    {
        // use ellipse radii by default, expand/contract by spread
        vec2 shadow_radius = radius;
        if (shadow_spread != 0.0)
        {
            shadow_radius = max(radius + shadow_spread, 0.0);
        }

        float blur_inset = shadow_blur_radius * 0.5;
        shadow_radius = max(shadow_radius - blur_inset, 0.0);

        // flip the shadow offset to get the correct shadow position
        // negative offset-x value places the shadow to the left of the element. Negative offset-y value places the shadow above the element
        vec2 shadow_offset_corrected = vec2(-shadow_offset.x, shadow_offset.y);
        float distance_shadow = calc_distance(vec_centered_pos + shadow_offset_corrected, shadow_radius);
        float shadow_alpha = shadow_color.a * (1.0 - smoothstep(-edge_softness, shadow_blur_radius + edge_softness, distance_shadow));

        if (shadow_type == 0.0)
        {
            // remove shadow inside the ellipse
            float mask = smoothstep(-2.0 * edge_softness, 0.0, dist);
            shadow_alpha *= mask;
        }

        final_color = blend_shadow(final_color, vec4(shadow_color.rgb, shadow_alpha));
    }

    gl_FragColor = final_color;
}
