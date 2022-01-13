package commands

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/urfave/cli/v2"

	"fyne.io/fyne/v2/cmd/fyne/internal/util"
)

// Server serve fyne wasm application over http
type Server struct {
	port              int
	srcDir, icon, dir string
}

// Serve return the cli command for serving fyne wasm application over http
func Serve() *cli.Command {
	s := &Server{}

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
				Value:       "Icon.png",
				Destination: &s.icon,
			},
			&cli.IntFlag{
				Name:        "port",
				Usage:       "The port to have the http server listen on (default: 8080).",
				Value:       8080,
				Destination: &s.port,
			},
		},
		Action: s.Server,
	}
}

// AddFlags adds the flags for interacting with the Installer.
//
// Deprecated: Access to the individual cli commands are being removed.
func (s *Server) AddFlags() {
	flag.IntVar(&s.port, "port", 8080, "The port to have the http server listen on (default: 8080).")
	flag.StringVar(&s.srcDir, "sourceDir", "", "The directory to package, if executable is not set")
	flag.StringVar(&s.icon, "icon", "Icon.png", "The name of the application icon file")
}

// PrintHelp prints the help for the install command.
//
// Deprecated: Access to the individual cli commands are being removed.
func (s *Server) PrintHelp(indent string) {
	fmt.Println(indent, "The serve command setup a local http server that is necessary to test WebAssembly package.")
	fmt.Println(indent, "Command usage: fyne serve [parameters]")
}

func (s *Server) requestPackage() {
	p := &Packager{
		os:     "wasm",
		srcDir: s.srcDir,
		icon:   s.icon,
	}

	p.Run([]string{})
	s.dir = p.dir
}

func (s *Server) serve() error {
	err := s.validate()
	if err != nil {
		return err
	}

	s.requestPackage()

	webDir := util.EnsureSubDir(s.dir, "wasm")
	fileServer := http.FileServer(http.Dir(webDir))

	http.Handle("/", fileServer)

	fmt.Printf("Serving %s on HTTP port: %v\n", webDir, s.port)
	err = http.ListenAndServe(":"+strconv.Itoa(s.port), nil)

	return err
}

// Run runs the install command.
//
// Deprecated: A better version will be exposed in the future.
func (s *Server) Run(_ []string) {
	err := s.serve()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
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
	return nil
}
