package data

import "io"

type ServerHandler interface {
	GetFullTar(tar io.Writer, idx io.Writer) error
	GetDiffTar(fidx io.Reader, tar io.Writer, idx io.Writer) error
}
