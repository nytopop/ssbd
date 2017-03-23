package data

import "io"

type VolumeHandler interface {
	GetW(id int64, v int) (io.WriteCloser, error)
	GetR(id int64, v int) (io.ReadCloser, error)
}

type VolStats struct {
	Capacity int64
	Free     int64
	Used     int64
	Alive    bool
}
