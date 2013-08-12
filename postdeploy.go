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
	"strings"
)

const (
	VERSION = "0.1.0"
)

var usage = func() {
	message("Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

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

	http.HandleFunc("/", deploymentHook)

	// start the server
	if err := http.ListenAndServe(Settings.Binding, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}

func deploymentHook(w http.ResponseWriter, r *http.Request) {

	// parse the request url
	deployRoute := "/deploy"
	requestUri := r.RequestURI

	fmt.Println(requestUri)

	// check the deploy hook
	if !strings.HasPrefix(requestUri, deployRoute) {
		error404Handler(w, r)
		return
	}

	// detect the provider
	provider := strings.Replace(requestUri, deployRoute, "", 1)
	provider = strings.TrimPrefix(provider, "/")

	fmt.Println(provider)

	// read the post body
	p, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// handler error
	}
	postBody := string(p)

	switch provider {
	case "bitbucket":
		bitbucket(w, r, postBody)
	default:
		error404Handler(w, r)
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
