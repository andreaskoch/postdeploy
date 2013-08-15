// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strings"
)

const (
	PARAMETER_REPOSITORYNAME = "repository"
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

func (parameters *ParameterList) Add(name, value string) {
	parameters.Values = append(parameters.Values, newCommandParameter(name, value))
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
