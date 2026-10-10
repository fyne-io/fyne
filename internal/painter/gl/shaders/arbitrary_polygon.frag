#version 110

#define MAX_VERTICES 16

uniform vec2 frame;
uniform vec4 bounds;
uniform float edgeSoftness;

uniform vec2 vertices[MAX_VERTICES];
uniform float cornerRadii[MAX_VERTICES];
uniform float vertexCount;

uniform vec4 fillColor;
uniform float strokeWidth;
uniform vec4 strokeColor;

const float INF = 1e10;
const float EPS = 1e-3;

// Signed distance to an arbitrary polygon with per-corner rounding.
// Based on Inigo Quilez's sdPolygon (MIT License) and the per-radius rounding technique.
float arbitrary_polygon_distance(vec2 p, int num)
{
    if (num < 3) return 1.0;

    // Phase 1: Compute rounded corner geometry
    // For each corner, compute the inscribed arc center and the tangent points
    // on the adjacent edges. These define the straight segments between corners.
    vec2 start_pts[MAX_VERTICES];
    vec2 end_pts[MAX_VERTICES];
    float arc_dist = INF;
    float sign_val = 1.0;

    for (int k = 0; k < MAX_VERTICES; k++)
    {
        if (k >= num) break;

        int i = k - 2;
        if (i < 0) i += num;
        int j = k - 1;
        if (j < 0) j += num;

        vec2 point1 = vertices[i];
        vec2 point2 = vertices[j];
        vec2 point3 = vertices[k];
        float radius = cornerRadii[j];

        vec2 pos = p - point2;
        vec2 a = normalize(point1 - point2);
        vec2 b = normalize(point3 - point2);
        float cross_ab = abs(a.x * b.y - a.y * b.x);

        // Avoid degenerate case where edges are parallel
        if (cross_ab < EPS)
        {
            start_pts[k] = point2;
            end_pts[k] = point2;
            continue;
        }

        vec2 center = (a + b) * radius / cross_ab;
        vec2 pos_shifted = pos - center;

        // Even-odd rule for arc region
        float c = radius * radius - pos_shifted.y * pos_shifted.y;
        if (c > 0.0)
        {
            c = sqrt(c);
            vec2 p1 = vec2(-c, pos_shifted.y);
            vec2 p2_arc = vec2(c, pos_shifted.y);
            float s1 = dot(p1, a);
            float s2 = dot(p1, b);
            float s3 = dot(p2_arc, a);
            float s4 = dot(p2_arc, b);
            if (pos_shifted.x < p1.x && s1 < 0.0 && s2 < 0.0)
            {
                sign_val = -sign_val;
            }
            if (pos_shifted.x < p2_arc.x && s4 < 0.0 && s3 < 0.0)
            {
                sign_val = -sign_val;
            }
        }

        // Distance to arc at this corner
        float s1_dot = dot(pos_shifted, a);
        float s2_dot = dot(pos_shifted, b);
        if (s1_dot < 0.0 && s2_dot < 0.0)
        {
            float d = abs(length(pos_shifted) - radius);
            arc_dist = min(arc_dist, d);
        }

        // Tangent points on the edges
        vec2 start = point2 + a * dot(center, a);
        vec2 end_pt = point2 + b * dot(center, b);
        start_pts[k] = start;
        end_pts[k] = end_pt;
    }

    // Phase 2: Distance to straight edge segments between tangent points
    float edge_dist = INF;
    for (int j2 = 0; j2 < MAX_VERTICES; j2++)
    {
        if (j2 >= num) break;

        int i2 = j2 - 1;
        if (i2 < 0) i2 += num;

        vec2 start = end_pts[i2];
        vec2 end_seg = start_pts[j2];
        vec2 e = end_seg - start;
        vec2 w = p - start;
        float h = clamp(dot(w, e) / dot(e, e), 0.0, 1.0);
        vec2 d = w - e * h;
        edge_dist = min(edge_dist, dot(d, d));

        // Even-odd rule for edge segments
        if ((w.y > 0.0) != (w.y > e.y))
        {
            if ((e.y * w.x < e.x * w.y) != (e.y < 0.0))
            {
                sign_val = -sign_val;
            }
        }
    }

    return min(arc_dist, sqrt(edge_dist)) * sign_val;
}

void main()
{
    // coordinates: (0.0) at rect top-left, +X right, +Y down
    vec2 p = vec2(gl_FragCoord.x, frame.y - gl_FragCoord.y) - bounds.xy;

    int num = int(vertexCount);
    float dist = arbitrary_polygon_distance(p, num);
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
