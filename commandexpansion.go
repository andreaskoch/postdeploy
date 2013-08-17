// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	PARAMETER_REPOSITORYNAME = "repository"
)

var (
	// Command values can only contain [a-z, 0-9, -, _]
	IsValidCommandValue = regexp.MustCompile(`^[a-z0-9_-]+$`)
)

func newParameterList() *ParameterList {
	return &ParameterList{
		Values: make([]*CommandParameter, 0),
	}
}

func newCommandParameter(name, value string) *CommandParameter {
	return &CommandParameter{
		Name:  name,
		Value: value,
	}
}

type ParameterList struct {
	Values []*CommandParameter
}

func (parameters *ParameterList) Add(name, value string) error {

	// Check if the command value is valid.
	// Until i found a better way to guard against command injection
	// i will only allow alpha-numeric characters
	if !IsValidCommandValue.MatchString(value) {
		return fmt.Errorf("%q is an invalid command value.", value)
	}

	parameters.Values = append(parameters.Values, newCommandParameter(name, value))
	return nil
}

func (parameters ParameterList) Expand(command string) string {
	for _, parameter := range parameters.Values {
		// replace placeholders auch as "{Repository}"" with the corresponding value
		placeholder := fmt.Sprintf("{%s}", parameter.Name)
		command = strings.Replace(command, placeholder, parameter.Value, -1)
	}

	return command
}

type CommandParameter struct {
	Name  string
	Value string
}
