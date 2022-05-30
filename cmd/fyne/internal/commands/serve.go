package commands

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/urfave/cli/v2"
)

// Server serve fyne wasm application over http
type Server struct {
	*appData
	port            int
	srcDir, dir, os string
}

// Serve return the cli command for serving fyne wasm application over http
func Serve() *cli.Command {
	s := &Server{appData: &appData{}}

	return &cli.Command{
		Name:        "serve",
		Usage:       "Package an application using WebAssembly and expose it via a web server.",
		Description: `The serve command packages an application using WebAssembly and expose it via a web server which port can be overridden with port.`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "sourceDir",
				Usage:       "The directory to package, if executable is not set",
				Value:       "",
				Destination: &s.srcDir,
			},
			&cli.StringFlag{
				Name:        "icon",
				Usage:       "The name of the application icon file.",
				Destination: &s.icon,
			},
			&cli.IntFlag{
				Name:        "port",
				Usage:       "The port to have the http server listen on (default: 8080).",
				Value:       8080,
				Destination: &s.port,
			},
			&cli.StringFlag{
				Name:        "target",
				Aliases:     []string{"os"},
				Usage:       "The web runtime to target (wasm, gopherjs, web).",
				Value:       "web",
				Destination: &s.os,
			},
		},
		Action: s.Server,
	}
}

func (s *Server) requestPackage() error {
	p := &Packager{
		os:     s.os,
		srcDir: s.srcDir,

		appData: s.appData,
	}

	err := p.Package()
	s.dir = p.dir
	return err
}

func (s *Server) serve() error {
	err := s.validate()
	if err != nil {
		return err
	}

	err = s.requestPackage()
	if err != nil {
		return err
	}

	webDir := util.EnsureSubDir(s.dir, s.os)
	fileServer := http.FileServer(http.Dir(webDir))

	http.Handle("/", fileServer)

	fmt.Printf("Serving %s on HTTP port: %v\n", webDir, s.port)
	err = http.ListenAndServe(":"+strconv.Itoa(s.port), nil)

	return err
}

// Server serve the wasm application over http
func (s *Server) Server(ctx *cli.Context) error {
	if ctx.Args().Len() != 0 {
		return errors.New("unexpected parameter after flags")
	}

	return s.serve()
}

func (s *Server) validate() error {
	if s.port == 0 {
		s.port = 8080
	}
	if s.port < 0 || s.port > 65535 {
		return fmt.Errorf("the port must be a strictly positive number and be strictly smaller than 65536 (Got %v)", s.port)
	}
	if s.os != "wasm" && s.os != "gopherjs" && s.os != "web" {
		return fmt.Errorf("unsupported web runtime (only wasm, gopherjs and web): %v", s.os)
	}
	return nil
}
