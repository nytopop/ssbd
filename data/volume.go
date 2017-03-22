package data

import (
	"os"

	"github.com/nytopop/ssbd/logs"
)

const (
	ErrNotDir     = logs.Err("Not a directory!")
	ErrNotExist   = logs.Err("Does not exist!")
	ErrPermission = logs.Err("Permissions don't work!")
)

type VolumeHandler interface {
	GetStats() VolStats
	Close()
}

type VolStats struct {
	Capacity int64
	Free     int64
	Used     int64
	Alive    bool
}

var VolumeHandlers map[int]VolumeHandler

type FileDir string

func OpenFileDir(dir string) (FileDir, error) {
	info, err := os.Stat(dir)
	switch {
	case os.IsNotExist(err):
		return FileDir(dir), ErrNotExist
	case os.IsPermission(err):
		return FileDir(dir), ErrPermission
	case info.IsDir() == true:
		return FileDir(dir), ErrNotDir
	}

	return FileDir(dir), err
}

func (d FileDir) Close() {
}

func (d FileDir) GetStats() VolStats {
	return VolStats{}
}
