// A regular n-gon by signed distance, with rotation, per-shape corner rounding
// and an inset stroke. radius.x is the outer radius, radius.z the corner radius,
// and shadowOffset.w the rotation in degrees.

// floorMod is a floored modulus: the result always takes the sign of y. The
// built-in fmod() truncates toward zero instead, which would break the angular
// offset below for angles left of the vertical.
float floorMod(float x, float y)
{
    return x - y * floor(x / y);
}

float2 rotate(float2 v, float a)
{
    float s = sin(-a);
    float c = cos(-a);
    return float2(c * v.x - s * v.y, s * v.x + c * v.y);
}

// Signed distance from a point to a regular polygon's edge.
float regular_distance(float2 p, float r, int s)
{
    float angle = PI / float(s);
    float angleCos = cos(angle);
    float angleSin = sin(angle);
    float angularOffset = floorMod(atan2(p.x, p.y), 2.0 * angle) - angle;
    float2 dist = length(p) * float2(cos(angularOffset), abs(sin(angularOffset)))
        - r * float2(angleCos, angleSin);
    dist.y += clamp(-dist.y, 0.0, r * angleSin);
    return length(dist) * sign(dist.x);
}

float4 main(PSIn input) : SV_TARGET
{
    float2 p = centredUp(input.pos.xy);

    float outerRadius = radius.x;
    float cornerRadius = radius.z;
    float sides = radius.w;
    float edgeSoftness = rectHalf.w;
    float strokeWidth = misc.x;

    p = rotate(p, radians(shadowOffset.w));
    float dist = regular_distance(p, outerRadius - cornerRadius, int(sides)) - cornerRadius;

    float4 finalColor = fillColor;
    if (strokeWidth > 0.0)
    {
        float fillMask = smoothstep(-strokeWidth + edgeSoftness, -strokeWidth - edgeSoftness, dist);
        finalColor = lerp(strokeColor, fillColor, fillMask);
    }

    float finalAlpha = smoothstep(edgeSoftness, -edgeSoftness, dist);
    return float4(finalColor.rgb, finalColor.a * finalAlpha);
}
