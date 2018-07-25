package elevate

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/getlantern/byteexec"
	"github.com/getlantern/elevate/bin"
)

func buildCommand(prompt string, icon string, name string, args ...string) (*exec.Cmd, error) {
	argsLen := len(args)
	if icon != "" {
		argsLen += 1
	}
	if prompt != "" {
		argsLen += 1
	}
	allArgs := make([]string, 0, argsLen)
	allArgs = append(allArgs, "-w") // wait for termination
	allArgs = append(allArgs, name)
	allArgs = append(allArgs, args...)
	cocoasudo, err := bin.Asset("elevate.exe")
	if err != nil {
		return nil, fmt.Errorf("Unable to load elevate.exe: %v", err)
	}
	_, program := filepath.Split(os.Args[0])
	be, err := byteexec.New(cocoasudo, program)
	if err != nil {
		return nil, fmt.Errorf("Unable to build byteexec for cocoasudo: %v", err)
	}

	return be.Command(allArgs...), nil
}
