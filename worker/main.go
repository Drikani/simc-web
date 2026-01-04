package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})

	for {
		// Blocking pop: wartet bis ein Job da ist
		res, err := rdb.BLPop(ctx, 0*time.Second, "job_queue").Result()
		if err != nil {
			log.Println("Redis error:", err)
			time.Sleep(1 * time.Second)
			continue
		}

		jobID := res[1]
		log.Println("Processing job:", jobID)

		profileFile := filepath.Join("/jobs", jobID+".simc")
		if _, err := os.Stat(profileFile); os.IsNotExist(err) {
			log.Println("Profile file does not exist:", profileFile)
			continue
		}

		rdb.Set(ctx, "job:"+jobID+":progress", "running", 0)

		cmd := exec.CommandContext(ctx, "/app/SimulationCraft/simc", profileFile)
		stdout, _ := cmd.StdoutPipe()
		stderr, _ := cmd.StderrPipe()

		if err := cmd.Start(); err != nil {
			log.Println("Failed to start SimC:", err)
			rdb.Set(ctx, "job:"+jobID+":progress", "failed", 0)
			continue
		}

		pushToRedis := func(scanner *bufio.Scanner) {
			for scanner.Scan() {
				line := scanner.Text()
				rdb.RPush(ctx, "job:"+jobID+":stream", line)
			}
		}
		go pushToRedis(bufio.NewScanner(stdout))
		go pushToRedis(bufio.NewScanner(stderr))

		if err := cmd.Wait(); err != nil {
			log.Println("SimC finished with error:", err)
			rdb.Set(ctx, "job:"+jobID+":progress", "failed", 0)
		} else {
			rdb.Set(ctx, "job:"+jobID+":progress", "done", 0)
		}
	}
}