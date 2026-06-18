package buildins

import (
	"fmt"
	"sync"
)

type JobStore struct {
	mu   sync.RWMutex
	jobs []*RunningJob
}

var DefaultJobStore = NewJobStore()

func NewJobStore() *JobStore {
	return &JobStore{
		jobs: []*RunningJob{},
	}
}

func (s *JobStore) Add(job *RunningJob) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.jobs = append(s.jobs, job)
}

func (s *JobStore) RunningCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	count := 0
	for _, job := range s.jobs {
		job.mu.RLock()
		if job.Status == Running {
			count++
		}
		job.mu.RUnlock()
	}

	return count
}

func jobsCMD(args []string) error {
	jobs := DefaultJobStore.jobs

	for _, job := range jobs {
		fmt.Printf("[%d]+  Running                 %s\n", job.JobNumber, job.CmdUsed)
	}
	return nil
}
