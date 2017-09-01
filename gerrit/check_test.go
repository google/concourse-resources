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
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/google/concourse-resources/internal/resource"
)

func testCheck(t *testing.T, src Source, ver Version) []Version {
	src.Url = testGerritUrl
	req := testRequest{Source: src, Version: ver}
	var versions []Version
	assert.NoError(t, resource.TestCheckFunc(t, req, &versions, check))
	return versions
}

func TestCheckSourceQuery(t *testing.T) {
	testCheck(t, Source{Query: "my query"}, Version{})
	assert.Equal(t, "my query", testGerritLastQ)
}

func TestCheckSourceCookies(t *testing.T) {
	cookies := "localhost\tFALSE\t/\tFALSE\t9999999999\tauth\tbar\n"
	testCheck(t, Source{Cookies: cookies}, Version{})
	assert.True(t, testGerritLastAuthenticated)
	cookie, err := testGerritLastRequest.Cookie("auth")
	assert.NoError(t, err)
	assert.Equal(t, "bar", cookie.Value)
}

func TestCheckSourceUsernamePassword(t *testing.T) {
	testCheck(t, Source{Username: "bob", Password: "dog"}, Version{})
	assert.True(t, testGerritLastAuthenticated)
	authHeader := testGerritLastRequest.Header.Get("authorization")
	assert.Equal(t, "Basic Ym9iOmRvZw==", authHeader) // == Base64("bob:dog")
}

func TestCheckSourceDigestAuth(t *testing.T) {
	testCheck(t, Source{Username: "bob", Password: "dog", DigestAuth: true}, Version{})
	assert.True(t, testGerritLastAuthenticated)
	authHeader := testGerritLastRequest.Header.Get("authorization")
	assert.Contains(t, authHeader, "Digest ")
}

func TestCheckWithoutVersion(t *testing.T) {
	versions := testCheck(t, Source{}, Version{})
	assert.Equal(t, "status:open", testGerritLastQ)
	assert.Equal(t, 1, testGerritLastN)

	assert.Len(t, versions, 1)
	assert.Equal(t, "testproject~testbranch~Itestchange1", versions[0].ChangeId)
	assert.Equal(t, "deadbeef0", versions[0].Revision)
	assert.True(t, time.Unix(100, 0).Equal(versions[0].Created))
}

func TestCheckWithNewVersions(t *testing.T) {
	versions := testCheck(t, Source{}, Version{
		ChangeId: "Itestchange1",
		Revision: "deadbeef0",
		Created:  time.Unix(1, 0),
	})
	assert.Len(t, versions, 10)
	assert.Equal(t, "(status:open) AND after:{1970-01-01 00:00:01}", testGerritLastQ)
	assert.Equal(t, 0, testGerritLastN)
	assert.Equal(t, "Itestchange1", testGerritLastChangeId)
}

func TestCheckWithoutNewVersions(t *testing.T) {
	versions := testCheck(t, Source{}, Version{
		ChangeId: "Itestchange1",
		Revision: "deadbeef0",
		Created:  time.Unix(50000, 0),
	})
	assert.Len(t, versions, 1)
	assert.Equal(t, "Itestchange1", versions[0].ChangeId)
}

func TestCheckVersionsSorted(t *testing.T) {
	versions := testCheck(t, Source{}, Version{
		ChangeId: "Itestchange1",
		Revision: "deadbeef0",
		Created:  time.Unix(2, 0),
	})
	assert.Len(t, versions, 10)
	assert.True(t, sort.SliceIsSorted(versions, func(i, j int) bool {
		return versions[i].Created.Before(versions[j].Created)
	}))
}

func TestCheckWithBadRevision(t *testing.T) {
	versions := testCheck(t, Source{}, Version{
		ChangeId: "Itestchange1",
		Revision: "badrevision",
	})
	assert.Equal(t, "Itestchange1", testGerritLastChangeId)

	assert.Len(t, versions, 9)
}

func TestCheckTimestampCaching(t *testing.T) {
	testCheck(t, Source{Query: "foo"}, Version{
		ChangeId: "Itestchange1",
		Revision: "deadbeaf0",
		Created:  time.Unix(100, 0),
	})
	assert.Equal(t, "(foo) AND after:{1970-01-01 00:01:40}", testGerritLastQ)

	testCheck(t, Source{Query: "foo"}, Version{
		ChangeId: "Itestchange1",
		Revision: "deadbeaf0",
		Created:  time.Unix(100, 0),
	})
	assert.Equal(t, "(foo) AND after:{1970-01-01 05:35:00}", testGerritLastQ)

	testCheck(t, Source{Query: "bar"}, Version{
		ChangeId: "Itestchange1",
		Revision: "deadbeaf0",
		Created:  time.Unix(100, 0),
	})
	assert.Equal(t, "(bar) AND after:{1970-01-01 00:01:40}", testGerritLastQ)
}
