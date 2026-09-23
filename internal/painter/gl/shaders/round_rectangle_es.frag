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
uniform vec2 frame;
uniform vec4 bounds; //x1 [0], y1 [1], x2 [2], y2 [3]; coords in the frame
uniform float strokeWidthHalf;
uniform vec2 rectSizeHalf;
uniform vec4 radius;
uniform float edgeSoftness;
/* colors params*/
uniform vec4 fillColor;
uniform vec4 strokeColor;
/* shadow params*/
uniform float addShadow;
uniform float shadowBlurRadius;
uniform float shadowSpread;
uniform vec2 shadowOffset;
uniform vec4 shadowColor;
uniform float shadowType;

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

void main()
{
    vec4 frag_rect_coords = vec4(bounds[0], bounds[2], frame.y - bounds[3], frame.y - bounds[1]);
    vec2 vec_centered_pos = (gl_FragCoord.xy - vec2(frag_rect_coords[0] + frag_rect_coords[1], frag_rect_coords[2] + frag_rect_coords[3]) * 0.5);

    float distance;
    float max_radius = max(max(radius.x, radius.y), max(radius.z, radius.w));
    vec4 final_color = fillColor;
    float final_alpha;

    // subtract a small threshold value to avoid calling calc_distance_all_quadrants when the largest corner radius is very close to half the length of the rectangle's shortest edge
    bool calc_all_quadrants = max_radius - 0.9 > min(rectSizeHalf.x, rectSizeHalf.y) + strokeWidthHalf;
    if (calc_all_quadrants)
    {
        // at least one corner radius is larger than half of the shorter edge
        distance = calc_distance_all_quadrants(vec_centered_pos, rectSizeHalf + strokeWidthHalf, radius);
        final_alpha = 1.0 - smoothstep(-edgeSoftness, edgeSoftness, distance);

        if (strokeWidthHalf > 0.0)
        {
            float color_blend = 1.0 - smoothstep(strokeWidthHalf * 2.0 - edgeSoftness, strokeWidthHalf * 2.0 + edgeSoftness, abs(distance));
            final_color = mix(fillColor, strokeColor, color_blend);
        }
    }
    else
    {
        distance = calc_distance(vec_centered_pos, rectSizeHalf, radius - strokeWidthHalf);
        final_alpha = 1.0 - smoothstep(strokeWidthHalf - edgeSoftness, strokeWidthHalf + edgeSoftness, distance);

        if (strokeWidthHalf > 0.0)
        {
            float color_blend = smoothstep(-strokeWidthHalf - edgeSoftness, -strokeWidthHalf + edgeSoftness, distance);
            final_color = mix(fillColor, strokeColor, color_blend);
        }
    }

    // final color
    final_color = vec4(final_color.rgb, final_color.a * final_alpha);

    if (addShadow == 1.0)
    {
        // use rectangle size by default
        vec2 shadow_size = rectSizeHalf + strokeWidthHalf;
        vec4 shadow_radius = radius;

        if (shadowSpread != 0.0)
        {
            // expand/contract by spread, adjust radii to match
            vec2 original_size = shadow_size;
            shadow_size = max(original_size + shadowSpread, 0.0);
            float ratio_x = (original_size.x > 0.0) ? (shadow_size.x / original_size.x) : 1.0;
            float ratio_y = (original_size.y > 0.0) ? (shadow_size.y / original_size.y) : 1.0;
            // scale all corner radii proportionally, use minimum ratio so radius never exceeds the shorter adjacent edge
            shadow_radius = max(radius * min(ratio_x, ratio_y), 0.0);
        }

        float blur_inset = shadowBlurRadius * 0.5;
        shadow_size = max(shadow_size - blur_inset, 0.0);
        shadow_radius = max(shadow_radius - blur_inset, 0.0);

        // apply shadow effect
        float distance_shadow;
        // flip the shadow offset to get the correct shadow position
        // negative offset-x value places the shadow to the left of the element. Negative offset-y value places the shadow above the element
        vec2 shadow_offset_corrected = vec2(-shadowOffset.x, shadowOffset.y);
        if (calc_all_quadrants)
        {
            distance_shadow = calc_distance_all_quadrants(vec_centered_pos + shadow_offset_corrected, shadow_size, shadow_radius);
        }
        else
        {
            distance_shadow = calc_distance(vec_centered_pos + shadow_offset_corrected, shadow_size, shadow_radius);
        }
        float shadow_alpha = shadowColor.a * (1.0 - smoothstep(-edgeSoftness, shadowBlurRadius + edgeSoftness, distance_shadow));

        if (shadowType == 0.0)
        {
            // remove shadow inside rectangle
            float d_shape;
            if (calc_all_quadrants)
            {
                // reuse the previously computed outer distance
                d_shape = distance;
            }
            else
            {
                d_shape = calc_distance(vec_centered_pos, rectSizeHalf + strokeWidthHalf, radius);
            }
            float mask = smoothstep(-2.0 * edgeSoftness, 0.0, d_shape);
            shadow_alpha *= mask;
        }

        final_color = blend_shadow(final_color, vec4(shadowColor.rgb, shadow_alpha));
    }

    gl_FragColor = final_color;
}
