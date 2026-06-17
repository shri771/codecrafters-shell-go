package buildins

import "sync"

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
		if job.Status == Running {
			count++
		}
	}

	return count
}

func jobsCMD(args []string) error {
	return nil
}
