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
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/google/concourse-resources/repo/internal"
)

func initRepo(repoDir string, src Source) error {
	if src.ManifestUrl == "" {
		return errors.New("manifest_url is required")
	}

	initArgs := []string{
		"--depth", "1",
		"--manifest-url", src.ManifestUrl,
	}

	if src.ManifestName != "" {
		initArgs = append(initArgs,
			fmt.Sprintf("--manifest-name=%s", src.ManifestName))
	}

	if src.ManifestBranch != "" {
		initArgs = append(initArgs,
			fmt.Sprintf("--manifest-branch=%s", src.ManifestBranch))
	}

	if len(src.Groups) > 0 {
		initArgs = append(initArgs,
			fmt.Sprintf("--groups=%s", strings.Join(src.Groups, ",")))
	}

	if src.MirrorPath != "" {
		initArgs = append(initArgs,
			fmt.Sprintf("--reference=%s", src.MirrorPath))
	}

	log.Printf("repo init %v", initArgs)
	output, err := internal.RepoInit(repoDir, initArgs...)
	log.Printf("repo init stdout:\n%s", output)
	internal.LogExecErrors("repo init", err)

	return err
}

func syncRepo(repoDir string, src Source, extraArgs ...string) error {
	syncArgs := append([]string{
		"sync",
		"--current-branch",
		"--no-tags",
		"--optimized-fetch",
	}, extraArgs...)

	if src.SyncJobs > 0 {
		syncArgs = append(syncArgs,
			fmt.Sprintf("--jobs=%d", src.SyncJobs))
	}

	if !src.SyncVerbose {
		syncArgs = append(syncArgs, "--quiet")
	}

	if len(src.Projects) > 0 {
		syncArgs = append(syncArgs, src.Projects...)
	}

	log.Printf("repo %v", syncArgs)
	output, err := internal.RepoRun(repoDir, syncArgs...)
	log.Printf("repo sync stdout:\n%s", output)
	internal.LogExecErrors("repo sync", err)

	return err
}

func getCurrentVersion(repoDir string) (ver Version, err error) {
	log.Println("repo manifest -r")
	output, err := internal.RepoRun(repoDir, "manifest", "--revision-as-HEAD")
	if !internal.LogExecErrors("repo manifest", err) {
		ver.Manifest = string(output)
	}
	return
}

/*
func getProjectRevisions(src Source) (map[string]string, error) {
	forallArgs := []string{"forall", "--abort-on-errors"}

	if len(src.Projects) > 0 {
		forallArgs = append(forallArgs, src.Projects...)
	}

	forallArgs = append(forallArgs,
		"--command", `echo -e "${REPO_PROJECT}\x00$(git rev-parse HEAD)"`)

	log.Printf("repo %v", forallArgs)
	forallOut, err := internal.RepoRun(repoDir, forallArgs...)
	if internal.LogExecErrors("repo forall", err) {
		return nil, err
	}

	projRevs := make(map[string]string)
	for _, line := range bytes.Split(forallOut, []byte{'\n'}) {
		// Each line is expected to look like: "<project>\x00<revision>"
		nulIdx := bytes.LastIndexByte(line, 0)
		if nulIdx == -1 {
			return nil, fmt.Errorf("error processing project revisions; "+
				"got line with no NUL: %q", line)
		}
		proj := string(line[:nulIdx])
		rev := string(line[nulIdx+1:])
		projRevs[proj] = rev
	}

	return projRevs, nil
}
*/
