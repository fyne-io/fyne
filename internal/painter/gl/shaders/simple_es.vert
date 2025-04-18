#version 100

#ifdef GL_ES
# ifdef GL_FRAGMENT_PRECISION_HIGH
precision highp float;
# else
precision mediump float;
#endif
precision mediump int;
precision lowp sampler2D;
#endif

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