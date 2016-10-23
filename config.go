// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"io"
	"os"
)

func getConfig(path string) *Config {

	// open the config file
	file, err := os.Open(path)
	if err != nil {
		return nil
	}

	// deserialize the config
	deserializer := NewConfigDeserializer()
	config, err := deserializer.Deserialize(file)
	if err != nil {
		return nil
	}

	return config
}

type Config struct {
	Hooks []*DeploymentHook `json:"hooks"`
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

type ConfigDeserializer struct{}

func NewConfigDeserializer() *ConfigDeserializer {
	return &ConfigDeserializer{}
}

func (ConfigDeserializer) Deserialize(reader io.Reader) (*Config, error) {
	decoder := json.NewDecoder(reader)
	var config *Config
	err := decoder.Decode(&config)
	return config, err
}
