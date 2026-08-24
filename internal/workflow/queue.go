package workflow

import (
	"fmt"
	"gymrecommend/internal/model"
	"sync"
)

type Queue struct {
	mu   sync.Mutex
	jobs []model.Recommendation
}

func NewQueue() *Queue { return &Queue{jobs: []model.Recommendation{}} }
func (q *Queue) Enqueue(r model.Recommendation) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if r.ID == "" {
		return fmt.Errorf("job id required")
	}
	q.jobs = append(q.jobs, r)
	return nil
}
func (q *Queue) Dequeue() (model.Recommendation, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.jobs) == 0 {
		return model.Recommendation{}, false
	}
	r := q.jobs[0]
	q.jobs = q.jobs[1:]
	return r, true
}
func (q *Queue) Size() int { q.mu.Lock(); defer q.mu.Unlock(); return len(q.jobs) }
func (q *Queue) Drain(fn func(model.Recommendation) error) error {
	for {
		r, ok := q.Dequeue()
		if !ok {
			return nil
		}
		if e := fn(r); e != nil {
			return e
		}
	}
}
