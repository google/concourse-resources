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
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testOutVersion = Version{
		ChangeId: "outChange",
		Revision: "outRev",
	}
)

func testOut(src Source, ver Version, params outParams) ResourceResponse {
	src.Url = testGerritUrl

	repoDir, err := ioutil.TempDir(testTempDir, "repo")
	if err != nil {
		panic(err)
	}

	err = ver.WriteToFile(filepath.Join(repoDir, gerritVersionFilename))
	if err != nil {
		panic(err)
	}
	params.Repository = filepath.Base(repoDir)

	return outMain(testJsonReader(outRequest{
		Source: src,
		Params: params,
	}), testTempDir)
}

func TestOutVersion(t *testing.T) {
	testOut(Source{}, testOutVersion, outParams{})
	assert.Equal(t, "outChange", testGerritLastChangeId)
	assert.Equal(t, "outRev", testGerritLastRevision)
}

func TestOutMessage(t *testing.T) {
	testOut(Source{}, testOutVersion, outParams{Message: "foo bar"})
	assert.Equal(t, "foo bar", testGerritLastReviewInput.Message)
}

func TestOutMessageFile(t *testing.T) {
	err := ioutil.WriteFile(
		filepath.Join(testTempDir, "message.txt"),
		[]byte("file msg"), 0600)
	assert.NoError(t, err)

	testOut(Source{}, testOutVersion, outParams{MessageFile: "message.txt"})
	assert.Equal(t, "file msg", testGerritLastReviewInput.Message)
}

func TestOutLabels(t *testing.T) {
	testOut(Source{}, testOutVersion, outParams{Labels: map[string]int{
		"Code-Review": 1,
		"Verified":    -1,
	}})
	assert.Equal(t, 1, testGerritLastReviewInput.Labels["Code-Review"])
	assert.Equal(t, -1, testGerritLastReviewInput.Labels["Verified"])
}
