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
uniform float strokeWidth;
uniform vec2 radius;
uniform float edgeSoftness;
uniform float angle;
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
    vec4 frag_rect_coords = vec4(bounds[0], bounds[2], frame.y - bounds[3], frame.y - bounds[1]);
    vec2 vec_centered_pos = (gl_FragCoord.xy - vec2(frag_rect_coords[0] + frag_rect_coords[1], frag_rect_coords[2] + frag_rect_coords[3]) * 0.5);

    vec_centered_pos = rotate(radians(angle)) * vec_centered_pos;

    float dist = calc_distance(vec_centered_pos, radius);
    vec4 final_color = fillColor;

    if (strokeWidth > 0.0)
    {
        vec2 innerRadius = radius - strokeWidth;
        float fill_mask = 0.0;
        if (innerRadius.x > 1.0 && innerRadius.y > 1.0)
        {
            // create a mask for the fill area (inside, shrunk by stroke width)
            float dist_inner = calc_distance(vec_centered_pos, innerRadius);
            fill_mask = smoothstep(edgeSoftness, -edgeSoftness, dist_inner);
        }

        // combine fill mask and colors (fill + stroke)
        final_color = mix(strokeColor, fillColor, fill_mask);
    }

    // smooth edges
    float final_alpha = smoothstep(edgeSoftness, -edgeSoftness, dist);

    // apply the final alpha to the combined color
    final_color = vec4(final_color.rgb, final_color.a * final_alpha);

    if (addShadow == 1.0)
    {
        // use ellipse radii by default, expand/contract by spread
        vec2 shadow_radius = radius;
        if (shadowSpread != 0.0)
        {
            shadow_radius = max(radius + shadowSpread, 0.0);
        }

        float blur_inset = shadowBlurRadius * 0.5;
        shadow_radius = max(shadow_radius - blur_inset, 0.0);

        // flip the shadow offset to get the correct shadow position
        // negative offset-x value places the shadow to the left of the element. Negative offset-y value places the shadow above the element
        vec2 shadow_offset_corrected = vec2(-shadowOffset.x, shadowOffset.y);
        float distance_shadow = calc_distance(vec_centered_pos + shadow_offset_corrected, shadow_radius);
        float shadow_alpha = shadowColor.a * (1.0 - smoothstep(-edgeSoftness, shadowBlurRadius + edgeSoftness, distance_shadow));

        if (shadowType == 0.0)
        {
            // remove shadow inside the ellipse
            float mask = smoothstep(-2.0 * edgeSoftness, 0.0, dist);
            shadow_alpha *= mask;
        }

        final_color = blend_shadow(final_color, vec4(shadowColor.rgb, shadow_alpha));
    }

    gl_FragColor = final_color;
}
