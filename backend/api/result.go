package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

func Result(rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		out := rdb.LRange(context.Background(), "job:"+id+":output", 0, -1).Val()
		for _, line := range out {
			w.Write([]byte(line + "\n"))
		}
	}
}