package buildins

import (
	"os/exec"
	"sync"
)

type JobStatus int

const (
	Running JobStatus = iota
	Stopped
	Done
	Failed
)

type RunningJob struct {
	mu          sync.RWMutex
	JobNumber   int
	Name        string
	CmdUsed     string
	PID         int
	Status      JobStatus
	IsDisplayed bool
	Cmd         *exec.Cmd
}

func CreateJob() *RunningJob {
	return &RunningJob{}
}

func (j *RunningJob) SetJobNumber(jobNumber int) {
	j.mu.Lock()
	defer j.mu.Unlock()

	j.JobNumber = jobNumber
}

func (j *RunningJob) GetJobNumber() int {
	j.mu.RLock()
	defer j.mu.RUnlock()

	return j.JobNumber
}

func (j *RunningJob) SetName(name string) {
	j.mu.Lock()
	defer j.mu.Unlock()

	j.Name = name
}

func (j *RunningJob) GetName() string {
	j.mu.RLock()
	defer j.mu.RUnlock()

	return j.Name
}

func (j *RunningJob) SetCmdUsed(cmdUsed string) {
	j.mu.Lock()
	defer j.mu.Unlock()

	j.CmdUsed = cmdUsed
}

func (j *RunningJob) GetCmdUsed() string {
	j.mu.RLock()
	defer j.mu.RUnlock()

	return j.CmdUsed
}

func (j *RunningJob) SetPID(pid int) {
	j.mu.Lock()
	defer j.mu.Unlock()

	j.PID = pid
}

func (j *RunningJob) GetPID() int {
	j.mu.RLock()
	defer j.mu.RUnlock()

	return j.PID
}

func (j *RunningJob) SetCmd(cmd *exec.Cmd) {
	j.mu.Lock()
	defer j.mu.Unlock()

	j.Cmd = cmd
}

func (j *RunningJob) GetCmd() *exec.Cmd {
	j.mu.RLock()
	defer j.mu.RUnlock()

	return j.Cmd
}

func (j *RunningJob) SetStatus(status JobStatus) {
	j.mu.Lock()
	defer j.mu.Unlock()

	j.Status = status
}

func (j *RunningJob) GetStatus() string {
	j.mu.RLock()
	defer j.mu.RUnlock()

	switch j.Status {
	case Running:
		return "Running"
	case Stopped:
		return "Stopped"
	case Done:
		return "Done"
	case Failed:
		return "Failed"
	default:
		return "Unknown"
	}
}

func (j *RunningJob) GetIsDisplayed() bool {
	j.mu.RLock()
	defer j.mu.RUnlock()

	return j.IsDisplayed
}

func (j *RunningJob) SetIsDisplayed(isDisplayed bool) {
	j.mu.Lock()
	defer j.mu.Unlock()

	j.IsDisplayed = isDisplayed
}
