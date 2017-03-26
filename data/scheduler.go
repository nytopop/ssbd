package data

import (
	"fmt"
	"log"
	"runtime"
	"sync"

	"github.com/nytopop/ssbd/logs"
	"github.com/nytopop/ssbd/models"
	"github.com/nytopop/ssbd/srv"
	"github.com/nytopop/ssbd/vol"
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
	EvtVerify
)

/* TODO
handle EvtVerify
get server handlers working + rescan
get volume handlers working + rescan

need data input, frontend forms, etc

fix the locks, they need to lock until done with handler
*/

type Scheduler struct {
	// scheduling
	dlock sync.Mutex
	db    models.Handler
	clock sync.RWMutex
	cron  *cron.Cron

	// event channels
	jobs chan models.Job
	evts chan int
	errs chan error

	// volume management
	vlock sync.RWMutex
	vols  map[int64]vol.Handler

	// server management
	slock sync.RWMutex
	srvs  map[int64]srv.Handler
}

func NewScheduler(db models.Handler) *Scheduler {
	s := &Scheduler{
		db:   db,
		cron: cron.New(),
		jobs: make(chan models.Job),
		evts: make(chan int),
		errs: make(chan error),
		vols: make(map[int64]vol.Handler),
		srvs: make(map[int64]srv.Handler),
	}

	go s.Run()

	s.ScanVols()
	s.ScanSrvs()
	s.ScanJobs()

	return s
}

// keeps track of all helper goroutines. Will block.
func (s *Scheduler) Run() {
	n := runtime.NumCPU()

	// spawn n spinners on s.jobs
	for i := 0; i < n; i++ {
		go func(x int) {
			for j := range s.jobs {
				fmt.Println("Job", j, "arrived on", x)
				s.RunJob(j)
			}
		}(i)
	}

	// spinner on s.evts
	go func() {
		for e := range s.evts {
			switch e {
			case EvtRescanVols:
				s.ScanVols()
			case EvtRescanSrvs:
				s.ScanSrvs()
			case EvtRescanJobs:
				s.ScanJobs()
			case EvtVerify:
				// go through list of backups, verify them all
				// get all runs with
			}
		}
	}()

	// spinner on s.errs
	go func() {
		for err := range s.errs {
			log.Fatalln(err)
			logs.Error(err)
		}
	}()
}

func (s *Scheduler) RunJob(j models.Job) {
	// Setup locks
	s.vlock.RLock()
	s.slock.RLock()
	defer s.vlock.RUnlock()
	defer s.slock.RUnlock()

	// get the current run's id & set status waiting
	r := models.Run{
		JobID:  j.JobID,
		Status: models.StatusWait,
	}
	s.dlock.Lock()
	rid, err := s.db.InsertRun(r)
	s.dlock.Unlock()
	if err != nil {
		s.errs <- err
		return
	}
	r.RunID = rid
	r.Status = models.StatusFail

	// we defer updating job status
	// so good/fail are correct
	defer func(tr *models.Run) {
		s.dlock.Lock()
		err := s.db.UpdateRun(*tr)
		if err != nil {
			s.errs <- err
		}
		s.dlock.Unlock()
	}(&r)

	// check volume is available,
	if _, ok := s.vols[j.VolumeID]; !ok {
		s.errs <- ErrNoVol
		return
	}
	vlm := s.vols[j.VolumeID]

	// check server is available,
	if _, ok := s.srvs[j.ServerID]; !ok {
		s.errs <- ErrNoSrv
		return
	}
	serv := s.srvs[j.ServerID]

	switch j.Style {
	case models.Full:
		// get tar,idx writers for 'runid' : tar,idx
		tar, err := vlm.GetW(rid, vol.Tar)
		if err != nil {
			s.errs <- err
			return
		}
		defer tar.Close()

		idx, err := vlm.GetW(rid, vol.Idx)
		if err != nil {
			s.errs <- err
			return
		}
		defer idx.Close()

		// serv . GetFullTar ( tar, idx )
		err = serv.GetFullTar(j.Directory, tar, idx)
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

		fidx, err := vlm.GetR(fid, vol.Idx)
		if err != nil {
			s.errs <- err
			return
		}
		defer fidx.Close()

		tar, err := vlm.GetW(rid, vol.Tar)
		if err != nil {
			s.errs <- err
			return
		}
		defer tar.Close()

		idx, err := vlm.GetW(rid, vol.Idx)
		if err != nil {
			s.errs <- err
			return
		}
		defer idx.Close()

		err = serv.GetDiffTar(fidx, tar, idx)
		if err != nil {
			s.errs <- err
			return
		}
	}

	// set status to StatusGood before returning
	// deferred update will use the new value because
	// r referenced by pointer
	r.Status = models.StatusGood
}

// vlock, vols
func (s *Scheduler) ScanVols() {
	s.vlock.Lock()
	defer s.vlock.Unlock()

	vols, err := s.db.GetVolumes()
	if err != nil {
		s.errs <- err
		return
	}

	nvols := make(map[int64]vol.Handler)
	for i := range vols {
		switch vols[i].Backend {
		// if this is a good volume
		case models.FileDir:
			f, err := vol.OpenFileDir("./run/bak")
			if err != nil {
				s.errs <- err
			} else {
				nvols[vols[i].VolumeID] = f
			}
		}
	}

	for k := range s.vols {
		s.vols[k].Close()
	}
	s.vols = nvols
}

// slock, srvs
func (s *Scheduler) ScanSrvs() {
	s.slock.Lock()
	defer s.slock.Unlock()

	srvs, err := s.db.GetServers()
	if err != nil {
		s.errs <- err
		return
	}

	nsrvs := make(map[int64]srv.Handler)
	for i := range srvs {
		switch srvs[i].Proto {
		case models.SrvSSH:
			h, err := srv.DialSSH(srvs[i].Address, srvs[i].Port)
			if err != nil {
				s.errs <- err
			} else {
				nsrvs[srvs[i].ServerID] = h
			}
			// attempt to ssh connect
			// spawn an SSH Srv handler
			// pop into
			// use a default private key for now, to make things simple
		case models.SrvFTP:
		case models.SrvHTTP:
		case models.SrvSSBD:
		default:
			// unknown
		}
	}

	for k := range s.srvs {
		s.srvs[k].Close()
	}
	s.srvs = nsrvs
}

// clock, cron
func (s *Scheduler) ScanJobs() {
	s.clock.Lock()
	defer s.clock.Unlock()

	jobs, err := s.db.GetJobs()
	if err != nil {
		s.errs <- err
		return
	}

	c := cron.New()
	for i := range jobs {
		var k = i
		err := c.AddFunc(jobs[k].Cron, func() {
			s.jobs <- jobs[k]
		})

		if err != nil {
			s.errs <- err
		}
	}

	err = c.AddFunc("@every 60s", func() {
		fmt.Println("Rescan event")
		s.evts <- EvtRescanVols
		s.evts <- EvtRescanSrvs
		s.evts <- EvtRescanJobs
	})
	if err != nil {
		s.errs <- err
	}

	s.cron.Stop()
	s.cron = c
	s.cron.Start()
}

func (s *Scheduler) Close() {
	// Cleanup

	// finish all pending backups
	// release all server handlers
	// release all volume handlers
}
