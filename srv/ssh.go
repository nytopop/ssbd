package srv

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"

	"github.com/pkg/sftp"

	"golang.org/x/crypto/ssh"
)

type SSH struct {
	c *ssh.Client
	s *sftp.Client
}

func DialSSH(addr []byte, port int) (*SSH, error) {
	uri := strconv.Itoa(int(addr[0])) + "." +
		strconv.Itoa(int(addr[1])) + "." +
		strconv.Itoa(int(addr[2])) + "." +
		strconv.Itoa(int(addr[3])) + ":" + strconv.Itoa(port)

	key, err := ioutil.ReadFile("/home/eric/.ssh/id_rsa")
	if err != nil {
		return nil, err
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	cfg := &ssh.ClientConfig{
		User: "eric",
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
	}

	c, err := ssh.Dial("tcp", uri, cfg)
	if err != nil {
		return nil, err
	}

	return &SSH{c: c}, nil
}

func (s *SSH) GetFullTar(dir string, tar, idx io.Writer) error {
	tw := bufio.NewWriter(tar)
	defer tw.Flush()
	iw := bufio.NewWriter(idx)
	defer iw.Flush()

	h := sha256.New()
	mw := io.MultiWriter(tar, h)

	sesh, err := s.c.NewSession()
	if err != nil {
		return err
	}
	sesh.Stdout = mw
	defer sesh.Close()

	cmd := "tar cpf - -g /tmp/backup.idx --level=0 " + dir
	err = sesh.Run(cmd)
	if err != nil {
		return err
	}

	fmt.Println(hex.EncodeToString(h.Sum(nil)))

	// get the index
	sf, err := sftp.NewClient(s.c)
	if err != nil {
		return err
	}
	defer sf.Close()

	cIdx, err := sf.Open("/tmp/backup.idx")
	if err != nil {
		return err
	}
	defer cIdx.Close()

	h.Reset()
	mw = io.MultiWriter(idx, h)

	_, err = cIdx.WriteTo(mw)
	if err != nil {
		return err
	}

	fmt.Println(hex.EncodeToString(h.Sum(nil)))

	return nil
}

func (s *SSH) GetDiffTar(fidx io.Reader, tar, idx io.Writer) error {
	return nil
}

func (s *SSH) Ping() error {
	return nil
}

func (s *SSH) Close() {
	s.c.Close()
}
