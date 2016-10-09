package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/jasonpuglisi/ircutil"
)

// Config stores user-defined configuration values loaded from a separate file.
type Config struct {
	// (Optional) Default command settings. Can be overrode for individual
	// commands. Default: All nested defaults
	Settings ircutil.Settings `json:"settings"`
	// List of servers that can be used in a client.
	Servers []ircutil.Server `json:"servers"`
	// List of users that can be used in a client.
	Users []ircutil.User `json:"users"`
	// List of clients. Connections between user and server.
	Clients []ircutil.Client `json:"clients"`
	// List of commands to be used for all clients. To use different commands
	// with different clients, run multiple instances of the program with
	// different configuration files.
	Commands []ircutil.Command `json:"commands"`
}

// getConfig opens a config file at the given path and parses it into a config
// struct with default values applied.
func getConfig(path string) (*Config, error) {
	// Attempt to open configuration file.
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Create configuration with default values before file data is parsed.
	config := newConfig()

	// Parse file data into configuration, update dependent default values, and
	// parse it once more. We do this so that any booleans specifically set to
	// false in the config are not overrode by default values.
	json.Unmarshal(raw, config)
	setConfigDefaults(config)
	json.Unmarshal(raw, config)

	// Return parsed and updated configuration.
	return config, nil
}

// newConfig creates a config struct and populates the default settings.
func newConfig() *Config {
	return &Config{
		Settings: ircutil.Settings{
			CaseSensitive: true,
			Symbol:        "!",
			Scope:         []string{"channel"},
		},
	}
}

// setConfigDefaults updates default configuration values that are dependent on
// parsed config file data.
func setConfigDefaults(config *Config) {
	// Update defaults for each server.
	for i := range config.Servers {
		s := &config.Servers[i]
		if s.Port == 0 {
			s.Port = 6697
		}
		if s.Secure == false && s.Port == 6697 {
			s.Secure = true
		}
	}

	// Update defaults for each user.
	for i := range config.Users {
		u := &config.Users[i]
		if len(u.User) < 1 {
			u.User = strings.ToLower(u.Nick)
		}
		if len(u.Real) < 1 {
			u.Real = u.Nick
		}
	}

	// Update details for each command's settings.
	for i := range config.Commands {
		cs := &config.Commands[i].Settings
		s := &config.Settings
		cs.CaseSensitive = s.CaseSensitive
		cs.Symbol = s.Symbol
		cs.Scope = make([]string, 2)
		copy(cs.Scope, s.Scope)
		cs.Admin = s.Admin
	}
}

// getServer searches a config struct for a server with a specified id. It
// returns the server if found, or an error otherwise.
func getServer(config *Config, id string) (*ircutil.Server, error) {
	for i := range config.Servers {
		if id == config.Servers[i].ID {
			return &config.Servers[i], nil
		}
	}
	return nil, errors.New("getting server: id not found")
}

// getUser searches a config struct for a user with a specified id. It returns
// the user if found, or an error otherwise.
func getUser(config *Config, id string) (*ircutil.User, error) {
	for i := range config.Users {
		if id == config.Users[i].ID {
			return &config.Users[i], nil
		}
	}
	return nil, errors.New("getting server: id not found")
}
