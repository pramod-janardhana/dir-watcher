package jobs

import (
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

var jobscheduler *jobScheduler = nil

type Job struct {
	id            int
	CronExpresion string
	Handler       func()
	lable         string
	startAt       time.Time
	endAt         time.Time
}

type jobScheduler struct {
	s    *cron.Cron
	jobs map[*Job]struct{}
}

// Start the cron scheduler in its own goroutine, or no-op if already started.
func (s *jobScheduler) Start() {
	s.s.Start()
}

// Stop stops the cron scheduler if it is running; otherwise it does nothing.
func (s *jobScheduler) Stop() {
	// wait for the cron jobs to stop
	s.s.Stop()
	jobscheduler = nil
}

// AddJob adds a job to the Cron to be run on the given schedule.
func (s *jobScheduler) AddJob(job *Job) error {
	jobId, err := s.s.AddFunc(job.CronExpresion, job.Handler)
	if err != nil {
		return err
	}

	job.id = int(jobId)
	s.jobs[job] = struct{}{}

	return nil
}

// RemoveJob removes the job from being run in the future.
func (s *jobScheduler) RemoveJob(job *Job) {
	s.s.Remove(cron.EntryID(job.id))
}

// RemoveAllJob removes all the job from being run in the future.
func (s *jobScheduler) RemoveAllJob() {
	for job := range s.jobs {
		s.s.Remove(cron.EntryID(job.id))
	}
}

// NewFileWatcher creates a new FileWatcher instance if it doesn't already exist
func NewJobScheduler() *jobScheduler {
	var lock = &sync.Mutex{}

	if jobscheduler == nil {
		lock.Lock()
		defer lock.Unlock()
		if jobscheduler == nil {

			// create a new job scheduler
			jobscheduler = &jobScheduler{
				s:    cron.New(),
				jobs: make(map[*Job]struct{}),
			}
		}
	}

	return jobscheduler
}

func GetJobScheduler() *jobScheduler {
	return jobscheduler
}
