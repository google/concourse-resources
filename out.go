// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"

	"golang.org/x/build/gerrit"
)

type outRequest struct {
	Source  `json:"source"`
	Params  outParams `json:"params"`
}

type outParams struct {
	Repository  string         `json:"repository"`
	Message     string         `json:"message"`
	MessageFile string         `json:"message_file"`
	Labels      map[string]int `json:"labels"`
}

func outMain(reqDecoder *json.Decoder, buildDir string) ResourceResponse {
	var req outRequest
	err := reqDecoder.Decode(&req)
	fatalErr(err, "error reading request")

	authMan := newAuthManager(req.Source)
	defer authMan.cleanup()

	var ver Version
	if req.Params.Repository == "" {
		log.Fatalln("param repository required")
	}
	gerritVersionPath := filepath.Join(
		buildDir, req.Params.Repository, gerritVersionFilename)
	err = ver.ReadFromFile(gerritVersionPath)
	fatalErr(err, "error reading %q", gerritVersionPath)

	message := req.Params.Message

	if messageFile := req.Params.MessageFile; messageFile != "" {
		message_bytes, err := ioutil.ReadFile(filepath.Join(buildDir, messageFile))
		if err == nil {
			message = string(message_bytes)
		} else {
			log.Printf("error reading message file %q: %v", messageFile, err)
			if message == "" {
				log.Fatalln("no fallback message; failing")
			} else {
				log.Printf("using fallback message %q", message)
			}
		}
	}

	c, err := gerritClient(req.Source, authMan)
	fatalErr(err, "error setting up gerrit client")

	ctx := context.Background()

	err = c.SetReview(ctx, ver.ChangeId, ver.Revision, gerrit.ReviewInput{
		Message: message,
		Labels:  req.Params.Labels,
	})
	fatalErr(err, "error sending review")

	return ResourceResponse{Version: ver}
}