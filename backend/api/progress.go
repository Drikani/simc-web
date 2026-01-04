package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

func Progress(rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		p := rdb.Get(context.Background(), "job:"+id+":progress").Val()
		w.Write([]byte(p))
	}
}