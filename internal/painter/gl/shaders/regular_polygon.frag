#version 110

uniform vec2 frame;
uniform vec4 bounds;
uniform float edgeSoftness;

uniform float outerRadius;
uniform float angle;
uniform float sides;

uniform vec4 fillColor;
uniform float cornerRadius;
uniform float strokeWidth;
uniform vec4 strokeColor;

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
    vec4 frag_rect_coords = vec4(bounds[0], bounds[2], frame.y - bounds[3], frame.y - bounds[1]);
    vec2 vec_centered_pos = (gl_FragCoord.xy - vec2(frag_rect_coords[0] + frag_rect_coords[1], frag_rect_coords[2] + frag_rect_coords[3]) * 0.5);

    vec_centered_pos = rotate(radians(angle)) * vec_centered_pos;
    float dist = regular_distance(vec_centered_pos, outerRadius - cornerRadius, int(sides)) - cornerRadius;
    vec4 final_color = fillColor;

    if (strokeWidth > 0.0)
    {
        // create a mask for the fill area (inside, shrunk by stroke width)
        float fill_mask = smoothstep(-strokeWidth + edgeSoftness, -strokeWidth - edgeSoftness, dist);

        // combine fill mask and colors (fill + stroke)
        final_color = mix(strokeColor, fillColor, fill_mask);
    }

    // smooth edges
    float final_alpha = smoothstep(edgeSoftness, -edgeSoftness, dist);

    // apply the final alpha to the combined color
    gl_FragColor = vec4(final_color.rgb, final_color.a * final_alpha);
}
