// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var (
	serializer *JSONSerializer
)

func init() {
	serializer = NewJSONSerializer()
}

func bitbucket(w http.ResponseWriter, r *http.Request, directory, command string) {

	// extract the payload from the request
	requestBody := r.PostFormValue("payload")

	// unescape the request body
	unescapedRequestBody, unescapeError := url.QueryUnescape(requestBody)
	if unescapeError != nil {
		message("Unable to unescape request body. Error: %s", unescapeError)
		error500Handler(w, r, unescapeError)
		return
	}

	// deserialize request
	postMessage, deserializeError := serializer.Deserialize(strings.NewReader(unescapedRequestBody))
	if deserializeError != nil {
		message("Unable to deserialize %s. Error: %s", unescapedRequestBody, deserializeError)
		error500Handler(w, r, deserializeError)
		return
	}

	// get the command parameters from the post message
	commandParameters := getParameterListFromBitbucketPost(postMessage)

	// expand the command parameters
	expandedCommand := commandParameters.Expand(command)

	// execute comand
	go execute(directory, expandedCommand)
}

func getParameterListFromBitbucketPost(postMessage *BitbucketPostMessage) *ParameterList {
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

type JSONSerializer struct{}

func NewJSONSerializer() *JSONSerializer {
	return &JSONSerializer{}
}

func (JSONSerializer) Serialize(writer io.Writer, deploymentRequest *BitbucketPostMessage) error {
	bytes, err := json.MarshalIndent(deploymentRequest, "", "\t")
	if err != nil {
		return err
	}

	writer.Write(bytes)
	return nil
}

func (JSONSerializer) Deserialize(reader io.Reader) (*BitbucketPostMessage, error) {
	decoder := json.NewDecoder(reader)
	var deploymentRequest *BitbucketPostMessage
	err := decoder.Decode(&deploymentRequest)
	return deploymentRequest, err
}

type BitbucketSerializer interface {
	Serialize(writer io.Writer, deploymentRequest *BitbucketPostMessage) error
}

type BitbucketDeserializer interface {
	Deserialize(reader io.Reader) (*BitbucketPostMessage, error)
}
