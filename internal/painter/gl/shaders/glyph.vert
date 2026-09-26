#version 110

attribute vec2 vert;
attribute vec2 vertTexCoord;
attribute vec4 vertColor;

varying vec2 fragTexCoord;
varying vec4 fragColor;

void main() {
    fragTexCoord = vertTexCoord;
    fragColor = vertColor;

    gl_Position = vec4(vert, 0, 1);
}
