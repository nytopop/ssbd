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
	InsertVolume(v Volume) error
	UpdateVolume(v Volume) error

	// Servers
	GetServers() ([]Server, error) // TODO
	InsertServer(s Server) error   // TODO
	UpdateServer(s Server) error   // TODO

	// Jobs
	//GetJobs() ([]Job, error) // TODO
	//InsertJob(j Job) error   // TODO
	//UpdateJob(j Job) error   // TODO

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
