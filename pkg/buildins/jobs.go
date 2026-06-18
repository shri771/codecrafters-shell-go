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

func (s *JobStore) JobMarkers() map[*RunningJob]string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	markers := map[*RunningJob]string{}
	found := 0
	for i := len(s.jobs) - 1; i >= 0; i-- {
		job := s.jobs[i]
		if job.GetStatus() != "running" {
			continue
		}

		if found == 0 {
			markers[job] = "+"
		} else if found == 1 {
			markers[job] = "-"
		}

		found++
		if found == 2 {
			break
		}
	}

	return markers
}

func jobsCMD(args []string) error {
	jobs := DefaultJobStore.jobs
	markers := DefaultJobStore.JobMarkers()

	for _, job := range jobs {
		marker := markers[job]
		if marker == "" {
			marker = " "
		}

		fmt.Printf("[%d]%s  Running                 %s &\n", job.JobNumber, marker, job.CmdUsed)
	}
	return nil
}
