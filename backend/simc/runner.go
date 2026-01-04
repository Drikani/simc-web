package simc

import (
	"bytes"
	"context"
	"os/exec"
	"time"
)

func Run(profile string) (*Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(
		ctx,
		"docker", "compose", "run", "--rm",
		"--no-deps",
		"simc",
		"simc",
	)

	cmd.Stdin = bytes.NewBufferString(profile)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	return &Result{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}, err
}