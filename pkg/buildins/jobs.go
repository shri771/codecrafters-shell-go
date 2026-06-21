package buildins

import (
	"fmt"
	"io"
	"sync"
	"time"
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

func (s *JobStore) Add(job *RunningJob) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	jobNumber := 1
	for _, existingJob := range s.jobs {
		if existingJob.GetJobNumber() >= jobNumber {
			jobNumber = existingJob.GetJobNumber() + 1
		}
	}

	job.SetJobNumber(jobNumber)
	s.jobs = append(s.jobs, job)

	return jobNumber
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

// RefreshStatuses is now a no-op — job statuses are updated by the goroutine
// launched in runInBackground via cmd.Wait().
func (s *JobStore) RefreshStatuses() {}

func (s *JobStore) ReapCompleted() {
	s.RefreshStatuses()
	markers := s.JobMarkers()

	s.mu.Lock()
	defer s.mu.Unlock()

	remaining := s.jobs[:0]
	for _, job := range s.jobs {
		status := job.GetStatus()
		if status != "Done" && status != "Failed" {
			remaining = append(remaining, job)
			continue
		}

		if !job.GetIsDisplayed() {
			marker := markers[job]
			if marker == "" {
				marker = " "
			}

			fmt.Printf("[%d]%s  %s                 %s\n",
				job.GetJobNumber(), marker, status, job.GetCmdUsed())
		}
	}

	s.jobs = remaining
}

func ReapCompletedJobs() {
	DefaultJobStore.ReapCompleted()
}

func (s *JobStore) JobMarkers() map[*RunningJob]string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	markers := map[*RunningJob]string{}
	found := 0
	for i := len(s.jobs) - 1; i >= 0; i-- {
		job := s.jobs[i]
		if job.GetIsDisplayed() {
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

func jobsCMD(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	// Give newly started background processes a chance to exit before taking
	// the non-blocking status snapshot.
	time.Sleep(10 * time.Millisecond)
	DefaultJobStore.RefreshStatuses()

	jobs := DefaultJobStore.jobs
	markers := DefaultJobStore.JobMarkers()

	for _, job := range jobs {
		marker := markers[job]
		if marker == "" {
			marker = " "
		}

		if job.GetIsDisplayed() {
			continue
		}

		status := job.GetStatus()
		suffix := ""
		if status == "Running" || status == "Stopped" {
			suffix = " &"
		}
		fmt.Fprintf(stdout, "[%d]%s  %s                 %s%s\n",
			job.GetJobNumber(), marker, status, job.GetCmdUsed(), suffix)

		if status == "Done" || status == "Failed" {
			job.SetIsDisplayed(true)
		}
	}
	return nil
}
