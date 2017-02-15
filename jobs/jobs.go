// Package jobs provides job scheduling functions for ssbd.
package jobs

import (
	"errors"
	"sync"
)

/* This module provides job scheduling and execution.

Interops with ssbd/data and jobs table in database.

Need:
Force rescan database for new jobs, otherwise scheduled rescan.


*/

var scheduler Scheduler

// Scheduler keeps track of running and scheduled jobs.
type Scheduler struct {
	wg sync.WaitGroup
}

// StartScheduler starts the scheduler. An error will be returned if
// the scheduler is already running.
func StartScheduler() error {
	if scheduler != (Scheduler{}) {
		return errors.New("Scheduler exists.")
	}

	var wg sync.WaitGroup
	scheduler = Scheduler{
		wg: wg,
	}

	go scheduler.Run()

	return nil
}

// Run starts the scheduler. This should only be called once!
func (s Scheduler) Run() {
	// wait and listen for incoming jobs on channel.
	// process on FIFO basis.
}

// Rescan scans the database for any new scheduled jobs, and pushes them.
func (s Scheduler) Rescan() error {
	return nil
}

// Status returns a status string for the scheduler. Display at top bar?
func (s Scheduler) Status() string {
	return ""
}

// Job represents a single scheduled job.
type Job struct {
}
