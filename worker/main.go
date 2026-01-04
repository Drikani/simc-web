package main

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"simc-worker/simc"
)

func main() {
	rdb := redis.NewClient(&redis.Options{Addr: "redis:6379"})
	ctx := context.Background()

	for {
		id, err := rdb.BRPop(ctx, 0, "queue:jobs").Result()
		if err != nil {
			log.Println(err)
			continue
		}

		jobID := id[1]
		profile := rdb.Get(ctx, "job:"+jobID+":profile").Val()

		rdb.HSet(ctx, "job:"+jobID, "status", "running")

		simc.Run(ctx, rdb, jobID, profile)

		rdb.HSet(ctx, "job:"+jobID, "status", "done")
		time.Sleep(100 * time.Millisecond)
	}
}