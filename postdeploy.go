// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	defaultBinding    = ":7070"
	defaultConfigPath = "deploy.conf.js"

	applicationName = "postdeploy"
	version         = "v2.1.0"
)

func main() {

	app := kingpin.New(applicationName, "A http service that listens for deployment requests and executes a set of predefined commands when the request arrives")
	app.Version(version)
	bindAddressParameter := app.Flag("binding", "The port and address you want to listen on").Short('b').Default(defaultBinding).OverrideDefaultFromEnvar("POSTDEPLOY_BINDING").String()
	configFileParameter := app.Flag("configfile", "The deployment configuration").Short('c').Default(defaultConfigPath).OverrideDefaultFromEnvar("POSTDEPLOY_CONFIGFILE").String()

	kingpin.MustParse(app.Parse(os.Args[1:]))

	if err := postDeploy(*bindAddressParameter, *configFileParameter); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func postDeploy(bindAddress, configFilePath string) error {

	configProvider := jsonConfigProvider{configFilePath}
	routeBuilder := deploymentRouteBuilder{}
	hookProvider := configBasedDeploymentHookProvider{configProvider, routeBuilder}
	server := deploymentServer{bindAddress, hookProvider}

	return server.Run()
}

type deploymentServer struct {
	bindAddress  string
	hookProvider deploymentHookProvider
}

func (server deploymentServer) Run() error {

	hooks, err := server.hookProvider.GetHooks()
	if err != nil {
		return err
	}

	for _, hook := range hooks {
		http.HandleFunc(hook.Route, hook.Handler)
	}

	if err := http.ListenAndServe(server.bindAddress, nil); err != nil {
		return err
	}

	return nil
}

type deploymentHandler struct {
	Route   string
	Handler http.HandlerFunc
}

type deploymentHookProvider interface {
	GetHooks() ([]deploymentHandler, error)
}

type configBasedDeploymentHookProvider struct {
	deploymentConfigProvider configProvider
	deploymentRouteBuilder   routeBuilder
}

func (hookProvider configBasedDeploymentHookProvider) GetHooks() ([]deploymentHandler, error) {

	config, err := hookProvider.deploymentConfigProvider.GetConfig()
	if err != nil {
		return nil, err
	}

	var handlers []deploymentHandler
	for _, hook := range config.Hooks {

		switch hook.Provider {

		case "bitbucket":

			handlers = append(handlers, deploymentHandler{
				Route: hookProvider.deploymentRouteBuilder.GetRoute(hook.Provider, hook.Route),
				Handler: func(w http.ResponseWriter, r *http.Request) {
					bitbucket(w, r, hook.Directory, hook.Commands)
				},
			})

		default:

			handlers = append(handlers, deploymentHandler{
				Route: hookProvider.deploymentRouteBuilder.GetRoute(hook.Provider, hook.Route),
				Handler: func(w http.ResponseWriter, r *http.Request) {
					generic(hook.Directory, hook.Commands)
				},
			})

		}

	}

	return handlers, nil
}

type routeBuilder interface {
	GetRoute(providerName, routeName string) string
}

type deploymentRouteBuilder struct {
}

func (deploymentRouteBuilder) GetRoute(providerName, routeName string) string {
	return fmt.Sprintf("/deploy/%s/%s", providerName, routeName)
}

func error500Handler(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, err)
}

func error404Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func execute(directory string, commands []Command) {

	for _, command := range commands {
		runCommand(os.Stdout, os.Stderr, directory, command)
	}

}

// Execute go in the specified go path with the supplied command arguments.
func runCommand(stdout, stderr io.Writer, workingDirectory string, command Command) {

	expandedWorkingDirectory := os.ExpandEnv(workingDirectory)
	expandedCommandName := os.ExpandEnv(command.Name)
	var expandedArguments []string
	for _, argument := range command.Args {
		expandedArguments = append(expandedArguments, os.ExpandEnv(argument))
	}

	cmd := exec.Command(expandedCommandName, expandedArguments...)

	cmd.Dir = expandedWorkingDirectory
	cmd.Env = os.Environ()

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	log.Printf("%s: %s %s", expandedWorkingDirectory, command.Name, strings.Join(command.Args, " "))

	err := cmd.Run()
	if err != nil {
		log.Printf("Error running %s: %v", command, err)
	}
}
