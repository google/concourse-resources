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

	"github.com/google/concourse-resources/internal/resource"
)

const (
	defaultQuery = "status:open"
)

var (
	updateStampTempDir = os.TempDir()
)

func init() {
	resource.RegisterCheckFunc(check)
}

func check(req resource.CheckRequest) error {
	var src Source
	var ver Version
	err := req.Decode(&src, &ver)
	if err != nil {
		return err
	}

	authMan := newAuthManager(src)
	defer authMan.cleanup()

	c, err := gerritClient(src, authMan)
	if err != nil {
		return fmt.Errorf("error setting up gerrit client: %v", err)
	}

	// Setup Gerrit query
	query := src.Query
	if query == "" {
		query = defaultQuery
	}

	var afterTime time.Time

	queryOpt := gerrit.QueryChangesOpt{}

	// If a version is requested, try to return that version in results.
	wantRequestedVersion := false

	var lastUpdate time.Time

	if ver.ChangeId == "" {
		// No version requested; fetch only the most recently updated change's
		// current revision.
		queryOpt.N = 1
		queryOpt.Fields = []string{"CURRENT_REVISION"}
	} else {
		// Check version requested; fetch changes updated since version was created.
		afterTime = ver.Created

		// As an optimization, try to read the latest change update timestamp from disk
		// and use that to filter instead.
		lastUpdate, err = readUpdatedStamp(src, ver)
		if err != nil {
			log.Println(err)
		}
		if !lastUpdate.IsZero() {
			afterTime = lastUpdate
		}

		query = fmt.Sprintf("(%s) AND after:{%s}",
			query, afterTime.UTC().Format(timeStampLayout))
		queryOpt.Fields = []string{"ALL_REVISIONS"}
		wantRequestedVersion = true
	}

	log.Printf("query: %q %+v", query, queryOpt)

	ctx := context.Background()
	changes, err := c.QueryChanges(ctx, query, queryOpt)
	if err != nil {
		return fmt.Errorf("error querying for changes: %v", err)
	}

	// Write latest change update timestamp to disk
	if len(changes) > 0 {
		lastChange := changes[len(changes)-1]
		if lastChange.Updated.Time().After(lastUpdate) {
			lastUpdate = lastChange.Updated.Time()
		}
		err = writeUpdatedStamp(src, ver, lastUpdate)
		if err != nil {
			log.Println(err)
		}
	}

	// Translate Gerrit changes into Versions
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
	for _, version := range versions {
		req.AddResponseVersion(version)
	}
	return nil
}

func updateStampFilename(src Source, ver Version) string {
	hash := sha1.New()
	fmt.Fprintf(hash, "%#v|%#v", src, ver)
	hashed := base32.StdEncoding.EncodeToString(hash.Sum([]byte{}))
	return filepath.Join(
		updateStampTempDir,
		fmt.Sprintf("concourse-gerrit-%s.stamp", hashed))
}

func readUpdatedStamp(src Source, ver Version) (t time.Time, err error) {
	f, err := os.Open(updateStampFilename(src, ver))
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		} else {
			err = fmt.Errorf("error opening update stamp file: %v", err)
		}
		return
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&t)
	if err != nil {
		err = fmt.Errorf("error reading update stamp file: %f", err)
	}
	return
}

func writeUpdatedStamp(src Source, ver Version, updated time.Time) error {
	f, err := os.OpenFile(updateStampFilename(src, ver), os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("error opening update stamp file: %v", err)
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(updated)
	if err != nil {
		return fmt.Errorf("error writing update stamp file: %v", err)
	}
	return nil
}
