package main

import (
	"go/build"
	"io"
	"path/filepath"
	"strings"
	"testing"
)

func Test_pkgModPath(t *testing.T) {
	const mockModulesTxt = `# github.com/go-gl/gl v0.0.0-20181026044259-55b76b7df9d2
github.com/go-gl/gl/v3.2-core/gl
# github.com/go-gl/glfw v0.0.0-20181213070059-819e8ce5125f
github.com/go-gl/glfw/v3.2/glfw
# github.com/goki/freetype v0.0.0-20181231101311-fa8a33aabaff
`
	type args struct {
		r      io.Reader
		module string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "No match",
			args: args{
				r:      strings.NewReader(mockModulesTxt),
				module: "fyne.io/fyne",
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "Match",
			args: args{
				r:      strings.NewReader(mockModulesTxt),
				module: glwfMod,
			},
			want:    filepath.Join(build.Default.GOPATH, "pkg/mod/github.com/go-gl/glfw@v0.0.0-20181213070059-819e8ce5125f"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cacheModPath(tt.args.r, tt.args.module)
			if (err != nil) != tt.wantErr {
				t.Errorf("cacheModPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("cacheModPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
