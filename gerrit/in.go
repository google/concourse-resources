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
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"golang.org/x/build/gerrit"

	"github.com/google/concourse-resources/internal"
)

const (
	gerritVersionFilename = ".gerrit_version.json"
)

var (
	defaultFetchProtocols = []string{"http", "anonymous http"}

	// For testing
	execGit               = realExecGit
)

type inParams struct {
	FetchProtocol string `json:"fetch_protocol"`
	FetchUrl      string `json:"fetch_url"`
}

func init() {
	internal.RegisterInFunc(in)
}

func in(rs *internal.ResourceContext, src Source, ver Version, params inParams) Version {
	authMan := newAuthManager(src)
	defer authMan.cleanup()

	c, err := gerritClient(src, authMan)
	fatalErr(err, "error setting up gerrit client")

	ctx := context.Background()

	// Fetch requested version from Gerrit
	change, rev, err := getVersionChangeRevision(c, ctx, ver)
	fatalErr(err, "")

	fetchArgs, err := resolveFetchArgs(params, rev)
	fatalErr(err, "could not resolve fetch args for change %q", change.ID)

	// Prepare destination repo and checkout requested revision
	gitFatalErr(rs.TargetDir, "init")
	gitFatalErr(rs.TargetDir, "config", "color.ui", "always")

	configArgs, err := authMan.gitConfigArgs()
	fatalErr(err, "error getting git config args")
	gitFatalErr(rs.TargetDir, configArgs...)

	gitFatalErr(rs.TargetDir, fetchArgs...)
	gitFatalErr(rs.TargetDir, "checkout", "FETCH_HEAD")

	// Build response metadata
	rs.AddMetadata("project", change.Project)
	rs.AddMetadata("subject", change.Subject)
	if rev.Uploader != nil {
		rs.AddMetadata("uploader", fmt.Sprintf("%s <%s>", rev.Uploader.Name, rev.Uploader.Email))
	}
	link, err := buildRevisionLink(src, change.ChangeNumber, rev.PatchSetNumber)
	if err == nil {
		rs.AddMetadata("link", link)
	} else {
		log.Printf("error building revision link: %v", err)
	}

	// Write gerrit_version.json
	gerritVersionPath := filepath.Join(rs.TargetDir, gerritVersionFilename)
	err = ver.WriteToFile(gerritVersionPath)
	fatalErr(err, "error writing %q", gerritVersionPath)

	// Ignore gerrit_version.json file in repo
	excludePath := filepath.Join(rs.TargetDir, ".git", "info", "exclude")
	err = os.MkdirAll(filepath.Dir(excludePath), 0755)
	if err == nil {
		f, err := os.OpenFile(excludePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err == nil {
			defer f.Close()
			_, err = fmt.Fprintf(f, "\n/%s\n", gerritVersionFilename)
		}
	}
	if err != nil {
		log.Printf("error adding %q to %q: %v", gerritVersionPath, excludePath, err)
	}

	return ver
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
