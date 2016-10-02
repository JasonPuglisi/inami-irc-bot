package main

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

// Config stores user-defined configuration values loaded from a separate file.
type Config struct {
	// (Optional) Default command settings. Can be overrode for individual
	// commands. Default: All nested defaults
	Settings `json:"settings"`
	// List of servers that can be used in a client.
	Servers []Server `json:"servers"`
	// List of users that can be used in a client.
	Users []User `json:"users"`
	// List of clients. Connections between user and server.
	Clients []Client `json:"clients"`
	// List of commands to be used for all clients. To use different commands
	// with different clients, run multiple instances of the program with
	// different configuration files.
	Commands []Command `json:"commands"`
}

// Settings stores command settings that can be used as defaults or overrode on
// a per-command basis.
type Settings struct {
	// (Optional) Whether or not a command is case-sensitive. Default: true
	CaseSensitive bool `json:"caseSensitive"`
	// (Optional) Symbol/string that must be prefixed to a command trigger for
	// it to be detected. Default: "!"
	Symbol string `json:"symbol"`
	// (Optional) Array of places a command can be used. Contains "channel" for
	// in a channel, and/or "direct" for in a private message to the client.
	// Default: ["channel"]
	Scope []string `json:"scope"`
	// (Optional) Whether or not a command can only be used by client admins.
	// Default: false
	Admin bool `json:"admin"`
}

// Server stores server connection settings.
type Server struct {
	// String ID of the server. Used to identify the server in a client.
	ID string `json:"id"`
	// Hostname of server. Can be an IP address or domain name.
	Host string `json:"host"`
	// (Optional) Port of server. Default: 6697
	Port int `json:"port"`
	// (Optional) Whether or not the server port should be connected to using
	// SSL/TLS. Default: true if port is 6697, false otherwise
	Secure bool `json:"secure"`
}

// User stores user settings.
type User struct {
	// String ID of user. Used to identify the user in a client.
	ID string `json:"id"`
	// Nickname of client.
	Nick string `json:"nick"`
	// (Optional) Username of client. Default: Same as nickname lowercased
	User string `json:"user"`
	// (Optional) Realname of client. Default: Same as nickname
	Real string `json:"real"`
}

// Client stores client/connection settings and credentials.
type Client struct {
	// ID of server to use in connection.
	Server string `json:"server"`
	// ID of user to use in connection.
	User string `json:"user"`
	// (Optional) List of channels to join upon connection and authentication
	// (if specified). Channels must be prefixed with "#" (e.g., "#channel").
	// Channels with a password must have a space between the channel name and
	// password (e.g., "#channel pass"). Default: []
	Channels []string `json:"channels"`
	// (Optional) String of user modes to be set upon connection and
	// authentication (if specified). Must have a "+" before all modes to be
	// set and a "-" before all modes to be unset (e.g., "+i-x"). Default: ""
	Modes string `json:"modes"`
	// (Optional) List of client admin nicknames able to run commands set to
	// admin-only. Default: []
	Admins []string `json:"admins"`
	// (Optional) Authentication credentials to used in connection.
	// Defaults: All nested defaults
	Authentication `json:"authentication"`
}

// Authentication stores authentication credentials for servers and nicknames.
type Authentication struct {
	// (Optional) Server password to connect to a server. Empty string for
	// none. Default: ""
	ServerPassword string `json:"serverPassword"`
	// (Optional) Nickserv password to identify user with nickserv. Empty
	// string for none. Default: ""
	Nickserv string `json:"nickserv"`
}

// Command stores command triggers, execution details, and settings, with the
// ability to override default settings.
type Command struct {
	// List of strings that will trigger the command. Triggers must not contain
	// the command symbol, as it will be checked for automatically.
	Triggers []string `json:"triggers"`
	// Function that will be executed when the command is triggered.
	Function string `json:"function"`
	// (Optional) String of arguments that must follow the command. Optional
	// arguments must have square brackets around them (e.g., "arg1 [arg2]").
	// Default: ""
	Arguments string `json:"arguments"`
	// (Optional) Command settings to override default command settings. See
	// default command settings (above) for descriptions of each setting.
	// Defaults: Default command settings
	Settings `json:"settings"`
}

// getConfig opens a config file at the given path and parses it into a config
// struct with default values applied.
func getConfig(path string) (*Config, error) {
	// Attempt to open configuration file.
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return &Config{}, err
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
		Settings: Settings{
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
		if u.User == "" {
			u.User = strings.ToLower(u.Nick)
		}
		if u.Real == "" {
			u.Real = u.Nick
		}
	}

	// Update details for each command's settings.
	for i := range config.Commands {
		cs := &config.Commands[i].Settings
		s := &config.Settings
		cs.CaseSensitive = s.CaseSensitive
		cs.Symbol = s.Symbol
		copy(cs.Scope, s.Scope)
		cs.Admin = s.Admin
	}
}
