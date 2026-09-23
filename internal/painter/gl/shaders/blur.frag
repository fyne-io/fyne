#version 110

uniform sampler2D tex;
uniform sampler2D kernelTex;

varying vec2 fragTexCoord;
uniform float radius;
uniform vec2 direction;
uniform float sampleScale;
uniform float cornerRadius;
uniform vec2 size;

float getKernel(int i, int kernelLen) {
    float u = (float(i) + 0.5) / float(kernelLen);
    return texture2D(kernelTex, vec2(u, 0.5)).r;
}

void main() {
    float alpha = 1.0;
    if (cornerRadius > 0.5) {
        vec2 pos = fragTexCoord * size;
        vec2 halfSize = size * 0.5;
        vec2 q = abs(pos - halfSize) - halfSize + cornerRadius;
        float dist = min(max(q.x, q.y), 0.0) + length(max(q, 0.0)) - cornerRadius;
        alpha = 1.0 - smoothstep(-0.5, 0.5, dist);
    }

    int length = 2 * int(radius) + 1;
    vec4 sum = vec4(0.0);

    for (int i = 0; i < length; ++i) {
        float offset = (float(i) - radius) * sampleScale;
        vec2 tc = fragTexCoord + direction * offset;
        sum += getKernel(i, length) * texture2D(tex, tc);
    }

    gl_FragColor = mix(texture2D(tex, fragTexCoord), sum, alpha);
}
