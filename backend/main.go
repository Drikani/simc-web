package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var ctx = context.Background()

type JobRequest struct {
	Profile string `json:"profile"`
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

		// Job in Redis anlegen
		redisClient.RPush(ctx, "job_queue", jobID)
		redisClient.Set(ctx, "job:"+jobID+":progress", "queued", 0)

		c.JSON(200, gin.H{"job_id": jobID})
	})

	// SSE: Live-Ergebnis streamen, sobald JSON existiert
	r.GET("/api/jobs/:id/stream", func(c *gin.Context) {
		id := c.Param("id")
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Flush()

		jsonFile := filepath.Join("/jobs", id+".json")
		simcFile := filepath.Join("/jobs", id+".simc")

		for {
			// Prüfen, ob Worker fertig ist
			progress, _ := redisClient.Get(ctx, "job:"+id+":progress").Result()
			if progress == "done" {
				// JSON-Datei lesen und an Client senden
				data, err := os.ReadFile(jsonFile)
				if err != nil {
					fmt.Fprintf(c.Writer, "data: %s\n\n", fmt.Sprintf(`{"error": "Failed to read result: %v"}`, err))
				} else {
					fmt.Fprintf(c.Writer, "data: %s\n\n", data)
				}
				c.Writer.Flush()

				// Dateien löschen
				os.Remove(simcFile)
				os.Remove(jsonFile)
				break
			} else if progress == "failed" {
				fmt.Fprintf(c.Writer, "data: %s\n\n", `{"error": "Job failed"}`)
				c.Writer.Flush()
				// Datei löschen, falls vorhanden
				os.Remove(simcFile)
				os.Remove(jsonFile)
				break
			}

			time.Sleep(200 * time.Millisecond)
		}
	})

	r.Run(":8080")
}