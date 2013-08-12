// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"os"
)

const (
	defaultBinding   = ":7070"
	defaultCommand   = ""
	defaultDirectory = ""
)

type PostDeploySettings struct {
	Binding   string
	Directory string
	Command   string
}

var Settings PostDeploySettings = PostDeploySettings{}

func init() {

	// use the current directory as the default path
	defaultDirectory, err := os.Getwd()
	if err != nil {
		defaultDirectory = "."
	}

	flag.StringVar(&Settings.Binding, "binding", defaultBinding, "The http binding")
	flag.StringVar(&Settings.Directory, "directory", defaultDirectory, "The working directory")
	flag.StringVar(&Settings.Command, "command", defaultCommand, "The command that will be executed when the deployment hook is triggered")
}
