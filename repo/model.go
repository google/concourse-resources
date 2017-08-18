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
	"os"
	"path/filepath"
)

var defaultMirrorDir = filepath.Join(os.TempDir(), "concourse-repo-mirror")

type Source struct {
	ManifestUrl    string `json:"manifest_url"`
	ManifestName   string `json:"manifest_name"`
	ManifestBranch string `json:"manifest_branch"`

	Groups []string `json:"groups"`
	Projects []string `json:"projects"`

	MirrorPath  string `json:"mirror_path"`
	SyncJobs    uint8  `json:"sync_jobs"`
	SyncVerbose bool   `json:"sync_verbose"`
}

type Version struct {
	Manifest string `json:"manifest"`
}

/*
// repo init -u <> --mirror --depth=1 --reference=<>
// repo sync --no-tags --current-branch
// repo forall -c 'echo -n "${REPO_PROJECT}:"; git rev-parse HEAD'
type Version struct {
	ProjectRevisions map[string]string
}

func (v Version) Equal(other Version) bool {
	for proj, ref := range v.ProjectRevisions {
		otherRef, hasProj := other.ProjectRevisions[proj]
		if !hasProj || otherRef != ref {
			return false
		}
	}
	for proj := range other.ProjectRevisions {
		_, hasProj := v.ProjectRevisions[proj]
		if !hasProj {
			return false
		}
	}
	return true
}
*/
