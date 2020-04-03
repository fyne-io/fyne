package main

import (
	"os"
	"fmt"
	"flag"
	"strconv"
	"net/http"

	 "github.com/pkg/errors"
)

type server struct {
	port	int
	srcDir,icon,dir  string
}

func (s *server) addFlags() {
	flag.IntVar(&s.port, "port", 8080, "The port to have the http server listen on (default: 8080).")
	flag.StringVar(&s.srcDir, "sourceDir", "", "The directory to package, if executable is not set")
	flag.StringVar(&s.icon, "icon", "Icon.png", "The name of the application icon file")
}

func (s *server) printHelp(indent string) {
	fmt.Println(indent, "The server command setup a local http server that is necessary to test WebAssembly package.")
	fmt.Println(indent, "Command usage: fyne server [parameters]")
}

func (s *server) requestPackage() {
	p := &packager{
		os: "wasm",
		srcDir: s.srcDir,
		icon: s.icon,
	}

	p.run([]string{})
	s.dir = p.dir
}

func (s *server) run(_ []string) {
	err := s.validate()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	s.requestPackage()

	webDir := ensureSubDir(s.dir, "wasm")

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	fmt.Printf("Serving %s on HTTP port: %v\n", webDir, s.port)
	http.ListenAndServe(":"+strconv.Itoa(s.port), nil)
}

func (s *server) validate() error {
	if s.port == 0 {
		s.port = 8080
	}
	if s.port < 0 || s.port > 65535 {
		return errors.Errorf("The port must be a strictly positive number and be strictly smaller than 65536 (Got %v).", s.port)
	}
	return nil
}
