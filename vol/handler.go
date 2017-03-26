package vol

import "io"

const (
	Tar = iota
	Idx
)

// Volume Handler
type Handler interface {
	GetW(id int64, v int) (io.WriteCloser, error)
	GetR(id int64, v int) (io.ReadCloser, error)
	GetStats() VolStats
	Ping() error
	Close()
}

type VolStats struct {
	Capacity int64
	Free     int64
	Used     int64
}
