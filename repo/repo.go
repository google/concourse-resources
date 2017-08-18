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
	"log"
	"strings"

	"github.com/google/concourse-resources/repo/internal"
)

var (
	testingRepo = false
	testCurrentVersion Version

	testLastRepoInitArgs []string
	testLastRepoSyncArgs []string
)

func repoInit(repoDir string, src Source) error {
	if src.ManifestUrl == "" {
		return errors.New("manifest_url is required")
	}

	opts := options{
		"manifest-url": src.ManifestUrl,
		"depth": "1",
	}

	if src.ManifestName != "" {
		opts["manifest-name"] = src.ManifestName
	}

	if src.ManifestBranch != "" {
		opts["manifest-branch"] = src.ManifestBranch
	}

	if len(src.Groups) > 0 {
		opts["groups"] = strings.Join(src.Groups, ",")
	}

	opts.merge(src.InitOptions)

	args := opts.args()

	if testingRepo {
		testLastRepoInitArgs = args
		return nil
	}

	log.Printf("repo init %v", args)
	output, err := internal.RepoInit(repoDir, args...)
	log.Printf("repo init stdout:\n%s", output)
	internal.LogExecErrors("repo init", err)

	return err
}

func repoSync(repoDir string, src Source) error {
	opts := options{
		"current-branch": true,
		"no-tags": true,
		"optimized-fetch": true,
	}
	opts.merge(src.SyncOptions)

	args := append([]string{"sync"}, opts.args()...)

	if testingRepo {
		testLastRepoSyncArgs = args
		return nil
	}

	log.Printf("repo %v", args)
	output, err := internal.RepoRun(repoDir, args...)
	log.Printf("repo sync stdout:\n%s", output)
	internal.LogExecErrors("repo sync", err)

	return err
}

func getCurrentVersion(repoDir string) (ver Version, err error) {
	if testingRepo {
		return testCurrentVersion, nil
	}
	log.Println("repo manifest -r")
	output, err := internal.RepoRun(repoDir, "manifest", "--revision-as-HEAD")
	if !internal.LogExecErrors("repo manifest", err) {
		ver.Manifest = string(output)
	}
	return
}
