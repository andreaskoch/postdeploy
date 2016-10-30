// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"os"
)

// The configProvider interface returns Config models.
type configProvider interface {
	// GetConfig returns a Config model or an error.
	GetConfig() (Config, error)
}

// jsonConfigProvider returns Config models from JSON files.
type jsonConfigProvider struct {
	path string
}

// GetConfig returns a Config model or an error.
func (provider jsonConfigProvider) GetConfig() (Config, error) {

	// open the config file
	file, err := os.Open(provider.path)
	if err != nil {
		return Config{}, err
	}

	// deserialize the config
	decoder := json.NewDecoder(file)
	var config Config
	if err := decoder.Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}

type Config struct {
	Hooks []DeploymentHook `json:"hooks"`
}

type DeploymentHook struct {
	Provider  string    `json:"provider"`
	Route     string    `json:"route"`
	Directory string    `json:"directory"`
	Commands  []Command `json:"commands"`
}

type Command struct {
	Name string   `json:"name"`
	Args []string `json:"args"`
}
