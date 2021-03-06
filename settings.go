// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
)

const (
	defaultBinding    = ":7070"
	defaultConfigPath = "deploy.conf.js"
)

type PostDeploySettings struct {
	Verbose bool
	Binding string
	Config  string
}

var Settings PostDeploySettings = PostDeploySettings{}

func init() {
	flag.BoolVar(&Settings.Verbose, "v", false, "Verbose mode")
	flag.StringVar(&Settings.Binding, "binding", defaultBinding, "The http binding")
	flag.StringVar(&Settings.Config, "config", defaultConfigPath, "The deployment configuration")
}
