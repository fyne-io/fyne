// Distance-field stroke of a linear, quadratic or cubic Bezier. The curve points
// ride in their own register(b1) buffer; rectHalf.z is the half stroke width and
// texParams.x selects the degree (0 linear, 1 quadratic, 2 cubic).
//
// The distance functions are adapted from Inigo Quilez's work,
// https://www.shadertoy.com/view/MlKcDD, MIT licensed.

cbuffer CurveConstants : register(b1) {
    float4 curveStartEnd;  // xy: start point, zw: end point
    float4 curveControls;  // xy: control point 1, zw: control point 2
};

static const float EPS = 1e-3;

float linear_distance(float2 p, float2 a, float2 b)
{
    float2 pa = p - a, ba = b - a;
    return length(pa - ba * clamp(dot(pa, ba) / dot(ba, ba), 0.0, 1.0));
}

float dot2(float2 v) { return dot(v, v); }

float cos_acos_3(float x)
{
    x = sqrt(0.5 + 0.5 * x);
    return x * (x * (x * (x * -0.008972 + 0.039071) - 0.107074) + 0.576975) + 0.5;
}

float quadratic_distance(float2 pos, float2 v0, float2 v1, float2 v2)
{
    float2 a = v1 - v0;
    float2 b = v0 - 2.0 * v1 + v2;
    float2 c = a * 2.0;
    float2 d = v0 - pos;

    float kk = 1.0 / dot(b, b);
    float kx = kk * dot(a, b);
    float ky = kk * (2.0 * dot(a, a) + dot(d, b)) / 3.0;
    float kz = kk * dot(d, a);

    float p = ky - kx * kx;
    float q = kx * (2.0 * kx * kx - 3.0 * ky) + kz;
    float h = q * q + 4.0 * p * p * p;

    float res;
    if (h >= 0.0)
    {
        h = sqrt(h);
        h = (q < 0.0) ? h : -h; // copysign
        float x = (h - q) / 2.0;
        float v = sign(x) * pow(abs(x), 1.0 / 3.0);
        float t = v - p / (v + EPS);

        // Newton refinement, stabilises near cancellation
        t -= (t * (t * t + 3.0 * p) + q) / (3.0 * t * t + 3.0 * p);

        t = clamp(t - kx, 0.0, 1.0);
        float2 w = d + (c + b * t) * t;
        res = dot2(w);
    }
    else
    {
        float z = sqrt(-p);
        float m = cos_acos_3(q / (p * z * 2.0));
        float n = sqrt(1.0 - m * m) * 1.732050808;

        float3 t = clamp(float3(m + m, -n - m, n - m) * z - kx, 0.0, 1.0);
        float2 qx = d + (c + b * t.x) * t.x;
        float2 qy = d + (c + b * t.y) * t.y;
        res = min(dot2(qx), dot2(qy));
    }
    return sqrt(res);
}

float cubic_distance(float2 p, float2 v0, float2 v1, float2 v2, float2 v3)
{
    const int NEWTON_ITERS = 4;

    float approxLength = length(v1 - v0) + length(v2 - v1) + length(v3 - v2);
    int steps = int(clamp(approxLength * 10.0, 16.0, 64.0));

    float2 a = -v0 + 3.0 * v1 - 3.0 * v2 + v3;
    float2 b = 3.0 * v0 - 6.0 * v1 + 3.0 * v2;
    float2 c = -3.0 * v0 + 3.0 * v1;
    float2 d = v0;

    float dt = 1.0 / float(steps);
    float2 prev = d;
    float bestSegDist = 1e3;
    float t0 = 0.0;

    for (int i = 1; i <= steps; i++)
    {
        float t = float(i) * dt;
        float2 cur = ((a * t + b) * t + c) * t + d;

        // closest point on segment [prev, cur]
        float2 seg = cur - prev;
        float h = clamp(dot(p - prev, seg) / max(dot(seg, seg), EPS), 0.0, 1.0);
        float2 proj = prev + h * seg;
        float dseg = length(p - proj);

        if (dseg < bestSegDist)
        {
            bestSegDist = dseg;
            // map segment-local h in [0,1] to global t in [0,1]
            t0 = lerp(float(i - 1) * dt, float(i) * dt, h);
        }
        prev = cur;
    }

    // Newton refinement from t0 on the true cubic
    float t = t0;
    for (int iter = 0; iter < NEWTON_ITERS; iter++)
    {
        float2 dpt = ((3.0 * a * t + 2.0 * b) * t) + c;
        float2 r = ((a * t + b) * t + c) * t + d - p;
        float f = dot(r, dpt);
        float df = dot(dpt, dpt) + dot(r, 6.0 * a * t + 2.0 * b);

        float st = f / max(df, EPS);
        // limit the step to avoid big jumps, and keep t in [0,1]
        t = clamp(t - clamp(st, -0.5, 0.5), 0.0, 1.0);
    }

    float2 pref = ((a * t + b) * t + c) * t + d;
    return min(length(pref - p), bestSegDist);
}

float4 main(PSIn input) : SV_TARGET
{
    // (0,0) at the rect top-left, +X right, +Y down: SV_Position and bounds
    // are both top-origin, so subtracting maps into the rect directly.
    float2 p = input.pos.xy - bounds.xy;

    float strokeWidthHalf = rectHalf.z;
    float edgeSoftness = rectHalf.w;
    int controlPoints = int(texParams.x);

    float dist;
    if (controlPoints == 1) {
        dist = quadratic_distance(p, curveStartEnd.xy, curveControls.xy, curveStartEnd.zw);
    } else if (controlPoints == 2) {
        dist = cubic_distance(p, curveStartEnd.xy, curveControls.xy, curveControls.zw, curveStartEnd.zw);
    } else {
        dist = linear_distance(p, curveStartEnd.xy, curveStartEnd.zw);
    }

    float alpha = 1.0 - smoothstep(strokeWidthHalf - edgeSoftness, strokeWidthHalf + edgeSoftness, dist);
    return float4(strokeColor.rgb, strokeColor.a * alpha);
}
