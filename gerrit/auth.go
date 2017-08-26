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

type credentialsAuth struct {
	username	string
	password	string
}

type cookiesAuth struct {
	cookies      string
	cookiesPath_ string
}

type AuthManager interface {
	GerritAuth() (gerrit.Auth, error)
	GitConfigArgs() ([]string, error)
	Anonymous() bool
	Cleanup()
}

func NewAuthManager(source Source) AuthManager {
	return newAuthManager(source, newCookieAuth, newCredentialsAuth, newAnonymousAuth)
}

func newCookieAuth(source *Source) AuthManager { return &cookiesAuth{cookies: source.Cookies } }

func newAnonymousAuth(_ *Source) AuthManager { return &cookiesAuth{} }

func newCredentialsAuth(source * Source) AuthManager { return &credentialsAuth{username: source.Username, password: source.Password } }

type authManagerCreator func(source *Source) AuthManager

func newAuthManager(source Source, authManagers ...authManagerCreator) AuthManager {
	var authManager AuthManager
	for _, creatorFunc := range authManagers {
		authManager = creatorFunc(&source)
		if !authManager.Anonymous() { break }
	}
	return authManager
}

func (r *cookiesAuth) cookiesPath() (string, error) {
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

func (r *cookiesAuth) GerritAuth() (gerrit.Auth, error) {
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

func (r *cookiesAuth) GitConfigArgs() ([]string, error) {
	cookiesPath, err := r.cookiesPath()
	if err != nil {
		return []string{}, err
	}
	return []string{"config", "http.cookieFile", cookiesPath}, nil
}

func (r *cookiesAuth) Anonymous() bool {
	return r.cookies == ""
}

func (r *cookiesAuth) Cleanup() {
	if r.cookiesPath_ != "" {
		err := os.Remove(r.cookiesPath_)
		if err != nil {
			log.Printf("error removing cookies file: %s", err)
		}
		r.cookiesPath_ = ""
	}
}

func (r *credentialsAuth) GerritAuth() (gerrit.Auth, error) {
	return gerrit.BasicAuth(r.username, r.password), nil
}

func (r *credentialsAuth) GitConfigArgs() ([]string, error) {
	return []string{}, nil
}

func (r *credentialsAuth) Anonymous() bool {
	return r.username == ""
}

func (r *credentialsAuth) Cleanup() { }
