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
	"fmt"
	"log"
	"net/url"
	"os/exec"
	"path"
	"path/filepath"

	"golang.org/x/build/gerrit"
)

const (
	gerritVersionFilename = ".gerrit_version.json"
)

var (
	defaultFetchProtocols = []string{"http", "anonymous http"}
	execGit               = realExecGit // For testing
)

type inRequest struct {
	Source  `json:"source"`
	Version `json:"version"`
	Params  inParams `json:"params"`
}

type inParams struct {
	FetchProtocol string `json:"fetch_protocol"`
	FetchUrl      string `json:"fetch_url"`
}

func inMain(reqDecoder *json.Decoder, destDir string) ResourceResponse {
	var req inRequest
	err := reqDecoder.Decode(&req)
	fatalErr(err, "error reading request")

	authMan := newAuthManager(req.Source)
	defer authMan.cleanup()

	c, err := gerritClient(req.Source, authMan)
	fatalErr(err, "error setting up gerrit client")

	ctx := context.Background()

	change, rev, err := getVersionChangeRevision(c, ctx, req.Version)
	fatalErr(err, "")

	fetchArgs, err := resolveFetchArgs(
		req.Params, rev)
	fatalErr(err, "could not resolve fetch args for change %q", change.ID)

	gitFatalErr(destDir, "init")
	gitFatalErr(destDir, "config", "color.ui", "always")

	configArgs, err := authMan.gitConfigArgs()
	fatalErr(err, "error getting git config args")
	gitFatalErr(destDir, configArgs...)

	gitFatalErr(destDir, fetchArgs...)
	gitFatalErr(destDir, "checkout", "FETCH_HEAD")

	metadata := make(metadataMap)
	metadata["project"] = change.Project
	metadata["subject"] = change.Subject
	if rev.Uploader != nil {
		metadata["uploader"] = fmt.Sprintf("%s <%s>", rev.Uploader.Name, rev.Uploader.Email)
	}
	link, err := buildRevisionLink(req.Source, change.ChangeNumber, rev.PatchSetNumber)
	if err == nil {
		metadata["link"] = link
	} else {
		log.Printf("error building revision link: %v", err)
	}

	gerritVersionPath := filepath.Join(destDir, gerritVersionFilename)
	err = req.Version.WriteToFile(gerritVersionPath)
	fatalErr(err, "error writing %q", gerritVersionPath)

	return ResourceResponse{Version: req.Version, Metadata: metadata}
}

func resolveFetchArgs(params inParams, rev *gerrit.RevisionInfo) ([]string, error) {
	fetchUrl := params.FetchUrl
	fetchRef := rev.Ref
	if fetchUrl == "" {
		fetchProtocol := params.FetchProtocol
		if fetchProtocol == "" {
			for _, proto := range defaultFetchProtocols {
				if _, ok := rev.Fetch[proto]; ok {
					fetchProtocol = proto
					break
				}
			}
		}
		fetchInfo, ok := rev.Fetch[fetchProtocol]
		if ok {
			fetchUrl = fetchInfo.URL
			fetchRef = fetchInfo.Ref
		} else {
			return []string{}, fmt.Errorf("no fetch info for protocol %q", fetchProtocol)
		}
	}
	return []string{"fetch", fetchUrl, fetchRef}, nil
}

func gitFatalErr(dir string, args ...string) {
	gitArgs := append([]string{"-C", dir}, args...)
	log.Printf("git %v", gitArgs)
	output, err := execGit(gitArgs...)
	log.Printf("git output:\n%s", output)
	fatalErr(err, "git failed")
}

func realExecGit(args ...string) ([]byte, error) {
	return exec.Command("git", args...).CombinedOutput()
}

func buildRevisionLink(src Source, changeNum int, psNum int) (string, error) {
	srcUrl, err := url.Parse(src.Url)
	if err != nil {
		return "", err
	}
	srcUrl.Path = path.Join(srcUrl.Path, fmt.Sprintf("c/%d/%d", changeNum, psNum))
	return srcUrl.String(), nil
}
