package models

import (
	"database/sql"
	"io/ioutil"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nytopop/ssbd/config"
	"github.com/nytopop/ssbd/logs"
)

const (
	ErrQueryFailed = logs.Err("Query failed to execute.")
	ErrScan        = logs.Err("Failed to scan query rows.")
	ErrNotFound    = logs.Err("Query not found.")
	ErrConFail     = logs.Err("DB connection failed.")
)

type Handler interface {
	// Volumes
	GetVolumes() ([]Volume, error)
	InsertVolume(v Volume) (int64, error)
	UpdateVolume(v Volume) error

	// Servers
	GetServers() ([]Server, error)
	InsertServer(s Server) (int64, error)
	UpdateServer(s Server) error

	// Jobs
	GetJobs() ([]Job, error)
	InsertJob(j Job) (int64, error)
	UpdateJob(j Job) error

	// RunHistory
	GetRuns() ([]Run, error)
	GetLastFullRunID(sid int64, dir string) (int64, error)
	InsertRun(r Run) (int64, error)
	UpdateRun(r Run) error

	// ActionHistory

	Close() error
}

type Client struct {
	DB *sql.DB
}

func NewClient() (*Client, error) {
	db, err := sql.Open("sqlite3", config.CFG.Srv.DB)
	if err != nil {
		return nil, err
	}

	schema, err := ioutil.ReadFile(config.CFG.Srv.Schema)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(string(schema)); err != nil {
		return nil, err
	}

	return &Client{DB: db}, nil
}

// Close closes the underlying database connection.
func (c *Client) Close() error {
	return c.DB.Close()
}
