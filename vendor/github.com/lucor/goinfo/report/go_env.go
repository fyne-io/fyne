package report

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/lucor/goinfo"
)

// GoEnv collects the info about the Go environment using the go env tool
type GoEnv struct {
	Filter []string // Filter defines a subset of vars to returns. Leave empty to return all.
}

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

	if len(i.Filter) == 0 {
		return info, nil
	}

	filteredInfo := goinfo.Info{}
	for _, key := range i.Filter {
		if v, ok := info[key]; ok {
			filteredInfo[key] = v
		}
	}

	return filteredInfo, nil
}
