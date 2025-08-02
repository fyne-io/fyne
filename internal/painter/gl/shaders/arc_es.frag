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

uniform float inner_radius;
uniform float outer_radius;
uniform float start_angle;
uniform float end_angle;
uniform vec4 fill_color;
uniform float corner_radius;
uniform float stroke_width;
uniform vec4 stroke_color;

const float PI = 3.141592653589793;

// Computes the signed distance for a rounded arc shape.
// Parameters:
//   position      - The 2D coordinate to evaluate (vec2).
//   inner radius  - The inner radius of the arc (float).
//   outer radius  - The outer radius of the arc (float).
//   start rad     - The starting angle of the arc in radians (float).
//   end rad       - The ending angle of the arc in radians (float).
//   corner radius - The radius for rounding the arc's corners (float).
// Returns:
//   The signed distance from the given position to the edge of the rounded arc.
//   Negative values are inside the arc, positive values are outside, and zero is on the edge.
float sdRoundedArc(vec2 p, float r1, float r2, float a0, float a1, float cr)
{
    // center the arc for simpler calculations
    float midAngle = (a0 + a1) / 2.0;
    float arcSpan  = a1 - a0;
    
    float cs = cos(midAngle);
    float sn = sin(midAngle);
    p = mat2(cs, -sn, sn, cs) * p;

    // calculate distance to a rounded box in a pseudo-polar space
    float r = length(p);

    // atan(y, x) for standard angle convention (0 degrees = right)
    float a = atan(p.y, p.x);

    vec2 boxHalfSize = vec2(arcSpan * 0.5 * r, (r2 - r1) * 0.5);
    vec2 q = vec2(a * r, r - (r1 + r2) * 0.5);
    
    // the inner corner radius cannot be larger than the inner radius itself
    float inner_cr = min(cr, r1);
    // the outer corner radius is just cr
    float outer_cr = cr;
    
    // interpolate between inner and outer corner radius based on the radial position
    // 't' goes from 0 (inner) to 1 (outer).
    float t = smoothstep(-boxHalfSize.y, boxHalfSize.y, q.y);
    float effective_cr = mix(inner_cr, outer_cr, t);

    // use the standard SDF for a 2D rounded box with the effective radius
    vec2 d = abs(q) - boxHalfSize + effective_cr;
    return length(max(d, 0.0)) + min(max(d.x, d.y), 0.0) - effective_cr;
}

void main()
{
    vec4 frag_rect_coords = vec4(rect_coords[0], rect_coords[1], frame_size.y - rect_coords[3], frame_size.y - rect_coords[2]);
    vec2 vec_centered_pos = (gl_FragCoord.xy - vec2(frag_rect_coords[0] + frag_rect_coords[1], frag_rect_coords[2] + frag_rect_coords[3]) * 0.5);

    // this logic correctly handles positive Counter-clockwise (CCW) and negative Clockwise (CW) direction
    // the sdRoundedArc function requires start_rad < end_rad to draw a CCW arc
    float start_rad;
    float end_rad;
    if (end_angle >= start_angle) { // CCW
        // the angles are already in the correct order for a CCW arc
        start_rad = radians(start_angle);
        end_rad = radians(end_angle);
    } else { // CW
        // a CW arc from start to end is the same as a CCW arc from end to start
        // swap them to satisfy the start_rad < end_rad requirement of the SDF
        start_rad = radians(end_angle);
        end_rad = radians(start_angle);
    }
    
    // check if the arc is a full circle (360 degrees or more).
    // the sdRoundedArc function creates segment at the start/end angle, which is undesirable for a complete circle
    float d;
    if (abs(end_rad - start_rad) >= 2.0 * PI - 0.001)
    {
        // full circle
        float r = length(vec_centered_pos);
        
        if (inner_radius < 0.5) {
            // no inner radius
            d = r - outer_radius;
        } else {
            float ring_center_radius = (inner_radius + outer_radius) * 0.5;
            float ring_thickness = (outer_radius - inner_radius) * 0.5;
            d = abs(r - ring_center_radius) - ring_thickness;
        }
    }
    else
    {
        d = sdRoundedArc(vec_centered_pos, inner_radius, outer_radius, start_rad, end_rad, corner_radius);
    }
    
    // create a mask for the fill area (inside, shrunk by stroke width)
    float fillMask = smoothstep(edge_softness, -edge_softness, d + stroke_width);

    // combine fill mask and colors (fill + stroke)
    vec4 color = mix(stroke_color, fill_color, fillMask);

    // smooth edges
    float fullMask = smoothstep(edge_softness, -edge_softness, d);
    
    // apply the final alpha to the combined color.
    gl_FragColor = vec4(color.rgb, color.a * fullMask);
}
