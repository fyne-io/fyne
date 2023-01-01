#version 110

/* scaled params */
uniform vec4 frame_size;  //width = x, height = y (z, w NOT USED)
uniform vec4 rect_coords; //x1 [0], x2 [1], y1 [2], y2 [3]; coords of the rect_frame
uniform float stroke_width_half;
uniform float radius;
/* colors params*/
uniform vec4 fill_color;
uniform vec4 stroke_color;



float calc_distance(vec2 p, vec2 b, float r)
{
    vec2 d = abs(p) - b + vec2(r);
	return min(max(d.x, d.y), 0.0) + length(max(d, 0.0)) - r;   
}

void main() {

    vec4 frag_rect_coords = vec4(rect_coords[0], rect_coords[1], frame_size.y - rect_coords[3], frame_size.y - rect_coords[2]);
    vec2 rect_size_half =  vec2(frag_rect_coords[1] - frag_rect_coords[0], frag_rect_coords[3] - frag_rect_coords[2]) * 0.5 - vec2(stroke_width_half);
    vec2 vec_centered_pos = (gl_FragCoord.xy - vec2(frag_rect_coords[0] + frag_rect_coords[1], frag_rect_coords[2] + frag_rect_coords[3]) * 0.5);

    float distance = calc_distance(vec_centered_pos, rect_size_half, radius - stroke_width_half);

    vec4 from_color = stroke_color; //Always the border color. If no border, this still should be set
    vec4 to_color = vec4(0.5, 0.5, 0.5, 0.5); //Outside color
    if ( frame_size.x == gl_FragCoord.x ){
        to_color = vec4(1.0, 0.0, 0.0, 0.5); //Outside color
    } else if ( frame_size.y == gl_FragCoord.y ){
        to_color = vec4(1.0, 0.0, 0.0, 0.5); //Outside color
    } 

    if (stroke_width_half > 0.0)
    {
        if (distance < 0.0)
        {
            to_color = fill_color;   
        } 
        
        distance = abs(distance) - stroke_width_half;
    }

    float blend_amount = smoothstep(-1.0, 1.0, distance);

    // final color
    gl_FragColor = mix(from_color, to_color, blend_amount);
    //gl_FragColor = vec4(vec3(blend_amount), 1.0);
    //gl_FragColor = vec4(vec3(abs(distance) / (2.0 * corner)), 1.0);

}
