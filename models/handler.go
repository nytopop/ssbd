package models

import (
	"database/sql"
	"io/ioutil"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nytopop/ssbd/config"
)

type Handler interface {
	// Users
	//	GetUsers() []models.User
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
