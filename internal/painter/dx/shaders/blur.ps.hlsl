// One pass of a separable Gaussian blur over a snapshot of what is already on
// screen. The painter runs it twice, horizontally then vertically, re-snapshotting
// in between. The kernel arrives as a 1-pixel-tall texture at t1.
//
// texParams holds {radius, cornerRadius, size.x, size.y} and rectHalf
// {direction.x, direction.y, sampleScale, unused} - not inset, which the
// vertex stage owns as the texture-coordinate source.
//
// Every sample is a SampleLevel rather than a Sample: the loop bound is dynamic,
// and HLSL forbids gradient-derived sampling under varying flow control. The
// textures are single-mip, so explicit LOD 0 is what Sample would have chosen.

Texture2D tex : register(t0);
SamplerState samp : register(s0);
Texture2D kernelTex : register(t1);
SamplerState kernelSamp : register(s1);

float getKernel(int i, int kernelLen)
{
    float u = (float(i) + 0.5) / float(kernelLen);
    return kernelTex.SampleLevel(kernelSamp, float2(u, 0.5), 0).r;
}

float4 main(PSIn input) : SV_TARGET
{
    float radius = texParams.x;
    float cornerRadius = texParams.y;
    float2 size = texParams.zw;
    float2 direction = rectHalf.xy;
    float sampleScale = rectHalf.z;

    float alpha = 1.0;
    if (cornerRadius > 0.5) {
        float2 pos = input.uv * size;
        float2 halfSize = size * 0.5;
        float2 q = abs(pos - halfSize) - halfSize + cornerRadius;
        float dist = min(max(q.x, q.y), 0.0) + length(max(q, 0.0)) - cornerRadius;
        alpha = 1.0 - smoothstep(-0.5, 0.5, dist);
    }

    int len = 2 * int(radius) + 1;
    float4 sum = float4(0.0, 0.0, 0.0, 0.0);
    for (int i = 0; i < len; ++i) {
        float offset = (float(i) - radius) * sampleScale;
        float2 tc = input.uv + direction * offset;
        sum += getKernel(i, len) * tex.SampleLevel(samp, tc, 0);
    }

    return lerp(tex.SampleLevel(samp, input.uv, 0), sum, alpha);
}
