// +build linux openbsd freebsd netbsd dragonfly

package report

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/lucor/goinfo"
	"golang.org/x/sys/execabs"
)

// Info returns the collected info
func (i *OS) Info() (goinfo.Info, error) {
	info, err := i.osRelease()
	if err != nil {
		return nil, err
	}

	arch, err := i.architecture()
	if err != nil {
		return info, err
	}
	info["architecture"] = arch

	kernel, err := i.kernel()
	if err != nil {
		return info, err
	}
	info["kernel"] = kernel
	return info, err
}

func (i *OS) osRelease() (goinfo.Info, error) {
	releaseFiles := []string{"/etc/os-release"}
	if matches, err := filepath.Glob("/etc/*release"); err != nil {
		releaseFiles = append(releaseFiles, matches...)
	}
	if matches, err := filepath.Glob("/etc/*version"); err != nil {
		releaseFiles = append(releaseFiles, matches...)
	}

	for _, releaseFile := range releaseFiles {
		b, err := ioutil.ReadFile(releaseFile)
		if err != nil {
			continue
		}
		return i.parseOsReleaseCmdOutput(b)
	}
	return nil, fmt.Errorf("could not found any release file: %#+v", releaseFiles)
}

func (i *OS) parseOsReleaseCmdOutput(data []byte) (goinfo.Info, error) {
	// fitlerKeys defines the key to return
	filterKeys := map[string]string{
		"NAME":     "name",
		"VERSION":  "version",
		"HOME_URL": "home_url",
	}
	info := goinfo.Info{}
	buf := bytes.NewBuffer(data)
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		tokens := strings.Split(scanner.Text(), "=")
		if len(tokens) != 2 {
			continue
		}
		key, ok := filterKeys[tokens[0]]
		if !ok {
			continue
		}
		info[key] = strings.Trim(tokens[1], `"`)
	}
	return info, scanner.Err()
}

func (i *OS) architecture() (string, error) {
	cmd := execabs.Command("uname", "-m")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("could not detect architecture using uname command: %w", err)
	}

	arch := strings.Trim(string(out), "\n")

	return arch, nil
}

func (i *OS) kernel() (string, error) {
	cmd := execabs.Command("uname", "-rsv")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("could not detect the kernel using uname command: %w", err)
	}

	kernel := strings.Trim(string(out), "\n")

	return kernel, nil
}
