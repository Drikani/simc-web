package main

import (
	"bufio"
	"bytes"
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

		var fullOutput bytes.Buffer

		if err := cmd.Start(); err != nil {
			log.Println("Failed to start SimC:", err)
			rdb.Set(ctx, "job:"+jobID+":progress", "failed", 0)
			continue
		}

		push := func(scanner *bufio.Scanner) {
			for scanner.Scan() {
				line := scanner.Text()

				// 🔴 Live-Stream
				rdb.RPush(ctx, "job:"+jobID+":stream", line)

				// 🟢 Finaler Output
				fullOutput.WriteString(line + "\n")
			}
		}

		go push(bufio.NewScanner(stdout))
		go push(bufio.NewScanner(stderr))

		if err := cmd.Wait(); err != nil {
			log.Println("SimC finished with error:", err)
			rdb.Set(ctx, "job:"+jobID+":progress", "failed", 0)
		} else {
			log.Println("SimC finished successfully:", jobID)

			// ✅ FINALEN OUTPUT SPEICHERN
			rdb.Set(ctx, "job:"+jobID+":result", fullOutput.String(), 0)

			// Cleanup
			os.Remove(profileFile)

			rdb.Set(ctx, "job:"+jobID+":progress", "done", 0)
		}
	}
}