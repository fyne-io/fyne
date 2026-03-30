#version 110

uniform sampler2D tex;

varying vec2 fragTexCoord;
uniform float radius;
uniform vec2 size;

// 50 on either side, scaled up to 4x pixel density plus 1
uniform float kernel[401];

void main() {
    vec2 inverseSize = vec2(1.0/size.x, 1.0/size.y);
    int length = 2 * int(radius) + 1;
    vec4 sum = vec4(0.0);

	for (int i = 0; i < length; ++i)
	{
		for (int j = 0; j < length; ++j)
		{
			vec2 tc = fragTexCoord + inverseSize * vec2(float(i) - radius, float(j) - radius);
			sum += kernel[i] * kernel[j] * texture2D(tex, tc);
		}
	}
    gl_FragColor = sum;
}
