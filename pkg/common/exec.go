package common

import (
	"bytes"
	"context"
	"os/exec"
)

func ExecCmd(ctx context.Context, args []string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)

	output, err := cmd.Output()

	return bytes.TrimSpace(output), err
}

func ExecCmdNoOutput(ctx context.Context, args []string) error {
	_, err := ExecCmd(ctx, args)
	return err
}
