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
	"bytes"
	"os/exec"
	"testing"
)

func TestRepoInit(t *testing.T) {
	output, err := RepoInit(".", "--help")
	if err != nil {
		t.Logf("error: %v; stdout: %q", err, output)
		if exitErr, ok := err.(*exec.ExitError); ok {
			t.Logf("stderr: %q", exitErr.Stderr)
		}
		t.FailNow()
	}
	if !bytes.Contains(output, []byte("Usage: repo init")) {
		t.Fatalf("expected 'Usage: repo init', got: %q", output[:20])
	}
}

func TestRepoInitError(t *testing.T) {
	output, err := RepoInit(".", "--badarg")
	if err == nil {
		t.Fatalf("expected error, got none; output: %q", output)
	}
	_, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected exec.ExitError, got %T", err)
	}
}
