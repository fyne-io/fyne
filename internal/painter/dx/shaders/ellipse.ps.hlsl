// Ellipse fill and stroke by signed distance, also used for circles.
//
// radius.xy holds the two semi-axes and shadowOffset.w the rotation in degrees.
float2 rotate(float2 v, float a)
{
    float s = sin(-a);
    float c = cos(-a);
    return float2(c * v.x - s * v.y, s * v.x + c * v.y);
}

float calc_distance(float2 p, float2 r)
{
    const float eps = 1e-3;
    r = max(r, eps);
    float2 f = p / r;
    return (dot(f, f) - 1.0) / max(length(2.0 * f / r), eps);
}

float4 main(PSIn input) : SV_TARGET
{
    float2 p = centredUp(input.pos.xy);

    float strokeWidth = misc.x;
    float addShadow = misc.y;
    float shadowBlurRadius = misc.z;
    float shadowSpread = misc.w;
    float shadowType = shadowOffset.z;
    float edgeSoftness = rectHalf.w;
    float2 rad = radius.xy;

    p = rotate(p, radians(shadowOffset.w));

    float dist = calc_distance(p, rad);
    float4 finalColor = fillColor;

    if (strokeWidth > 0.0)
    {
        float2 innerRadius = rad - strokeWidth;
        float fillMask = 0.0;
        if (innerRadius.x > 1.0 && innerRadius.y > 1.0)
        {
            float distInner = calc_distance(p, innerRadius);
            fillMask = smoothstep(edgeSoftness, -edgeSoftness, distInner);
        }
        finalColor = lerp(strokeColor, fillColor, fillMask);
    }

    float finalAlpha = smoothstep(edgeSoftness, -edgeSoftness, dist);
    finalColor = float4(finalColor.rgb, finalColor.a * finalAlpha);

    if (addShadow == 1.0)
    {
        float2 shadowRadius = rad;
        if (shadowSpread != 0.0) {
            shadowRadius = max(rad + shadowSpread, 0.0);
        }

        float blurInset = shadowBlurRadius * 0.5;
        shadowRadius = max(shadowRadius - blurInset, 0.0);

        float2 shadowOffsetCorrected = float2(-shadowOffset.x, shadowOffset.y);
        float distShadow = calc_distance(p + shadowOffsetCorrected, shadowRadius);
        float shadowAlpha = shadowColor.a *
            (1.0 - smoothstep(-edgeSoftness, shadowBlurRadius + edgeSoftness, distShadow));

        if (shadowType == 0.0) {
            shadowAlpha *= smoothstep(-2.0 * edgeSoftness, 0.0, dist);
        }

        finalColor = blend_shadow(finalColor, float4(shadowColor.rgb, shadowAlpha));
    }

    return finalColor;
}
