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
	"io/ioutil"
	"log"
	"path/filepath"

	"golang.org/x/build/gerrit"

	"github.com/google/concourse-resources/internal"
)

type outParams struct {
	Repository  string         `json:"repository"`
	Message     string         `json:"message"`
	MessageFile string         `json:"message_file"`
	Labels      map[string]int `json:"labels"`
}

func init() {
	internal.RegisterOutFunc(out)
}

func out(rs *internal.ResourceContext, src Source, params outParams) Version {
	authMan := newAuthManager(src)
	defer authMan.cleanup()

	// Read gerrit_version.json
	var ver Version
	if params.Repository == "" {
		log.Fatalln("param repository required")
	}
	gerritVersionPath := filepath.Join(
		rs.TargetDir, params.Repository, gerritVersionFilename)
	err := ver.ReadFromFile(gerritVersionPath)
	fatalErr(err, "error reading %q", gerritVersionPath)

	// Build comment message
	message := params.Message

	if messageFile := params.MessageFile; messageFile != "" {
		message_bytes, err := ioutil.ReadFile(filepath.Join(rs.TargetDir, messageFile))
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

	// Send review
	c, err := gerritClient(src, authMan)
	fatalErr(err, "error setting up gerrit client")

	ctx := context.Background()

	err = c.SetReview(ctx, ver.ChangeId, ver.Revision, gerrit.ReviewInput{
		Message: message,
		Labels:  params.Labels,
	})
	fatalErr(err, "error sending review")

	return ver
}
