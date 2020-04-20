package report

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/lucor/goinfo"
)

// Info returns the collected info about the OS
func (i *OS) Info() (goinfo.Info, error) {
	cmd := exec.Command("cmd", "/C", "wmic os get /value")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("could not detect os info using wmic command: %w", err)
	}

	return i.parseWmicCmdOutput(out)
}

func (i *OS) parseWmicCmdOutput(data []byte) (goinfo.Info, error) {
	// fitlerKeys defines the key to return
	filterKeys := map[string]string{
		"Caption":        "name",
		"Version":        "version",
		"OSArchitecture": "architecture",
	}
	info := goinfo.Info{}
	buf := bytes.NewBuffer(data)
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), "\r\n")
		tokens := strings.Split(line, "=")
		if len(tokens) != 2 {
			continue
		}
		key, ok := filterKeys[tokens[0]]
		if !ok {
			continue
		}
		info[key] = tokens[1]
	}
	return info, scanner.Err()
}
