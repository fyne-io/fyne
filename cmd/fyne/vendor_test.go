package main

import (
	"go/build"
	"io"
	"path/filepath"
	"strings"
	"testing"
)

func Test_pkgModPath(t *testing.T) {
	const mockModulesTxtMissing = `# github.com/go-gl/gl v0.0.0-20181026044259-55b76b7df9d2
github.com/go-gl/gl/v3.2-core/gl
# github.com/goki/freetype v0.0.0-20181231101311-fa8a33aabaff
`
	const mockModulesTxtMatch = `# github.com/go-gl/gl v0.0.0-20181026044259-55b76b7df9d2
github.com/go-gl/gl/v3.2-core/gl
# github.com/go-gl/glfw/v3.3/glfw v0.0.0-20181213070059-819e8ce5125f
github.com/go-gl/glfw/v3.3/glfw
# github.com/goki/freetype v0.0.0-20181231101311-fa8a33aabaff
`
	const mockModulesTxtOld = `# github.com/go-gl/glfw v0.0.0-20181213070059-819e8ce5125f
github.com/go-gl/glfw/v3.2/glfw
`
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    string
		target  string
		wantErr bool
	}{
		{
			name: "No match",
			args: args{
				r: strings.NewReader(mockModulesTxtMissing),
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "Match",
			args: args{
				r: strings.NewReader(mockModulesTxtMatch),
			},
			want:    filepath.Join(build.Default.GOPATH, "pkg/mod/github.com/go-gl/glfw/v3.3/glfw@v0.0.0-20181213070059-819e8ce5125f"),
			target:  "github.com/go-gl/glfw/v3.3/glfw",
			wantErr: false,
		},
		{
			name: "Match (Old)",
			args: args{
				r: strings.NewReader(mockModulesTxtOld),
			},
			want:    filepath.Join(build.Default.GOPATH, "pkg/mod/github.com/go-gl/glfw@v0.0.0-20181213070059-819e8ce5125f/v3.2/glfw"),
			target:  "github.com/go-gl/glfw/v3.2/glfw",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, target, err := cacheModPath(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("cacheModPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("cacheModPath() = %v, want %v", got, tt.want)
			}
			if target != tt.target {
				t.Errorf("target was %v, want %v", target, tt.target)
			}
		})
	}
}
