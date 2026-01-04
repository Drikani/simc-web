package simc

import (
	"bufio"
	"bytes"
	"context"
	"os/exec"
	"strings"

	"github.com/redis/go-redis/v9"
)

func Run(ctx context.Context, rdb *redis.Client, jobID, profile string) {
	cmd := exec.Command("simc")
	cmd.Stdin = bytes.NewBufferString(profile)

	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()

		rdb.RPush(ctx, "job:"+jobID+":output", line)

		if strings.Contains(line, "%") {
			rdb.Set(ctx, "job:"+jobID+":progress", line, 0)
		}
	}

	cmd.Wait()
}