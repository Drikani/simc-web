package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"simc-backend/api"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	r := chi.NewRouter()
	r.Post("/api/jobs", api.Submit(rdb))
	r.Get("/api/jobs/{id}", api.Status(rdb))
	r.Get("/api/jobs/{id}/progress", api.Progress(rdb))
	r.Get("/api/jobs/{id}/result", api.Result(rdb))

	log.Println("Backend listening on :8080")
	http.ListenAndServe(":8080", r)
}