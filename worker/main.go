package main

import (
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
		// Blocking pop: wartet auf einen Job
		res, err := rdb.BLPop(ctx, 0*time.Second, "job_queue").Result()
		if err != nil {
			log.Println("Redis error:", err)
			time.Sleep(1 * time.Second)
			continue
		}

		jobID := res[1]
		log.Println("Processing job:", jobID)

		profileFile := filepath.Join("/jobs", jobID+".simc")
		outputFile := filepath.Join("/jobs", jobID+".json")

		if _, err := os.Stat(profileFile); os.IsNotExist(err) {
			log.Println("Profile file does not exist:", profileFile)
			rdb.Set(ctx, "job:"+jobID+":progress", "failed", 0)
			continue
		}

		rdb.Set(ctx, "job:"+jobID+":progress", "running", 0)

		// SimC mit json2 direkt aufrufen
		cmd := exec.CommandContext(ctx, "/app/SimulationCraft/simc", profileFile, "json2="+outputFile)
		if err := cmd.Run(); err != nil {
			log.Println("SimC execution failed:", err)
			rdb.Set(ctx, "job:"+jobID+":progress", "failed", 0)
			continue
		}

		log.Printf("SimC finished successfully:", jobID)
		rdb.Set(ctx, "job:"+jobID+":progress", "done", 0)
	}
}