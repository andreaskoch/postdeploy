// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const (
	VERSION = "0.1.0"
)

var (
	deploymentHookPattern = regexp.MustCompile(`/deploy/([^/]+)/([^/]+)`)

	config *Config
)

func main() {

	// print application info
	message("%s (Version: %s)\n\n", os.Args[0], VERSION)

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
	message("Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func deploymentHookHandler(w http.ResponseWriter, r *http.Request) {

	// parse the request url
	requestUri := r.RequestURI
	message("Request URI: %s", requestUri)

	// check the deploy hook
	isMatch, matches := isMatch(requestUri, deploymentHookPattern)
	if !isMatch || len(matches) < 2 {
		error404Handler(w, r)
		return
	}

	// detect the provider
	provider := matches[1]
	message("Provider: %s", provider)

	// deted the route
	route := matches[2]
	message("Route: %s", route)

	// read the post body
	p, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// handler error
	}
	postBody := string(p)

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
		bitbucket(w, r, postBody, theHook.Directory, theHook.Command)
	default:
		generic(theHook.Directory, theHook.Command)
	}
}

func error500Handler(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, err)
}

func error404Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func execute(directory, commandText string) {

	message("Executing command %q in directory %q", commandText, directory)

	// get the command
	command := getCmd(directory, commandText)

	// execute the command
	if err := command.Start(); err != nil {
		fmt.Println(err)
	}

	// wait for the command to finish
	command.Wait()

	fmt.Println()
}

func getCmd(directory, commandText string) *exec.Cmd {
	if commandText == "" {
		return nil
	}

	components := strings.Split(commandText, " ")

	// get the command name
	commandName := components[0]

	// get the command arguments
	arguments := make([]string, 0)
	if len(components) > 1 {
		arguments = components[1:]
	}

	// create the command
	command := exec.Command(commandName, arguments...)

	// set the working directory
	command.Dir = directory

	// redirect command io
	redirectCommandIO(command)

	return command
}

func redirectCommandIO(cmd *exec.Cmd) (*os.File, error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println(err)
	}

	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)

	//direct. Masked passwords work OK!
	cmd.Stdin = os.Stdin
	return nil, err
}

func message(text string, args ...interface{}) {

	// append newline character
	if !strings.HasSuffix(text, "\n") {
		text += "\n"
	}

	fmt.Printf(text, args...)
}

// isMatch returns a flag indicating whether the supplied
// text and pattern do match and if yet, the matched text.
func isMatch(text string, pattern *regexp.Regexp) (isMatch bool, matches []string) {
	matches = pattern.FindStringSubmatch(text)
	return matches != nil, matches
}
