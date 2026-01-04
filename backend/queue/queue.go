package queue

import (
	"sync"
)

type Queue struct {
	jobs map[string]*Job
	mu   sync.RWMutex
	ch   chan *Job
}

func NewQueue(buffer int) *Queue {
	return &Queue{
		jobs: make(map[string]*Job),
		ch:   make(chan *Job, buffer),
	}
}

func (q *Queue) Add(job *Job) {
	q.mu.Lock()
	q.jobs[job.ID] = job
	q.mu.Unlock()

	q.ch <- job
}

func (q *Queue) Get(id string) (*Job, bool) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	job, ok := q.jobs[id]
	return job, ok
}

func (q *Queue) Channel() <-chan *Job {
	return q.ch
}