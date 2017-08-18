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

	"github.com/google/concourse-resources/internal/resource"
)

func in(req resource.InRequest) error {
	var src Source
	var ver Version
	err := req.Decode(&src, &ver, nil)
	if err != nil {
		return err
	}

	err = initRepo(req.TargetDir(), src)
	if err != nil {
		return err
	}

	manifestFile, err := ioutil.TempFile("", "concourse-repo-manifest")
	if err != nil {
		return err
	}
	defer os.Remove(manifestFile.Name())

	_, err = manifestFile.WriteString(ver.Manifest)
	if err != nil {
		return err
	}
	manifestArg := fmt.Sprintf("--manifest-name=%s", manifestFile.Name())

	err = syncRepo(req.TargetDir(), src, manifestArg)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	resource.RegisterInFunc(in)
}
