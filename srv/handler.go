package srv

import (
	"io"

	"github.com/nytopop/ssbd/logs"
)

const (
	ErrAuth = logs.Err("Error server authorization failed.")
	ErrDead = logs.Err("Error server not responding or blocked port.")
)

// Server Handler
type Handler interface {
	GetFullTar(dir string, tar, idx io.Writer) error
	GetDiffTar(fidx io.Reader, tar, idx io.Writer) error
	Ping() error
	Close()
}
