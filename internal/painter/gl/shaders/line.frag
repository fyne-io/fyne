#version 110

uniform vec4 color;
uniform float lineWidth;
uniform float feather;

varying vec2 delta;

void main() {
    float alpha = color.a;
    float distance = length(delta);

    if (feather == 0.0 || distance <= lineWidth - feather) {
        gl_FragColor = color;
    } else {
        gl_FragColor = vec4(color.r, color.g, color.b, mix(color.a, 0.0, (distance - (lineWidth - feather)) / feather));
    }
}
