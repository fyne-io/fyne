#version 110

uniform vec2 frame_size;
uniform vec4 rect_coords;
uniform float edge_softness;
uniform float stroke_width_half;
uniform vec4 stroke_color;

uniform vec2 start_point;
uniform vec2 end_point;
uniform vec2 control_point1; // used for quadratic and cubic
uniform vec2 control_point2; // used for cubic only
uniform float num_control_points; // 0: segment, 1: quadratic, 2: cubic

const float EPS = 1e-3;

float segment_distance(vec2 p, vec2 a, vec2 b)
{
	vec2 pa = p - a, ba = b - a;
    return length(pa - ba * clamp(dot(pa, ba) / dot(ba, ba), 0.0, 1.0));
}

// Inigo Quilez sdf functions: https://iquilezles.org/articles/distfunctions2d/ (The MIT License)
float dot2( vec2 v ) { return dot(v,v); }
float cro( vec2 a, vec2 b ) { return a.x*b.y-a.y*b.x; }
float cos_acos_3( float x ) { x=sqrt(0.5+0.5*x); return x*(x*(x*(x*-0.008972+0.039071)-0.107074)+0.576975)+0.5; }
float quadratic_distance(vec2 pos, vec2 v0, vec2 v1, vec2 v2)
{
    vec2 a = v1 - v0;
    vec2 b = v0 - 2.0*v1 + v2;
    vec2 c = a * 2.0;
    vec2 d = v0 - pos;

    float kk = 1.0 / dot(b, b);
    float kx = kk * dot(a, b);
    float ky = kk * (2.0 * dot(a, a) + dot(d, b)) / 3.0;
    float kz = kk * dot(d, a);

    float p = ky - kx*kx;
    float q = kx * (2.0*kx*kx - 3.0*ky) + kz;
    float h = q*q + 4.0*p*p*p;

    float res, sgn;

    if (h >= 0.0)
    {
        h = sqrt(h);
        h = (q < 0.0) ? h : -h; // copysign
        float x = (h - q) / 2.0;
        float v = sign(x) * pow(abs(x), 1.0 / 3.0);
        float t = v - p/(v + EPS);

		// Newton refinement (stabilizes near cancellation)
        t -= (t * (t*t + 3.0*p) + q) / (3.0*t*t + 3.0*p);

        t = clamp(t - kx, 0.0, 1.0);
        vec2 w = d + (c + b*t) * t;
        res = dot2(w);
        sgn = cro(c + 2.0*b*t, w);
    }
    else
    {
        float z = sqrt(-p);
        float m = cos_acos_3(q / (p*z*2.0));
        float n = sqrt(1.0 - m*m) * 1.732050808;

        vec3 t = clamp(vec3(m + m, -n - m, n - m) * z - kx, 0.0, 1.0);
        vec2 qx = d + (c + b*t.x) * t.x;
        float dx = dot2(qx), sx = cro(a + b*t.x, qx);
        vec2 qy = d + (c + b*t.y) * t.y;
        float dy = dot2(qy), sy = cro(a + b*t.y, qy);

        if (dx < dy)
        {
            res = dx;
            sgn = sx;
        }
        else
        {
            res = dy;
            sgn = sy;
        }
    }

    return sqrt(res) * sign(sgn);
}

float cubic_distance(vec2 p, vec2 v0, vec2 v1, vec2 v2, vec2 v3)
{
    const int STEPS = 65; // number of polyline segments for initial guess
    const int NEWTON_ITERS = 4;

    vec2 a = -v0 + 3.0*v1 - 3.0*v2 + v3;
    vec2 b =  3.0*v0 - 6.0*v1 + 3.0*v2;
    vec2 c = -3.0*v0 + 3.0*v1;
    vec2 d =  v0;

    float dt = 1.0 / float(STEPS);
    vec2 prev = d;
    float best_seg_dist = 1e3;
    float best_sgn = 0.0;
    float t0 = 0.0;

    for (int i = 1; i <= STEPS; i++)
    {
        float t = float(i) * dt;
        vec2 cur = ((a*t + b)*t + c)*t + d;

        // closest point on segment [prev, cur]
        vec2 seg = cur - prev;
        float h = clamp(dot(p - prev, seg) / max(dot(seg, seg), EPS), 0.0, 1.0);
        vec2 proj = prev + h * seg;
        float dseg = length(p - proj);

        if (dseg < best_seg_dist)
        {
            best_seg_dist = dseg;
            // map segment-local h in [0,1] to global t in [0,1]
            t0 = mix(float(i - 1) * dt, float(i) * dt, h);
            // compute sign determinant (matching quadratic convention: cro(tangent, proj - p))
            best_sgn = cro(seg, proj - p);
        }

        prev = cur;
    }

    // Newton refinement from t0 on the true cubic
    float t = t0;
    for (int iter = 0; iter < NEWTON_ITERS; iter++)
    {
        vec2 dpt = ((3.0*a*t + 2.0*b)*t) + c;
        vec2 r = ((a*t + b)*t + c)*t + d - p;
        float f  = dot(r, dpt);
        float df = dot(dpt, dpt) + dot(r, 6.0*a*t + 2.0*b);

        float step = f / max(df, EPS);
        // limit step to avoid big jumps; clamp t to [0,1]
        t = clamp(t - clamp(step, -0.5, 0.5), 0.0, 1.0);
    }

    vec2 pref = ((a*t + b)*t + c)*t + d;
    float dnewton = length(pref - p);

    // return the signed distance corresponding to the smallest unsigned distance
    return (dnewton < best_seg_dist) ? dnewton * sign(cro(((3.0*a*t + 2.0*b)*t) + c, pref - p)) : best_seg_dist * sign(best_sgn);
}

void main() {
    // coordinates: (0.0) at rect top-left, +X right, +Y down
    vec2 p = vec2(gl_FragCoord.x, frame_size.y - gl_FragCoord.y) - rect_coords.xz;

    float dist;
    if (int(num_control_points) == 1) {
        dist = quadratic_distance(p, start_point, control_point1, end_point);
    } else if (int(num_control_points) == 2) {
        dist = cubic_distance(p, start_point, control_point1, control_point2, end_point);
    } else {
        dist = segment_distance(p, start_point, end_point);
    }

    float alpha = 1.0 - smoothstep(stroke_width_half - edge_softness, stroke_width_half + edge_softness, abs(dist));
    gl_FragColor = vec4(stroke_color.rgb, stroke_color.a * alpha);
}
