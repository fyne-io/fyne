#version 110

// Lines are batched, so colour, width and feather arrive per vertex; they are
// the same across a line's vertices, so the varyings carry them unchanged.
varying vec2 delta;
varying vec4 color;
varying vec2 style;

void main() {
    float lineWidth = style.x;
    float feather = style.y;
    float distance = length(delta);

    if (feather == 0.0 || distance <= lineWidth - feather) {
        gl_FragColor = color;
    } else {
        gl_FragColor = vec4(color.r, color.g, color.b, mix(color.a, 0.0, (distance - (lineWidth - feather)) / feather));
    }
}
