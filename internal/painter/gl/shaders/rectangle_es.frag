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
/* shadow params*/
uniform float add_shadow;
uniform float shadow_softness;
uniform vec2 shadow_offset;
uniform vec4 shadow_color;
uniform float shadow_type;

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
    vec4 color = fill_color;

    if (add_shadow == 1.0)
    {
        vec2 frag_pos = gl_FragCoord.xy + shadow_offset;
        vec2 p = vec2(
            clamp(frag_pos.x, rect_coords[0], rect_coords[1]),
            clamp(frag_pos.y, frame_size.y - rect_coords[3], frame_size.y - rect_coords[2])
        );

        float distance_shadow = smoothstep(0.0, shadow_softness, length(frag_pos - p));
        float shadow_alpha = shadow_color.a * (1.0 - distance_shadow);

        if (shadow_type == 0.0)
        {
            // remove shadow inside rectangle
            vec2 frag_pos = gl_FragCoord.xy;
            float d_h = min(frag_pos.x - rect_coords[0], rect_coords[1] - frag_pos.x);
            float d_v = min(frag_pos.y - frame_size.y + rect_coords[3], frame_size.y - rect_coords[2] - frag_pos.y);
            float mask = smoothstep(0.0, -0.5, min(d_h, d_v));
            shadow_alpha *= mask;
        }

        if (gl_FragCoord.x > rect_coords[1]){
            color[3] = 0.0;
        } else if (gl_FragCoord.x < rect_coords[0]){
            color[3] = 0.0;
        } else if (gl_FragCoord.y < frame_size.y - rect_coords[3]){
            color[3] = 0.0;
        } else if (gl_FragCoord.y > frame_size.y - rect_coords[2]){
            color[3] = 0.0;
        }

        color = blend_shadow(color, vec4(shadow_color.rgb, shadow_alpha));
    }

    // discard if outside rectangle coords, necessary to draw thin stroke and mitigate inconsistent borders issue
    if (gl_FragCoord.x < rect_coords[0] || gl_FragCoord.x > rect_coords[1] || gl_FragCoord.y < frame_size.y - rect_coords[3] || gl_FragCoord.y > frame_size.y - rect_coords[2])
    {
        if (add_shadow == 0.0)
        {
            discard;
        }
    }
    else
    {
        if (gl_FragCoord.x >= rect_coords[1] - stroke_width)
        {
            color = stroke_color;
        }
        else if (gl_FragCoord.x <= rect_coords[0] + stroke_width)
        {
            color = stroke_color;
        }
        else if (gl_FragCoord.y <= frame_size.y - rect_coords[3] + stroke_width)
        {
            color = stroke_color;
        }
        else if (gl_FragCoord.y >= frame_size.y - rect_coords[2] - stroke_width)
        {
            color = stroke_color;
        }
    }

    gl_FragColor = color;
}
