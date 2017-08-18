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
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/concourse-resources/internal/resource"
)

func in(req resource.InRequest) error {
	var src Source
	var ver Version
	err := req.Decode(&src, &ver, nil)
	if err != nil {
		return err
	}

	err = repoInit(req.TargetDir(), src)
	if err != nil {
		return err
	}
	
	manifestFile := filepath.Join(req.TargetDir(), ".repo", "manifest.xml")

	err = os.Remove(manifestFile)
	if err != nil {
		return fmt.Errorf("error removing manifest symlink: %v", err)
	}

	err = ioutil.WriteFile(manifestFile, []byte(ver.Manifest), 0644)
	if err != nil {
		return fmt.Errorf("error writing snapshot manifest: %v", err)
	}

	err = repoSync(req.TargetDir(), src)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	resource.RegisterInFunc(in)
}
