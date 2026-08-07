// Signed-distance rounded rectangle with per-corner radii, stroke and shadow.
float calc_distance(float2 p, float2 b, float4 r)
{
    r.xy = (p.x > 0.0) ? r.xy : r.zw;
    r.x  = (p.y > 0.0) ? r.x  : r.y;

    float2 d = abs(p) - b + r.x;
    return min(max(d.x, d.y), 0.0) + length(max(d, 0.0)) - r.x;
}

// distance across all quadrants; correct when a radius exceeds the short edge
float calc_distance_all_quadrants(float2 p, float2 size, float4 rad)
{
    float2 d = abs(p) - size;
    float dist = length(max(d, 0.0)) + min(max(d.x, d.y), 0.0);

    float2 pTL = p - float2(rad.z - size.x, size.y - rad.z);
    if (pTL.x < 0.0 && pTL.y > 0.0) dist = max(dist, length(pTL) - rad.z);

    float2 pTR = p - float2(size.x - rad.x, size.y - rad.x);
    if (pTR.x > 0.0 && pTR.y > 0.0) dist = max(dist, length(pTR) - rad.x);

    float2 pBR = p - float2(size.x - rad.y, rad.y - size.y);
    if (pBR.x > 0.0 && pBR.y < 0.0) dist = max(dist, length(pBR) - rad.y);

    float2 pBL = p - float2(rad.w - size.x, rad.w - size.y);
    if (pBL.x < 0.0 && pBL.y < 0.0) dist = max(dist, length(pBL) - rad.w);

    return dist;
}

float4 main(PSIn input) : SV_TARGET
{
    float2 p = centredUp(input.pos.xy);

    float strokeWidthHalf = rectHalf.z;
    float edgeSoftness = rectHalf.w;
    float2 rectSizeHalf = rectHalf.xy;
    float addShadow = misc.y;
    float shadowBlurRadius = misc.z;
    float shadowSpread = misc.w;
    float shadowType = shadowOffset.z;

    float dist;
    float maxRadius = max(max(radius.x, radius.y), max(radius.z, radius.w));
    float4 finalColor = fillColor;
    float finalAlpha;

    bool allQuadrants = maxRadius - 0.9 > min(rectSizeHalf.x, rectSizeHalf.y) + strokeWidthHalf;
    if (allQuadrants)
    {
        dist = calc_distance_all_quadrants(p, rectSizeHalf + strokeWidthHalf, radius);
        finalAlpha = 1.0 - smoothstep(-edgeSoftness, edgeSoftness, dist);

        if (strokeWidthHalf > 0.0)
        {
            float colorBlend = 1.0 - smoothstep(strokeWidthHalf * 2.0 - edgeSoftness,
                strokeWidthHalf * 2.0 + edgeSoftness, abs(dist));
            finalColor = lerp(fillColor, strokeColor, colorBlend);
        }
    }
    else
    {
        dist = calc_distance(p, rectSizeHalf, radius - strokeWidthHalf);
        finalAlpha = 1.0 - smoothstep(strokeWidthHalf - edgeSoftness, strokeWidthHalf + edgeSoftness, dist);

        if (strokeWidthHalf > 0.0)
        {
            float colorBlend = smoothstep(-strokeWidthHalf - edgeSoftness,
                -strokeWidthHalf + edgeSoftness, dist);
            finalColor = lerp(fillColor, strokeColor, colorBlend);
        }
    }

    finalColor = float4(finalColor.rgb, finalColor.a * finalAlpha);

    if (addShadow == 1.0)
    {
        float2 shadowSize = rectSizeHalf + strokeWidthHalf;
        float4 shadowRadius = radius;

        if (shadowSpread != 0.0)
        {
            float2 originalSize = shadowSize;
            shadowSize = max(originalSize + shadowSpread, 0.0);
            float ratioX = (originalSize.x > 0.0) ? (shadowSize.x / originalSize.x) : 1.0;
            float ratioY = (originalSize.y > 0.0) ? (shadowSize.y / originalSize.y) : 1.0;
            shadowRadius = max(radius * min(ratioX, ratioY), 0.0);
        }

        float blurInset = shadowBlurRadius * 0.5;
        shadowSize = max(shadowSize - blurInset, 0.0);
        shadowRadius = max(shadowRadius - blurInset, 0.0);

        float2 shadowOffsetCorrected = float2(-shadowOffset.x, shadowOffset.y);
        float distShadow;
        if (allQuadrants) {
            distShadow = calc_distance_all_quadrants(p + shadowOffsetCorrected, shadowSize, shadowRadius);
        } else {
            distShadow = calc_distance(p + shadowOffsetCorrected, shadowSize, shadowRadius);
        }
        float shadowAlpha = shadowColor.a *
            (1.0 - smoothstep(-edgeSoftness, shadowBlurRadius + edgeSoftness, distShadow));

        if (shadowType == 0.0)
        {
            float dShape;
            if (allQuadrants) {
                dShape = dist;
            } else {
                dShape = calc_distance(p, rectSizeHalf + strokeWidthHalf, radius);
            }
            shadowAlpha *= smoothstep(-2.0 * edgeSoftness, 0.0, dShape);
        }

        finalColor = blend_shadow(finalColor, float4(shadowColor.rgb, shadowAlpha));
    }

    return finalColor;
}
