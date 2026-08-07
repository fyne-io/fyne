// A signed-distance rounded arc with an optional cutout, in the unit-circle
// convention where angles are measured from the positive X axis. radius carries
// {innerRadius, outerRadius, cornerRadius} and texParams.xy the start and end
// angle in degrees - this shader samples no texture, so those slots are free.
// Signed distance to a rounded arc: negative inside, positive outside.
float sd_rounded_arc(float2 p, float r1, float r2, float a0, float a1, float cr)
{
    // Centre the arc on the positive X axis so the span is symmetric.
    float midAngle = (a0 + a1) / 2.0;
    float arcSpan = abs(a1 - a0);

    float cs = cos(midAngle);
    float sn = sin(midAngle);
    // Rotate p by -midAngle. The rows are (cs, sn) and (-sn, cs), written out
    // to keep the orientation unambiguous.
    p = float2(cs * p.x + sn * p.y, -sn * p.x + cs * p.y);

    // Distance to a rounded box in pseudo-polar space.
    float r = length(p);
    float a = atan2(p.y, p.x);

    float2 boxHalfSize = float2(arcSpan * 0.5 * r, (r2 - r1) * 0.5);
    float2 q = float2(a * r, r - (r1 + r2) * 0.5);

    // Clamp the inner corner radius to half the smaller dimension so opposite
    // corners cannot overlap: thickness (r2-r1), and inner arc length.
    float innerCr = min(cr, 0.5 * min(r2 - r1, arcSpan * r1));
    float outerCr = cr;

    // Blend inner to outer corner radius with radial position.
    float t = smoothstep(-boxHalfSize.y, boxHalfSize.y, q.y);
    float effectiveCr = lerp(innerCr, outerCr, t);

    float2 dist = abs(q) - boxHalfSize + effectiveCr;
    return length(max(dist, 0.0)) + min(max(dist.x, dist.y), 0.0) - effectiveCr;
}

float4 main(PSIn input) : SV_TARGET
{
    float2 p = centredUp(input.pos.xy);

    float innerRadius = radius.x;
    float outerRadius = radius.y;
    float cornerRadius = radius.z;
    float startRad = radians(texParams.x);
    float endRad = radians(texParams.y);
    float edgeSoftness = rectHalf.w;
    float strokeWidth = misc.x;

    float dist;
    if (abs(endRad - startRad) >= 2.0 * PI - 0.001)
    {
        // A full circle: sd_rounded_arc would cut a seam at the start/end angle.
        float r = length(p);
        if (innerRadius < 0.5)
        {
            dist = r - outerRadius;
        }
        else
        {
            float ringCentre = (innerRadius + outerRadius) * 0.5;
            float ringThickness = (outerRadius - innerRadius) * 0.5;
            dist = abs(r - ringCentre) - ringThickness;
        }
    }
    else
    {
        dist = sd_rounded_arc(p, innerRadius, outerRadius, startRad, endRad, cornerRadius);
    }

    float4 finalColor = fillColor;
    if (strokeWidth > 0.0)
    {
        float fillMask = smoothstep(edgeSoftness, -edgeSoftness, dist + strokeWidth);
        finalColor = lerp(strokeColor, fillColor, fillMask);
    }

    float finalAlpha = smoothstep(edgeSoftness, -edgeSoftness, dist);
    return float4(finalColor.rgb, finalColor.a * finalAlpha);
}
