package elevate

import (
	"os/exec"
)

// No op
func buildCommand(prompt string, icon string, name string, args ...string) (*exec.Cmd, error) {
	return nil, nil
}
