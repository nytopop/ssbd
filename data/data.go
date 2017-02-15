// Package data provides data indexing, archiving, and gathering functions.
package data

import (
	"bufio"
	"io"
	"os"

	"github.com/pkg/sftp"

	"golang.org/x/crypto/ssh"
)

/* Notes and TODO
test that a differential backup tracks new files
save a sha256 hash of index and tar files
  backup/
  data.tar data.sha256 index index.sha256
fix directory manipulation
BUG string concatenation for running commands is a BAD IDEA
*/

// GetDiffTar creates a differential backup of src on c, comparing it with
// the file index in full, saving the resulting tar and index to dst.
func GetDiffTar(c *ssh.Client, full, src, dst string) error {
	// Copy the index to remote node
	fullidx, err := os.Open(full + "/index")
	if err != nil {
		return err
	}
	defer fullidx.Close()

	sf, err := sftp.NewClient(c)
	if err != nil {
		return err
	}
	defer sf.Close()

	ridx, err := sf.Create("/tmp/backup.index")
	if err != nil {
		return err
	}
	defer ridx.Close()

	_, err = io.Copy(ridx, fullidx)
	if err != nil {
		return err
	}

	// Make local files available
	err = os.MkdirAll(dst, 0700)
	if err != nil {
		return err
	}

	tarfile, err := os.Create(dst + "/data.tar")
	if err != nil {
		return err
	}
	defer tarfile.Close()

	// Get tar from client
	s, err := c.NewSession()
	if err != nil {
		return err
	}

	w := bufio.NewWriter(tarfile)
	defer w.Flush()
	s.Stdout = w

	command := "tar -cpf - -g /tmp/backup.index " + src
	err = s.Run(command)
	if err != nil {
		return err
	}

	// we leave the index
	return nil
}

// GetFullTar will tar up src on c, and pipe the tar file into dst locally.
func GetFullTar(c *ssh.Client, src, dst string) error {
	// Make local files available
	err := os.MkdirAll(dst, 0700)
	if err != nil {
		return err
	}

	tarfile, err := os.Create(dst + "/data.tar")
	if err != nil {
		return err
	}
	defer tarfile.Close()

	// Get tar from client
	s, err := c.NewSession()
	if err != nil {
		return err
	}
	w := bufio.NewWriter(tarfile)
	defer w.Flush()
	s.Stdout = w

	command := "tar cpf - -g /tmp/backup.index --level=0 " + src
	err = s.Run(command)
	if err != nil {
		return err
	}

	// Get index from client
	idxfile, err := os.Create(dst + "/index")
	if err != nil {
		return err
	}
	defer idxfile.Close()

	sf, err := sftp.NewClient(c)
	if err != nil {
		return err
	}
	defer sf.Close()

	idx, err := sf.Open("/tmp/backup.index")
	if err != nil {
		return err
	}
	defer idx.Close()

	// Write index locally
	_, err = idx.WriteTo(idxfile)
	if err != nil {
		return err
	}

	return nil
}
