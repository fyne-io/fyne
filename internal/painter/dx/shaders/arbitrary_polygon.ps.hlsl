// Signed distance to an arbitrary polygon with per-corner rounding, based on
// Inigo Quilez's sdPolygon (MIT licensed) plus the per-radius rounding technique.
//
// Every array element in a constant buffer occupies its own 16-byte register,
// so parallel vertices[] and cornerRadii[] arrays would burn 512 bytes to carry
// 12 useful ones. Both are packed into a single float4 array instead: xy is the
// vertex, z is its corner radius.

#define MAX_VERTICES 16

cbuffer PolygonConstants : register(b1) {
    float4 polyVerts[MAX_VERTICES]; // xy: vertex, z: corner radius
    float4 polyInfo;                // x: vertex count
};

static const float INF = 1e10;
static const float POLY_EPS = 1e-3;

float arbitrary_polygon_distance(float2 p, int num)
{
    if (num < 3) return 1.0;

    // Phase 1: rounded corner geometry. For each corner compute the inscribed arc
    // centre and the tangent points on the adjacent edges, which bound the
    // straight segments between corners.
    float2 startPts[MAX_VERTICES];
    float2 endPts[MAX_VERTICES];
    float arcDist = INF;
    float signVal = 1.0;

    for (int k = 0; k < MAX_VERTICES; k++)
    {
        if (k >= num) break;

        int i = k - 2;
        if (i < 0) i += num;
        int j = k - 1;
        if (j < 0) j += num;

        float2 point1 = polyVerts[i].xy;
        float2 point2 = polyVerts[j].xy;
        float2 point3 = polyVerts[k].xy;
        float radiusJ = polyVerts[j].z;

        float2 pos = p - point2;
        float2 a = normalize(point1 - point2);
        float2 b = normalize(point3 - point2);
        float crossAB = abs(a.x * b.y - a.y * b.x);

        // degenerate case: the edges are parallel
        if (crossAB < POLY_EPS)
        {
            startPts[k] = point2;
            endPts[k] = point2;
            continue;
        }

        float2 centre = (a + b) * radiusJ / crossAB;
        float2 posShifted = pos - centre;

        // even-odd rule for the arc region
        float c = radiusJ * radiusJ - posShifted.y * posShifted.y;
        if (c > 0.0)
        {
            c = sqrt(c);
            float2 p1 = float2(-c, posShifted.y);
            float2 p2Arc = float2(c, posShifted.y);
            float s1 = dot(p1, a);
            float s2 = dot(p1, b);
            float s3 = dot(p2Arc, a);
            float s4 = dot(p2Arc, b);
            if (posShifted.x < p1.x && s1 < 0.0 && s2 < 0.0)
            {
                signVal = -signVal;
            }
            if (posShifted.x < p2Arc.x && s4 < 0.0 && s3 < 0.0)
            {
                signVal = -signVal;
            }
        }

        // distance to the arc at this corner
        float s1Dot = dot(posShifted, a);
        float s2Dot = dot(posShifted, b);
        if (s1Dot < 0.0 && s2Dot < 0.0)
        {
            arcDist = min(arcDist, abs(length(posShifted) - radiusJ));
        }

        startPts[k] = point2 + a * dot(centre, a);
        endPts[k] = point2 + b * dot(centre, b);
    }

    // Phase 2: distance to the straight segments between tangent points.
    float edgeDist = INF;
    for (int j2 = 0; j2 < MAX_VERTICES; j2++)
    {
        if (j2 >= num) break;

        int i2 = j2 - 1;
        if (i2 < 0) i2 += num;

        float2 start = endPts[i2];
        float2 endSeg = startPts[j2];
        float2 e = endSeg - start;
        float2 w = p - start;
        float h = clamp(dot(w, e) / dot(e, e), 0.0, 1.0);
        float2 d = w - e * h;
        edgeDist = min(edgeDist, dot(d, d));

        // even-odd rule for the edge segments
        if ((w.y > 0.0) != (w.y > e.y))
        {
            if ((e.y * w.x < e.x * w.y) != (e.y < 0.0))
            {
                signVal = -signVal;
            }
        }
    }

    return min(arcDist, sqrt(edgeDist)) * signVal;
}

float4 main(PSIn input) : SV_TARGET
{
    // (0,0) at the rect top-left, +X right, +Y down: SV_Position and bounds
    // are both top-origin, so subtracting maps into the rect directly.
    float2 p = input.pos.xy - bounds.xy;

    float edgeSoftness = rectHalf.w;
    float strokeWidth = misc.x;

    float dist = arbitrary_polygon_distance(p, int(polyInfo.x));

    float4 finalColor = fillColor;
    if (strokeWidth > 0.0)
    {
        float fillMask = smoothstep(-strokeWidth + edgeSoftness, -strokeWidth - edgeSoftness, dist);
        finalColor = lerp(strokeColor, fillColor, fillMask);
    }

    float finalAlpha = smoothstep(edgeSoftness, -edgeSoftness, dist);
    return float4(finalColor.rgb, finalColor.a * finalAlpha);
}
