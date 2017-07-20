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
	"log"
	"os"

	"golang.org/x/build/gerrit"
)

var (
	cookiesTempDir = ""
)

type authManager struct {
	cookies string
	cookiesPath_ string
}

func newAuthManager(source Source) *authManager {
	return &authManager{
		cookies: source.Cookies,
	}
}

func (r *authManager) cookiesPath() (string, error) {
	if r.cookies == "" {
		return "", nil
	}
	if r.cookiesPath_ == "" {
		f, err := ioutil.TempFile(cookiesTempDir, "concourse-gerrit-cookies")
		if err != nil {
			return "", err
		}
		r.cookiesPath_ = f.Name()

		_, err = f.WriteString(r.cookies)
		if err != nil {
			return "", err
		}
	}
	return r.cookiesPath_, nil
}

func (r *authManager) gerritAuth() (gerrit.Auth, error) {
	if r.cookies != "" {
		cookiesPath, err := r.cookiesPath()
		if err != nil {
			return nil, err
		}
		return gerrit.GitCookieFileAuth(cookiesPath), nil
	} else {
		return gerrit.NoAuth, nil
	}
}

func (r *authManager) gitConfigArgs() ([]string, error) {
	cookiesPath, err := r.cookiesPath()
	if err != nil {
		return []string{}, err
	}
	return []string{"config", "http.cookieFile", cookiesPath}, nil
}

func (r *authManager) anonymous() bool {
	return r.cookies == ""
}

func (r *authManager) cleanup() {
	if r.cookiesPath_ != "" {
		err := os.Remove(r.cookiesPath_)
		if err != nil {
			log.Printf("error removing cookies file: %s", err)
		}
		r.cookiesPath_ = ""
	}
}


