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
	"encoding/json"
	"os"
	"time"
)

const (
	timeStampLayout = "2006-01-02 15:04:05.999999999"
)

type Source struct {
	Url        string `json:"url"`
	Query      string `json:"query"`
	Cookies    string `json:"cookies"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	DigestAuth bool   `json:"digest_auth"`
}

type Version struct {
	ChangeId string    `json:"change_id"`
	Revision string    `json:"revision"`
	Created  time.Time `json:"created"`
}

func (v Version) Equal(o Version) bool {
	return v.ChangeId == o.ChangeId &&
		v.Revision == o.Revision &&
		v.Created.Equal(o.Created)
}

func (v Version) WriteToFile(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(v)
}

func (v *Version) ReadFromFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(v)
}

type VersionList []Version

func (vl VersionList) Len() int {
	return len(vl)
}

func (vl VersionList) Less(i, j int) bool {
	return vl[i].Created.Before(vl[j].Created)
}

func (vl VersionList) Swap(i, j int) {
	vl[i], vl[j] = vl[j], vl[i]
}
