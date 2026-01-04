package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func Submit(rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Profile string `json:"profile"`
		}
		json.NewDecoder(r.Body).Decode(&body)

		id := uuid.NewString()
		ctx := context.Background()

		rdb.HSet(ctx, "job:"+id, map[string]any{
			"status": "queued",
			"created_at": time.Now().Unix(),
		})
		rdb.Set(ctx, "job:"+id+":profile", body.Profile, 0)
		rdb.LPush(ctx, "queue:jobs", id)

		json.NewEncoder(w).Encode(map[string]string{"job_id": id})
	}
}