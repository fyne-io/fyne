#version 110

uniform float alpha;

attribute vec3 vert;
attribute vec2 vertTexCoord;

varying vec2 fragTexCoord;
varying float fragAlpha;

void main() {
    fragTexCoord = vertTexCoord;
    fragAlpha = alpha;

    gl_Position = vec4(vert, 1);
}