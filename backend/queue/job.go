package queue

import "time"

type Status string

const (
	StatusQueued  Status = "queued"
	StatusRunning Status = "running"
	StatusDone    Status = "done"
	StatusError   Status = "error"
)

type Job struct {
	ID        string
	Profile   string
	Status    Status
	Output    string
	Error     string
	CreatedAt time.Time
	StartedAt time.Time
	EndedAt   time.Time
}