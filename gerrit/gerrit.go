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
	"fmt"
	"golang.org/x/build/gerrit"
)

func gerritClient(src Source, authMan *authManager) (*gerrit.Client, error) {
	if src.Url == "" {
		return nil, fmt.Errorf("source url is required")
	}
	auth, err := authMan.gerritAuth()
	if err != nil {
		return nil, err
	}
	return gerrit.NewClient(src.Url, auth), nil

}

func getVersionChangeRevision(
	client *gerrit.Client,
	ctx context.Context,
	ver Version,
	extraFields ...string,
) (*gerrit.ChangeInfo, *gerrit.RevisionInfo, error) {
	if ver.ChangeId == "" {
		return nil, nil, fmt.Errorf("version change_id required")
	}
	if ver.Revision == "" {
		return nil, nil, fmt.Errorf("version revision required")
	}

	change, err := client.GetChange(
		ctx, ver.ChangeId,
		gerrit.QueryChangesOpt{
			Fields: append([]string{"ALL_REVISIONS", "DETAILED_ACCOUNTS"}, extraFields...),
		})
	if err != nil {
		return nil, nil, fmt.Errorf(
			"error getting change %q: %v", ver.ChangeId, err)
	}

	revision, ok := change.Revisions[ver.Revision]
	if !ok {
		return nil, nil, fmt.Errorf(
			"no revision %q on change %q", ver.Revision, ver.ChangeId)
	}

	return change, &revision, nil
}
