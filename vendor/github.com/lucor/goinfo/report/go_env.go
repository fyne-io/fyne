package report

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/lucor/goinfo"
)

// GoEnv collects the info about the Go environment using the go env tool
type GoEnv struct{}

// Summary return the summary
func (i *GoEnv) Summary() string {
	return "Go environment info"
}

// Info returns the collected info
func (i *GoEnv) Info() (goinfo.Info, error) {
	var info goinfo.Info
	cmd := exec.Command("go", "env", "-json")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("could not detect go env info: %w", err)
	}

	err = json.Unmarshal(out, &info)
	if err != nil {
		return nil, fmt.Errorf("could not decode the go env json info: %w", err)
	}
	return info, nil
}
