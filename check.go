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
	"context"
	"crypto/sha1"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"golang.org/x/build/gerrit"
)

const (
	defaultQuery = "status:open"
)

var (
	updateStampTempDir = os.TempDir()
)

type checkRequest struct {
	Source  `json:"source"`
	Version `json:"version"`
}

func checkMain(reqDecoder *json.Decoder) []Version {
	var req checkRequest
	err := reqDecoder.Decode(&req)
	fatalErr(err, "error reading request")

	authMan := newAuthManager(req.Source)
	defer authMan.cleanup()

	c, err := gerritClient(req.Source, authMan)
	fatalErr(err, "error setting up gerrit client")

	ctx := context.Background()

	query := req.Source.Query
	if query == "" {
		query = defaultQuery
	}

	var afterTime time.Time

	queryOpt := gerrit.QueryChangesOpt{}

	// If a version is requested, try to return that version in results.
	wantRequestedVersion := false

	var lastUpdate time.Time

	ver := req.Version
	if ver.ChangeId == "" {
		// No version requested; fetch only the most recently updated change's
		// current revision.
		queryOpt.N = 1
		queryOpt.Fields = []string{"CURRENT_REVISION"}
	} else {
		// Check version requested; fetch changes updated since version was created.
		afterTime = req.Version.Created

		// As an optimization, try to read the last change update timestamp from disk
		// and use that to filter instead.
		lastUpdate = readUpdatedStamp(req)
		if !lastUpdate.IsZero() {
			afterTime = lastUpdate
		}

		query = fmt.Sprintf("(%s) AND after:{%s}",
			query, afterTime.UTC().Format(timeStampLayout))
		queryOpt.Fields = []string{"ALL_REVISIONS"}
		wantRequestedVersion = true
	}

	log.Printf("query: %q %+v", query, queryOpt)

	changes, err := c.QueryChanges(ctx, query, queryOpt)
	fatalErr(err, "error querying for changes")

	if len(changes) > 0 {
		lastChange := changes[len(changes)-1]
		if lastChange.Updated.Time().After(lastUpdate) {
			lastUpdate = lastChange.Updated.Time()
		}
		writeUpdatedStamp(req, lastUpdate)
	}

	versions := VersionList{}
	for _, change := range changes {
		for revision, revisionInfo := range change.Revisions {
			include := false
			created := revisionInfo.Created.Time()
			if wantRequestedVersion && change.ID == ver.ChangeId && revision == ver.Revision {
				include = true
				wantRequestedVersion = false
			} else {
				include = created.After(afterTime)
			}
			if include {
				versions = append(versions, Version{
					ChangeId: change.ID,
					Revision: revision,
					Created:  created,
				})
			}
		}
	}
	if wantRequestedVersion {
		// Confirm the requested version still exists
		_, _, err := getVersionChangeRevision(c, ctx, ver)
		if err == nil {
			versions = append(versions, ver)
		} else {
			log.Printf("failed to fetch requested version: %v", err)
		}
	}
	sort.Sort(versions)
	return versions
}

func updateStampFilename(req checkRequest) string {
	hash := sha1.New()
	fmt.Fprintf(hash, "%#v", req)
	hashed := base32.StdEncoding.EncodeToString(hash.Sum([]byte{}))
	return filepath.Join(
		updateStampTempDir,
		fmt.Sprintf("concourse-gerrit-%s.stamp", hashed))
}

func readUpdatedStamp(req checkRequest) (t time.Time) {
	f, err := os.Open(updateStampFilename(req))
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("error opening update stamp file: %s", err)
		}
		return
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&t)
	if err != nil {
		log.Printf("error reading update stamp file: %s", err)
	}
	return t
}

func writeUpdatedStamp(req checkRequest, updated time.Time) {
	f, err := os.OpenFile(updateStampFilename(req), os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Printf("error opening update stamp file: %s", err)
		return
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(updated)
	if err != nil {
		log.Printf("error writing update stamp file: %s", err)
	}
}
