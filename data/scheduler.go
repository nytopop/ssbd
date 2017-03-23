package data

import (
	"fmt"
	"sync"

	"github.com/nytopop/ssbd/logs"
	"github.com/nytopop/ssbd/models"
	"github.com/robfig/cron"
)

const (
	ErrNoVol      = logs.Err("Error no volume for job.")
	ErrNoSrv      = logs.Err("Error no server for job.")
	ErrDB         = logs.Err("Error accessing database.")
	ErrNoFull     = logs.Err("Error no existing full backup.")
	EvtRescanVols = iota
	EvtRescanSrvs
	EvtRescanJobs
)

type Scheduler struct {
	// scheduling
	db    models.Handler
	wg    sync.WaitGroup
	clock sync.RWMutex
	cron  *cron.Cron

	// event channels
	jobs chan models.Job
	evts chan int
	errs chan error

	// volume management
	vlock sync.RWMutex
	vols  map[int64]VolumeHandler

	// server management
	slock sync.RWMutex
	srvs  map[int64]ServerHandler
}

func NewScheduler(db models.Handler) *Scheduler {
	return &Scheduler{
		db:   db,
		cron: cron.New(),
		jobs: make(chan models.Job),
		evts: make(chan int),
		errs: make(chan error),
		vols: make(map[int64]VolumeHandler),
		srvs: make(map[int64]ServerHandler),
	}
}

// keeps track of all helper goroutines. Will block.
func (s *Scheduler) Run() {
	// spinner on s.jobs
	go func() {
		for j := range s.jobs {
			// do we do this in a goroutine?
			s.RunJob(j)
		}
	}()

	// spinner on s.evts
	go func() {
		for e := range s.evts {
			fmt.Println(e)
			switch e {
			case EvtRescanVols:
				s.ScanVols()
			case EvtRescanSrvs:
				s.ScanSrvs()
			case EvtRescanJobs:
				s.ScanJobs()
			}
		}
	}()

	// spinner on s.errs
	go func() {
		for err := range s.errs {
			logs.Error(err)
		}
	}()
}

func (s *Scheduler) RunJob(j models.Job) {
	// check volume is available,
	s.vlock.RLock()
	if _, ok := s.vols[j.VolumeID]; !ok {
		s.vlock.RUnlock()
		s.errs <- ErrNoVol
		return
	}
	vol := s.vols[j.VolumeID]
	s.vlock.RUnlock()

	// check server is available,
	s.slock.RLock()
	if _, ok := s.srvs[j.ServerID]; !ok {
		s.slock.RUnlock()
		s.errs <- ErrNoSrv
		return
	}
	srv := s.srvs[j.ServerID]
	s.slock.RUnlock()

	r := models.Run{
		JobID:  j.JobID,
		Status: models.StatusWait,
	}
	rid, err := s.db.InsertRun(r)
	if err != nil {
		s.errs <- ErrDB
		return
	}

	switch j.Style {
	case models.Full:
		// get tar,idx writers for 'runid' : tar,idx
		tar, err := vol.GetW(rid, Tar)
		if err != nil {
			s.errs <- err
			return
		}
		defer tar.Close()

		idx, err := vol.GetW(rid, Idx)
		if err != nil {
			s.errs <- err
			return
		}
		defer idx.Close()

		// srv . GetFullTar ( tar, idx )
		err = srv.GetFullTar(tar, idx)
		if err != nil {
			s.errs <- err
			return
		}
	case models.Diff:
		// get idx reader from last full backup : fidx
		// get tar/idx writers for 'runid' : tar,idx

		fid, err := s.db.GetLastFullRunID(j.ServerID, j.Directory)
		if err != nil {
			s.errs <- ErrNoFull
			return
		}

		fidx, err := vol.GetR(fid, Idx)
		if err != nil {
			s.errs <- err
			return
		}
		defer fidx.Close()

		tar, err := vol.GetW(rid, Tar)
		if err != nil {
			s.errs <- err
			return
		}
		defer tar.Close()

		idx, err := vol.GetW(rid, Idx)
		if err != nil {
			s.errs <- err
			return
		}
		defer idx.Close()

		err = srv.GetDiffTar(fidx, tar, idx)
		if err != nil {
			s.errs <- err
			return
		}
	}

	// cleanup the run, set status to StatusGood
	// update run in db
}

// vlock, vols
func (s *Scheduler) ScanVols() {
	vols, err := s.db.GetVolumes()
	if err != nil {
		s.errs <- err
		return
	}

	vols = make([]models.Volume, 1)

	nvols := make(map[int64]VolumeHandler)
	for i := range vols {
		// if this is a good volume
		f, err := OpenFileDir("./run/bak")
		if err != nil {
			s.errs <- err
		} else {
			nvols[vols[i].VolumeID] = f
		}
	}

	s.vlock.Lock()
	s.vols = nvols
	s.vlock.Unlock()
}

// slock, srvs
func (s *Scheduler) ScanSrvs() {
}

// clock, cron
func (s *Scheduler) ScanJobs() {
	jobs, err := s.db.GetJobs()
	if err != nil {
		s.errs <- err
		return
	}
	jobs = make([]models.Job, 1)

	c := cron.New()

	for i := range jobs {
		err := c.AddFunc("@every 8s", func() {
			s.jobs <- jobs[i]
		})

		if err != nil {
			s.errs <- err
		}
	}

	s.cron.Stop()
	s.clock.Lock()
	s.cron = c
	s.clock.Unlock()
	s.cron.Start()
}

func (s *Scheduler) Close() {
}
