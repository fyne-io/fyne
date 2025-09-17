#version 110

// Note: This shader operates in the unit circle coordinate system, where angles are measured from the positive X axis.
// To adapt the arc orientation or coordinate system, adjust the start_angle and end_angle uniforms accordingly.

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

// Computes the signed distance for a rounded arc shape
// Parameters:
//   position      - The 2D coordinate to evaluate (vec2)
//   inner radius  - The inner radius of the arc (float)
//   outer radius  - The outer radius of the arc (float)
//   start rad     - The starting angle of the arc in radians (float)
//   end rad       - The ending angle of the arc in radians (float)
//   corner radius - The radius for rounding the arc's corners (float)
// Returns:
//   The signed distance from the given position to the edge of the rounded arc
//   Negative values are inside the arc, positive values are outside, and zero is on the edge
float sd_rounded_arc(vec2 p, float r1, float r2, float a0, float a1, float cr)
{
    // center the arc for simpler calculations
    float mid_angle = (a0 + a1) / 2.0;
    float arc_span = abs(a1 - a0);
    
    float cs = cos(mid_angle);
    float sn = sin(mid_angle);
    p = mat2(cs, -sn, sn, cs) * p;

    // calculate distance to a rounded box in a pseudo-polar space
    float r = length(p);

    // atan(y, x) for standard angle convention (0 degrees = right)
    float a = atan(p.y, p.x);

    vec2 box_half_size = vec2(arc_span * 0.5 * r, (r2 - r1) * 0.5);
    vec2 q = vec2(a * r, r - (r1 + r2) * 0.5);

    // the inner corner radius is clamped to half of the smaller dimension:
    //   thickness (r2 - r1), to prevent inner/outer corners on the same side from overlapping
    //   inner length (arc_span * r1), to prevent the start/end inner corners from overlapping
    float inner_cr = min(cr, 0.5 * min(r2 - r1, arc_span * r1));
    // the outer corner radius is just cr
    float outer_cr = cr;
    
    // interpolate between inner and outer corner radius based on the radial position
    // 't' goes from 0 (inner) to 1 (outer).
    float t = smoothstep(-box_half_size.y, box_half_size.y, q.y);
    float effective_cr = mix(inner_cr, outer_cr, t);

    // use the standard SDF for a 2D rounded box with the effective radius
    vec2 dist = abs(q) - box_half_size + effective_cr;
    return length(max(dist, 0.0)) + min(max(dist.x, dist.y), 0.0) - effective_cr;
}

void main()
{
    vec4 frag_rect_coords = vec4(rect_coords[0], rect_coords[1], frame_size.y - rect_coords[3], frame_size.y - rect_coords[2]);
    vec2 vec_centered_pos = (gl_FragCoord.xy - vec2(frag_rect_coords[0] + frag_rect_coords[1], frag_rect_coords[2] + frag_rect_coords[3]) * 0.5);
    float start_rad = radians(start_angle);
    float end_rad = radians(end_angle);
    
    // check if the arc is a full circle (360 degrees or more)
    // the sd_rounded_arc function creates segment at the start/end angle, which is undesirable for a complete circle
    float dist;
    if (abs(end_rad - start_rad) >= 2.0 * PI - 0.001)
    {
        // full circle
        float r = length(vec_centered_pos);
        
        if (inner_radius < 0.5)
        {
            // no inner radius
            dist = r - outer_radius;
        }
        else
        {
            float ring_center_radius = (inner_radius + outer_radius) * 0.5;
            float ring_thickness = (outer_radius - inner_radius) * 0.5;
            dist = abs(r - ring_center_radius) - ring_thickness;
        }
    }
    else
    {
        dist = sd_rounded_arc(vec_centered_pos, inner_radius, outer_radius, start_rad, end_rad, corner_radius);
    }

    vec4 final_color = fill_color;

    if (stroke_width > 0.0)
    {
        // create a mask for the fill area (inside, shrunk by stroke width)
        float fill_mask = smoothstep(edge_softness, -edge_softness, dist + stroke_width);

        // combine fill mask and colors (fill + stroke)
        final_color = mix(stroke_color, fill_color, fill_mask);
    }

    // smooth edges
    float final_alpha = smoothstep(edge_softness, -edge_softness, dist);
    
    // apply the final alpha to the combined color
    gl_FragColor = vec4(final_color.rgb, final_color.a * final_alpha);
}
