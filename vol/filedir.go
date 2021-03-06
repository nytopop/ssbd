package vol

import (
	"io"
	"os"
	"strconv"

	"github.com/nytopop/ssbd/logs"
)

const (
	ErrNotDir      = logs.Err("Not a directory!")
	ErrNotExist    = logs.Err("Does not exist!")
	ErrPermission  = logs.Err("Permissions don't work!")
	ErrWriteHandle = logs.Err("Error opening a write handle.")
	ErrReadHandle  = logs.Err("Error opening a read handle.")
	ErrUnknownV    = logs.Err("Error handle type unknown.")
)

type FileDir string

func OpenFileDir(dir string) (FileDir, error) {
	info, err := os.Stat(dir)
	switch {
	case os.IsNotExist(err):
		return FileDir(dir), ErrNotExist
	case os.IsPermission(err):
		return FileDir(dir), ErrPermission
	case !info.IsDir():
		return FileDir(dir), ErrNotDir
	}

	return FileDir(dir), err
}

func (d FileDir) GetW(id int64, v int) (io.WriteCloser, error) {
	// BUG TODO : don't use string concat >:|
	path := string(d) + "/" + strconv.Itoa(int(id))
	err := os.MkdirAll(path, os.ModeDir|0777)
	if err != nil {
		return nil, err
	}

	switch v {
	case Tar:
		path += "/dat.tar"
	case Idx:
		path += "/dat.idx"
	default:
		return nil, ErrUnknownV
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, ErrWriteHandle
	}

	return f, nil
}

func (d FileDir) GetR(id int64, v int) (io.ReadCloser, error) {
	// check if path exists
	// check if d + id + v + /dat.tar|/dat.idx
	return nil, nil
}

func (d FileDir) GetStats() VolStats {
	return VolStats{}
}

func (d FileDir) Ping() error {
	return nil
}

func (d FileDir) Close() {
}
