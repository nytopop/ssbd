// Package config provides server configuration for ssbd.
package config

import gcfg "gopkg.in/gcfg.v1"

// CFG contains the currently active configuration.
var CFG Config

// Config represents the available configuration options for ssbd.
type Config struct {
	Srv struct {
		DB          string
		Schema      string
		ResourceDir string
		ErrorLog    string
		AccessLog   string
		Listen      string
	}
}

// LoadConfig loads a configuration file at path into CFG.
func LoadConfig(path string) error {
	return gcfg.ReadFileInto(&CFG, path)
}
