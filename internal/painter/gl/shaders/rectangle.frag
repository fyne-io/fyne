#version 110

/* scaled params */
uniform vec2 frame;
uniform vec4 bounds; //x1 [0], y1 [1], x2 [2], y2 [3]; coords in the frame
uniform float strokeWidth;
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
    vec4 color = fillColor;

    if (addShadow == 1.0)
    {
        vec2 frag_pos = gl_FragCoord.xy + vec2(-shadowOffset.x, shadowOffset.y);
        vec2 center = vec2((bounds[0] + bounds[2]) * 0.5, frame.y - (bounds[1] + bounds[3]) * 0.5);
        // expand/contract rectangle bounds by spread on all sides
        vec2 half_size = vec2(bounds[2] - bounds[0], bounds[3] - bounds[1]) * 0.5 + vec2(shadowSpread);

        vec2 d = abs(frag_pos - center) - half_size;
        float distance_shadow = smoothstep(-shadowBlurRadius * 0.5, shadowBlurRadius * 0.5, length(max(d, 0.0)) + min(max(d.x, d.y), 0.0));
        float shadow_alpha = shadowColor.a * (1.0 - distance_shadow);

        if (shadowType == 0.0)
        {
            // remove shadow inside rectangle (uses original rect, not spread rect)
            vec2 frag_pos = gl_FragCoord.xy;
            float d_h = min(frag_pos.x - bounds[0], bounds[2] - frag_pos.x);
            float d_v = min(frag_pos.y - frame.y + bounds[3], frame.y - bounds[1] - frag_pos.y);
            float mask = smoothstep(0.0, -0.5, min(d_h, d_v));
            shadow_alpha *= mask;
        }

        if (gl_FragCoord.x > bounds[2]){
            color[3] = 0.0;
        } else if (gl_FragCoord.x < bounds[0]){
            color[3] = 0.0;
        } else if (gl_FragCoord.y < frame.y - bounds[3]){
            color[3] = 0.0;
        } else if (gl_FragCoord.y > frame.y - bounds[1]){
            color[3] = 0.0;
        }

        color = blend_shadow(color, vec4(shadowColor.rgb, shadow_alpha));
    }

    // discard if outside rectangle coords, necessary to draw thin stroke and mitigate inconsistent borders issue
    if (gl_FragCoord.x < bounds[0] || gl_FragCoord.x > bounds[2] || gl_FragCoord.y < frame.y - bounds[3] || gl_FragCoord.y > frame.y - bounds[1])
    {
        if (addShadow == 0.0)
        {
            discard;
        }
    }
    else
    {
        if (gl_FragCoord.x >= bounds[2] - strokeWidth)
        {
            color = strokeColor;
        }
        else if (gl_FragCoord.x <= bounds[0] + strokeWidth)
        {
            color = strokeColor;
        }
        else if (gl_FragCoord.y <= frame.y - bounds[3] + strokeWidth)
        {
            color = strokeColor;
        }
        else if (gl_FragCoord.y >= frame.y - bounds[1] - strokeWidth)
        {
            color = strokeColor;
        }
    }

    gl_FragColor = color;
}
