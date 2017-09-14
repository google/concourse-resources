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
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/concourse-resources/internal/resource"
)

var checkRepoDir = filepath.Join(os.TempDir(), "concourse-repo")

func check(req resource.CheckRequest) error {
	var src Source
	var ver Version
	err := req.Decode(&src, &ver)
	if err != nil {
		return err
	}

	// Create and init repo if it doesn't exist.
	_, err = os.Stat(checkRepoDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(checkRepoDir, 0755)
		if err != nil {
			return fmt.Errorf("error creating repo dir %s: %v", checkRepoDir, err)
		}

		err = repoInit(checkRepoDir, src)
		if err != nil {
			// Clean up the repo dir so we can try init again later
			os.RemoveAll(checkRepoDir)
			return err
		}
	} else if err != nil {
		return fmt.Errorf("error accessing repo dir %s: %v", checkRepoDir, err)
	}

	err = repoSync(checkRepoDir, src)
	if err != nil {
		return err
	}

	newVer, err := getCurrentVersion(checkRepoDir)
	if (ver != Version{}) && (ver != newVer) {
		req.AddResponseVersion(ver)
	}
	req.AddResponseVersion(newVer)

	return nil
}

func init() {
	resource.RegisterCheckFunc(check)
}
