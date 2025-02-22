package cli

import (
	"bytes"
	"fmt"
	"os/exec"
)

func RunCommand(args ...string) (string, error) {
	cmd := exec.Command("huggingface-cli", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("command failed: %s, error: %v", out.String(), err)
	}

	return out.String(), nil
}
