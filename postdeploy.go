// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var (
	deploymentHookPattern = regexp.MustCompile(`/deploy/([^/]+)/([^/]+)`)
	applicationName       = "postdeploy"
	version               = "v2.1.0"

	config *Config
)

func main() {

	// print application info
	message("%s (Version: %s)\n\n", applicationName, version)

	// print usage information if no arguments are supplied
	if len(os.Args) == 1 {
		usage()
		os.Exit(1)
	}

	// parse the flags
	flag.Parse()

	// get the config
	config = getConfig(Settings.Config)
	if config == nil {
		message("Unable to load config from %q", Settings.Config)
		os.Exit(2)
	}

	// attach the deployment handler
	http.HandleFunc("/", deploymentHookHandler)

	// start the server
	if err := http.ListenAndServe(Settings.Binding, nil); err != nil {
		message("%s", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func usage() {
	message("%s is a http service that listens for deployment requests and executes a predefined command when the request arrives", applicationName)
	message("")
	message("Usage:\n")
	message("  %s -binding \":7070\" -config \"postdeploy.conf.js\"", applicationName)
	message("")
	message("Parameters:\n")
	flag.PrintDefaults()
}

func deploymentHookHandler(w http.ResponseWriter, r *http.Request) {

	// parse the request url
	requestUri := r.RequestURI

	if Settings.Verbose {
		message("Request URI: %s", requestUri)
	}

	// check the deploy hook
	isMatch, matches := isMatch(requestUri, deploymentHookPattern)
	if !isMatch || len(matches) < 2 {
		error404Handler(w, r)
		return
	}

	// detect the provider
	provider := matches[1]
	message("Provider: %s", provider)

	// detect the route
	route := matches[2]
	message("Route: %s", route)

	// find a matching hook
	var theHook *DeploymentHook
	for _, hook := range config.Hooks {
		if hook.Provider == provider && hook.Route == route {
			theHook = hook
			break
		}
	}

	if theHook == nil {
		message("No matching hook for provider %q and route %q", provider, route)
		error404Handler(w, r)
		return
	}

	// execute the handler
	switch theHook.Provider {
	case "bitbucket":
		bitbucket(w, r, theHook.Directory, theHook.Commands)
	default:
		generic(theHook.Directory, theHook.Commands)
	}
}

func error500Handler(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, err)
}

func error404Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func execute(directory string, commands []Command) {

	if directory == "" {
		directory = getWorkingDirectory()
	}

	for _, command := range commands {
		message("Executing command %q in directory %q", command.Name, directory)
		runCommand(os.Stdout, os.Stderr, directory, command)
	}

}

// Execute go in the specified go path with the supplied command arguments.
func runCommand(stdout, stderr io.Writer, workingDirectory string, command Command) {

	// set the go path
	cmd := exec.Command(command.Name, command.Args...)

	cmd.Dir = workingDirectory
	cmd.Env = os.Environ()

	// execute the command
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if Settings.Verbose {
		log.Printf("Running %s", command)
	}

	err := cmd.Run()
	if err != nil {
		log.Printf("Error running %s: %v", command, err)
	}
}

func message(text string, args ...interface{}) {

	// append newline character
	if !strings.HasSuffix(text, "\n") {
		text += "\n"
	}

	fmt.Printf(text, args...)
}

// getWorkingDirectory returns the current working directory path or fails.
func getWorkingDirectory() string {
	goPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}

	return goPath
}

// isMatch returns a flag indicating whether the supplied
// text and pattern do match and if yet, the matched text.
func isMatch(text string, pattern *regexp.Regexp) (isMatch bool, matches []string) {
	matches = pattern.FindStringSubmatch(text)
	return matches != nil, matches
}
