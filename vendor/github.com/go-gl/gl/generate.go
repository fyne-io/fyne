// +build gen

//go:generate glow generate -out=./v2.1/gl/ -api=gl -version=2.1 -xml=../glow/xml/
//go:generate glow generate -out=./all-core/gl/ -api=gl -version=all -profile=core -lenientInit -xml=../glow/xml/
//go:generate glow generate -out=./v3.2-core/gl/ -api=gl -version=3.2 -profile=core -xml=../glow/xml/
//go:generate glow generate -out=./v3.3-core/gl/ -api=gl -version=3.3 -profile=core -xml=../glow/xml/
//go:generate glow generate -out=./v4.1-core/gl/ -api=gl -version=4.1 -profile=core -xml=../glow/xml/
//go:generate glow generate -out=./v4.2-core/gl/ -api=gl -version=4.2 -profile=core -xml=../glow/xml/
//go:generate glow generate -out=./v4.3-core/gl/ -api=gl -version=4.3 -profile=core -xml=../glow/xml/
//go:generate glow generate -out=./v4.4-core/gl/ -api=gl -version=4.4 -profile=core -xml=../glow/xml/
//go:generate glow generate -out=./v4.5-core/gl/ -api=gl -version=4.5 -profile=core -xml=../glow/xml/
//go:generate glow generate -out=./v4.6-core/gl/ -api=gl -version=4.6 -profile=core -xml=../glow/xml/
//go:generate glow generate -out=./v3.2-compatibility/gl/ -api=gl -version=3.2 -profile=compatibility -xml=../glow/xml/
//go:generate glow generate -out=./v3.3-compatibility/gl/ -api=gl -version=3.3 -profile=compatibility -xml=../glow/xml/
//go:generate glow generate -out=./v4.1-compatibility/gl/ -api=gl -version=4.1 -profile=compatibility -xml=../glow/xml/
//go:generate glow generate -out=./v4.2-compatibility/gl/ -api=gl -version=4.2 -profile=compatibility -xml=../glow/xml/
//go:generate glow generate -out=./v4.3-compatibility/gl/ -api=gl -version=4.3 -profile=compatibility -xml=../glow/xml/
//go:generate glow generate -out=./v4.4-compatibility/gl/ -api=gl -version=4.4 -profile=compatibility -xml=../glow/xml/
//go:generate glow generate -out=./v4.5-compatibility/gl/ -api=gl -version=4.5 -profile=compatibility -xml=../glow/xml/
//go:generate glow generate -out=./v4.6-compatibility/gl/ -api=gl -version=4.6 -profile=compatibility -xml=../glow/xml/
//go:generate glow generate -out=./v3.1/gles2/ -api=gles2 -version=3.1 -xml=../glow/xml/

// This is an empty pseudo-package with the sole purpose of containing go generate directives
// that generate all gl binding packages inside this repository.
package gl
