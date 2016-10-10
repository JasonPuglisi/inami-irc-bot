package configutil

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/jasonpuglisi/ircutil"
)

// GetData opens a data file at the given path and parses it into a data
// struct.
func GetData(path string) (ircutil.Data, error) {
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
	clientPrefix, scope, owner, group, key := ircutil.GetClientPrefix(client),
		keys[0], keys[1], keys[2], keys[3]

	// Error if scope is invalid.
	if scope != "user" && scope != "channel" && scope != "client" {
		return "", errors.New("getting data: invalid scope")
	}

	// Return value.
	buildMap(client, keys)
	return client.Data[clientPrefix][scope][owner][group][key], nil
}

// SetValue sets a value in persistent data using a keys array in the same
// format and with the same scopes as GetValue.
func SetValue(client *ircutil.Client, keys []string, value string) error {
	// Error if number of key parameters is wrong.
	if len(keys) != 4 {
		return errors.New("setting data: invalid number of key parameters")
	}

	// Set individual key parameters.
	clientPrefix, scope, owner, group, key := ircutil.GetClientPrefix(client),
		keys[0], keys[1], keys[2], keys[3]

	// Error if scope is invalid.
	if scope != "user" && scope != "channel" && scope != "client" {
		return errors.New("setting data: invalid scope")
	}

	// Set value and write data file.
	buildMap(client, keys)
	client.Data[clientPrefix][scope][owner][group][key] = value
	return writeData(client)
}

// buildMap ensures all levels of a map exist, and creates them if necessary.
// It uses a keys array in the same format and with the same scopes as
// GetValue.
func buildMap(client *ircutil.Client, keys []string) {
	// Set data and key values.
	data, clientPrefix, scope, owner, group := client.Data,
		ircutil.GetClientPrefix(client), keys[0], keys[1], keys[2]

	// Build each level of the map.
	if data[clientPrefix] == nil {
		data[clientPrefix] = map[string]map[string]map[string]map[string]string{}
	}
	if data[clientPrefix][scope] == nil {
		data[clientPrefix][scope] = map[string]map[string]map[string]string{}
	}
	if data[clientPrefix][scope][owner] == nil {
		data[clientPrefix][scope][owner] = map[string]map[string]string{}
	}
	if data[clientPrefix][scope][owner][group] == nil {
		data[clientPrefix][scope][owner][group] = map[string]string{}
	}
}

// UpdateScope sets scope and owner appropriately in a keys array.
func UpdateScope(keys []string, source string, target string) {
	if ircutil.IsChannel(target) {
		keys[0], keys[1] = "channel", target
	} else {
		keys[0], keys[1] = "user", ircutil.GetNick(source)
	}
}
