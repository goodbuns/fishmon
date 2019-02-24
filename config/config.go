// Package config provides configuration file parsing for mapping temperature
// probes to Adafruit.IO feeds.
package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"

	"github.com/liftM/fishmon/pkg/ds18b20"
)

// Configuration errors.
var (
	ErrNoSuchProbe = errors.New("probe configuration not found")
)

// File stores the contents of a configuration file.
type File struct {
	Version string
	Probes  map[ds18b20.ID]Probe
}

// Probe stores the configuration for a single temperature probe.
type Probe struct {
	Name    string
	FeedKey string `json:"feed"`
}

// New parses a configuration file.
func New(filename string) (*File, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "could not read fishmon configuration file")
	}
	var file File
	err = json.Unmarshal(bytes, &file)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal fishmon configuration file")
	}

	return &file, nil
}
