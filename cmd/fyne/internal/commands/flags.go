package commands

import (
	"fmt"
	"strings"
)

type keyValueFlag struct {
	m map[string]string
}

func (k *keyValueFlag) Set(value string) error {
	if k.m == nil {
		k.m = make(map[string]string)
	}
	parts := strings.Split(value, "=")

	if len(parts) < 2 {
		return fmt.Errorf("expected format: key=value, got %s", value)
	}

	k.m[parts[0]] = strings.Join(parts[1:], "=")
	return nil
}

func (k *keyValueFlag) String() string {
	result := ""

	for key, value := range k.m {
		if result == "" {
			result = "\"" + key + "=" + value + "\""
		} else {
			result = result + ", \"" + key + "=" + value + "\""
		}
	}

	return result
}
