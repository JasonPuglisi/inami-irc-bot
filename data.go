package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/jasonpuglisi/ircutil"
)

// getData opens a data file at the given path and parses it into a data
// struct.
func getData(path string) (ircutil.Data, error) {
	// Attempt to open data file.
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Parse file data into data struct.
	data := &ircutil.Data{}
	err = json.Unmarshal(raw, data)
	if err != nil {
		return nil, err
	}

	// Return parsed data.
	return *data, nil
}

// writeData writes a data struct to a data file at the given path.
func writeData(client *ircutil.Client) error {
	// Parse data structure to string.
	raw, err := json.Marshal(client.Data)
	if err != nil {
		return err
	}

	// Write string to data file.
	err = ioutil.WriteFile(*client.DataFile, []byte(raw), 0644)
	return err
}

// GetValue gets a value from persistent data using a keys array in the format
// [ "scope", "owner", "data_group", "key" ]. Scope must be "user", "channel",
// or "client".
func GetValue(client *ircutil.Client, keys []string) (string, error) {
	// Error if number of key parameters is wrong.
	if len(keys) != 4 {
		return "", errors.New("getting data: invalid number of key parameters")
	}

	// Set individual key parameters.
	scope, owner, group, key := keys[0], keys[1], keys[2], keys[3]

	// Error if scope is invalid.
	if scope != "user" && scope != "channel" && scope != "client" {
		return "", errors.New("getting data: invalid scope")
	}

	// Return value.
	buildMap(client, keys)
	return client.Data[scope][owner][group][key], nil
}

// SetValue sets a value in persistent data using a keys array in the same
// format and with the same scopes as GetValue.
func SetValue(client *ircutil.Client, keys []string, value string) error {
	// Error if number of key parameters is wrong.
	if len(keys) != 4 {
		return errors.New("setting data: invalid number of key parameters")
	}

	// Set individual key parameters.
	scope, owner, group, key := keys[0], keys[1], keys[2], keys[3]

	// Error if scope is invalid.
	if scope != "user" && scope != "channel" && scope != "client" {
		return errors.New("setting data: invalid scope")
	}

	// Set value and write data file.
	buildMap(client, keys)
	client.Data[scope][owner][group][key] = value
	return writeData(client)
}

// buildMap ensures all levels of a map exist, and creates them if necessary.
// It uses a keys array in the same format and with the same scopes as
// GetValue.
func buildMap(client *ircutil.Client, keys []string) {
	// Set data and key values.
	data, scope, owner, group := client.Data, keys[0], keys[1], keys[2]

	// Build each level of the map.
	if data[scope] == nil {
		data[scope] = map[string]map[string]map[string]string{}
	}
	if data[scope][owner] == nil {
		data[scope][owner] = map[string]map[string]string{}
	}
	if data[scope][owner][group] == nil {
		data[scope][owner][group] = map[string]string{}
	}
}
