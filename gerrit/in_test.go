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
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/google/concourse-resources/internal/resource"
)

var (
	testInVersion = Version{
		ChangeId: "Itestchange1",
		Revision: "deadbeef0",
		Created:  time.Unix(12345, 0),
	}

	testInDestDir string
)

func testIn(t *testing.T, src Source, ver Version, params inParams) (Version, []resource.MetadataField) {
	src.Url = testGerritUrl

	var err error
	testInDestDir, err = ioutil.TempDir(testTempDir, "repo")
	if err != nil {
		panic(err)
	}

	src.Url = testGerritUrl
	req := testRequest{Source: src, Version: ver, Params: params}
	var resp testResourceResponse
	assert.NoError(t, resource.TestInFunc(t, req, &resp, testInDestDir, in))
	return resp.Version, resp.Metadata
}

func TestInResponse(t *testing.T) {
	ver, metadata := testIn(t, Source{}, testInVersion, inParams{})
	assert.True(t, testInVersion.Equal(ver), "%v != %v", testInVersion, ver)
	assert.Contains(t, metadata, resource.MetadataField{Name: "project", Value: "testproject"})
	assert.Contains(t, metadata, resource.MetadataField{Name: "change subject", Value: "Test Subject"})
	assert.Contains(t, metadata, resource.MetadataField{Name: "revision uploader", Value: "Testy McTestface <testy@example.com>"})
	assert.Contains(t, metadata, resource.MetadataField{Name: "revision link", Value: fmt.Sprintf("%s/c/1/1", testGerritUrl)})
	assert.Contains(t, metadata, resource.MetadataField{Name: "commit author", Value: "Testy McTestface <testy@example.com>"})
	assert.Contains(t, metadata, resource.MetadataField{Name: "commit subject", Value: "Test Subject"})
	assert.Contains(t, metadata, resource.MetadataField{Name: "commit message", Value: "Commit message"})
}

func TestInGitInit(t *testing.T) {
	var initDir string
	mockGitWithArg("init", func(args []string, idx int) {
		for i := 0; i < idx; i++ {
			if args[i] == "-C" {
				initDir = args[i+1]
				break
			}
		}
	})

	testIn(t, Source{}, testInVersion, inParams{})
	assert.Equal(t, testInDestDir, initDir)
}

func TestInGitFetch(t *testing.T) {
	var fetchUrl, fetchRef string
	mockGitWithArg("remote", func(args []string, idx int) {
		fetchUrl = args[idx+3]
	})
	mockGitWithArg("fetch", func(args []string, idx int) {
		fetchRef = args[idx+2]
	})

	testIn(t, Source{}, testInVersion, inParams{})
	assert.Equal(t, fmt.Sprintf("%s/testproject.git", testGerritUrl), fetchUrl)
	assert.Equal(t, "refs/changes/1/1/1", fetchRef)
}

func TestInGitFetchProtocol(t *testing.T) {
	var fetchUrl, fetchRef string
	mockGitWithArg("remote", func(args []string, idx int) {
		fetchUrl = args[idx+3]
	})
	mockGitWithArg("fetch", func(args []string, idx int) {
		fetchRef = args[idx+2]
	})

	testIn(t, Source{}, testInVersion, inParams{FetchProtocol: "fake"})
	assert.Equal(t, "fake://example.com", fetchUrl)
	assert.Equal(t, "fake/ref", fetchRef)
}

func TestInGitFetchUrl(t *testing.T) {
	var fetchUrl, fetchRef string
	mockGitWithArg("remote", func(args []string, idx int) {
		fetchUrl = args[idx+3]
	})
	mockGitWithArg("fetch", func(args []string, idx int) {
		fetchRef = args[idx+2]
	})

	testIn(t, Source{}, testInVersion, inParams{FetchUrl: "some://otherurl"})
	assert.Equal(t, "some://otherurl", fetchUrl)
	assert.Equal(t, "refs/changes/1/1/1", fetchRef)
}

func TestInGitCookies(t *testing.T) {
	var cookiesPath string
	var cookiesFileData []byte
	mockGitWithArg("http.cookieFile", func(args []string, idx int) {
		var err error
		cookiesPath = args[idx+1]
		cookiesFileData, err = ioutil.ReadFile(cookiesPath)
		assert.NoError(t, err)
	})

	cookies := "localhost\tFALSE\t/\tFALSE\t9999999999\tauth\tbar\n"
	testIn(t, Source{Cookies: cookies}, testInVersion, inParams{})
	assert.Equal(t, cookies, string(cookiesFileData))

	// Cookie file should be deleted
	_, err := os.Stat(cookiesPath)
	assert.True(t, os.IsNotExist(err), "%s wasn't deleted", cookiesPath)
}

func TestInGitUsernamePassword(t *testing.T) {
	var credsPath string
	var credsOutput []byte
	mockGitWithArg("credential.helper", func(args []string, idx int) {
		var err error
		credHelper := args[idx+1]
		credsPath = credHelper[strings.LastIndex(credHelper, " ")+1:]
		credsOutput, err = exec.Command("sh", "-c", credHelper[1:]).CombinedOutput()
		assert.NoError(t, err, string(credsOutput))

		// Confirm credsPath exists right now.
		_, err = os.Stat(credsPath)
		assert.NoError(t, err)
	})

	password := `$(${'"\'\"` + "`"
	testIn(t, Source{Username: "bob", Password: password}, testInVersion, inParams{})
	assert.Equal(t, fmt.Sprintf("username=bob\npassword=%s\n", password), string(credsOutput))

	// Creds file should be deleted
	_, err := os.Stat(credsPath)
	assert.True(t, os.IsNotExist(err), "%s wasn't deleted", credsPath)
}

func TestInGerritVersionFile(t *testing.T) {
	testIn(t, Source{}, testInVersion, inParams{})

	var ver Version
	versionPath := filepath.Join(testInDestDir, gerritVersionFilename)
	assert.NoError(t, ver.ReadFromFile(versionPath))
	assert.True(t, testInVersion.Equal(ver), "%v != %v", testInVersion, ver)
}
