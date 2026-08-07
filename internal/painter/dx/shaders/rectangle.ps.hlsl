// Square corners, optional inset stroke, optional drop or inner shadow.
float4 main(PSIn input) : SV_TARGET
{
    float2 fc = input.pos.xy;
    float strokeWidth = misc.x;
    float addShadow = misc.y;
    float shadowBlurRadius = misc.z;
    float shadowSpread = misc.w;
    float shadowType = shadowOffset.z;

    float4 col = fillColor;

    if (addShadow == 1.0)
    {
        float2 fragPos = float2(fc.x - shadowOffset.x, fc.y - shadowOffset.y);
        float2 centre = float2((bounds.x + bounds.z) * 0.5, (bounds.y + bounds.w) * 0.5);
        float2 halfSize = float2(bounds.z - bounds.x, bounds.w - bounds.y) * 0.5 + shadowSpread;

        float2 d = abs(fragPos - centre) - halfSize;
        float distanceShadow = smoothstep(-shadowBlurRadius * 0.5, shadowBlurRadius * 0.5,
            length(max(d, 0.0)) + min(max(d.x, d.y), 0.0));
        float shadowAlpha = shadowColor.a * (1.0 - distanceShadow);

        if (shadowType == 0.0)
        {
            float dH = min(fc.x - bounds.x, bounds.z - fc.x);
            float dV = min(bounds.w - fc.y, fc.y - bounds.y);
            shadowAlpha *= smoothstep(0.0, -0.5, min(dH, dV));
        }

        if (fc.x > bounds.z || fc.x < bounds.x || fc.y > bounds.w || fc.y < bounds.y) {
            col.a = 0.0;
        }

        col = blend_shadow(col, float4(shadowColor.rgb, shadowAlpha));
    }

    if (fc.x < bounds.x || fc.x > bounds.z || fc.y < bounds.y || fc.y > bounds.w)
    {
        if (addShadow == 0.0) {
            discard;
        }
    }
    else
    {
        if (fc.x >= bounds.z - strokeWidth) {
            col = strokeColor;
        } else if (fc.x <= bounds.x + strokeWidth) {
            col = strokeColor;
        } else if (fc.y >= bounds.w - strokeWidth) {
            col = strokeColor;
        } else if (fc.y <= bounds.y + strokeWidth) {
            col = strokeColor;
        }
    }

    return col;
}
