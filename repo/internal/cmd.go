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

package internal

import (
	"log"
	"os"
	"os/exec"
)

func RepoInit(repoDir string, initArgs ...string) (output []byte, err error) {
	// Change to repoDir for this function.
	origDir, err := os.Getwd()
	if err != nil {
		return
	}
	err = os.Chdir(repoDir)
	if err != nil {
		return
	}
	defer os.Chdir(origDir)

	repoBin, err := Asset("repo")
	if err != nil {
		return
	}
	args := append([]string{"python", "-", "init"}, initArgs...)
	cmd := exec.Command("/usr/bin/env", args...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return
	}

	go func() {
		defer stdin.Close()
		stdin.Write(repoBin)
	}()

	return cmd.Output()
}

func RepoRun(repoDir string, repoArgs ...string) (output []byte, err error) {
	// Change to repoDir for this function.
	origDir, err := os.Getwd()
	if err != nil {
		return
	}
	err = os.Chdir(repoDir)
	if err != nil {
		return
	}
	defer os.Chdir(origDir)

	return exec.Command(".repo/repo/repo", repoArgs...).Output()
}

func LogExecErrors(prefix string, err error) bool {
	if err != nil {
		if err != nil {
			log.Printf("%s failed: %v", prefix, err)
			if exitErr, ok := err.(*exec.ExitError); ok {
				log.Printf("%s stderr:\n%s", prefix, exitErr.Stderr)
			}
		}
		return true
	}
	return false
}
