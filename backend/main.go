package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"strings"

	"simc-backend/parser"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var ctx = context.Background()

type JobRequest struct {
	Profile string `json:"profile"`
}

type ParseRequest struct {
	Output string `json:"output"`
}


func main() {
	r := gin.Default()

	// Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})

	// POST: neuen Job anlegen
	r.POST("/api/jobs", func(c *gin.Context) {
		var req JobRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		jobID := fmt.Sprintf("%d", time.Now().UnixNano())
		filePath := filepath.Join("/jobs", jobID+".simc")
		if err := os.WriteFile(filePath, []byte(req.Profile), 0644); err != nil {
			c.JSON(500, gin.H{"error": "Failed to write profile"})
			return
		}

		redisClient.RPush(ctx, "job_queue", jobID)
		redisClient.Set(ctx, "job:"+jobID+":progress", "queued", 0)

		c.JSON(200, gin.H{"job_id": jobID})
	})

	// GET: Fortschritt
	r.GET("/api/jobs/:id/progress", func(c *gin.Context) {
		id := c.Param("id")
		progress, _ := redisClient.Get(ctx, "job:"+id+":progress").Result()
		c.JSON(200, gin.H{"progress": progress})
	})

	// GET: Result für Formatierten Output
	r.GET("/api/jobs/:id/result", func(c *gin.Context) {
		id := c.Param("id")

		output, err := redisClient.Get(ctx, "job:"+id+":result").Result()
		if err != nil || strings.TrimSpace(output) == "" {
			c.JSON(404, gin.H{"error": "no result yet"})
			return
		}

		c.JSON(200, gin.H{
			"output": output,
		})
	})


	// SSE: Live-Output streamen
	r.GET("/api/jobs/:id/stream", func(c *gin.Context) {
		id := c.Param("id")
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Flush()

		for {
			lines, _ := redisClient.LRange(ctx, "job:"+id+":stream", 0, -1).Result()
			for _, line := range lines {
				fmt.Fprintf(c.Writer, "event: output\ndata: %s\n\n", line)
			}
			redisClient.Del(ctx, "job:"+id+":stream")
			c.Writer.Flush()

			progress, _ := redisClient.Get(ctx, "job:"+id+":progress").Result()
			if progress == "done" || progress == "failed" {
				fmt.Fprintf(c.Writer, "event: status\ndata: %s\n\n", progress)
				c.Writer.Flush()
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	})

	// NEU: SimC Parser Route
	r.POST("/api/parse-simc", func(c *gin.Context) {
		var req ParseRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid input"})
			return
		}

		if strings.TrimSpace(req.Output) == "" {
			c.JSON(400, gin.H{"error": "empty output"})
			return
		}

		result := parser.ParseSimC(req.Output)
		c.JSON(200, result)
	})

	r.Run(":8080")
}