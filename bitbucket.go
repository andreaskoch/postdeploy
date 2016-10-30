// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func bitbucket(w http.ResponseWriter, r *http.Request, directory string, commands []Command) {

	// extract the payload from the request
	requestBody := r.PostFormValue("payload")

	// unescape the request body
	unescapedRequestBody, unescapeError := url.QueryUnescape(requestBody)
	if unescapeError != nil {
		log.Printf("Error: Unable to unescape request body. Error: %s", unescapeError)
		error500Handler(w, r, unescapeError)
		return
	}

	// deserialize request
	decoder := json.NewDecoder(strings.NewReader(unescapedRequestBody))
	var postMessage BitbucketPostMessage
	if err := decoder.Decode(&postMessage); err != nil {
		log.Printf("Error: Unable to deserialize %s. Error: %s", unescapedRequestBody, err)
		error500Handler(w, r, err)
		return
	}

	// get the command parameters from the post message
	commandParameters := getParameterListFromBitbucketPost(postMessage)

	// expand the command parameters
	expandedCommands := commandParameters.Expand(commands)

	// execute comand
	go execute(directory, expandedCommands)
}

func getParameterListFromBitbucketPost(postMessage BitbucketPostMessage) *ParameterList {
	parameterList := newParameterList()

	// add the repository name
	parameterList.Add(PARAMETER_REPOSITORYNAME, postMessage.Repository.Slug)

	return parameterList
}

type BitbucketPostMessage struct {
	CanonUrl   string              `json:"canon_url"`
	Commits    []BitbucketCommit   `json:"commits"`
	Repository BitbucketRepository `json:"repository"`
	User       string              `json:"user"`
}

type BitbucketCommit struct {
	Author  string `json:"author"`
	Branch  string `json:"branch"`
	Message string `json:"message"`
}

type BitbucketRepository struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Owner       string `json:"owner"`
	AbsoluteUrl string `json:"absolute_url"`
	Website     string `json:"website"`
}
