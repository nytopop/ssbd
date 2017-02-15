// Package models provides models for ssbd.
package models

import (
	"database/sql"
	"io/ioutil"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nytopop/ssbd/config"
)

var db *sql.DB

// InitDB opens the sqlite database and executes a schema against it.
func InitDB() error {
	var err error
	db, err = sql.Open("sqlite3", config.CFG.Srv.DB)
	if err != nil {
		return err
	}

	schema, err := ioutil.ReadFile(config.CFG.Srv.Schema)
	if err != nil {
		return err
	}

	if _, err := db.Exec(string(schema)); err != nil {
		return err
	}

	return nil
}
