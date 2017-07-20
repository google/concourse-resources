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
	"log"
	"os"
	"path/filepath"
)

var (
	Build = "untagged"
)

func main() {
	log.Printf("gerrit-resource build %s", Build)

	reqDecoder := json.NewDecoder(os.Stdin)

	var resp interface{}

	switch filepath.Base(os.Args[0]) {
	case "check":
		assertArgsLen(0)
		resp = checkMain(reqDecoder)
	case "in":
		assertArgsLen(1)
		resp = inMain(reqDecoder, os.Args[1])
	case "out":
		assertArgsLen(1)
		resp = outMain(reqDecoder, os.Args[1])
	default:
		log.Fatalf("unknown resource command %q", os.Args[0])
	}

	err := json.NewEncoder(os.Stdout).Encode(resp)
	fatalErr(err, "error writing response")
}

func assertArgsLen(expected int) {
	if len(os.Args) - 1 != expected {
		log.Fatalf(
			"%s takes exactly %d args; got %d",
			os.Args[0], expected, len(os.Args) - 1)
	}
}

func fatalErr(err error, fmt string, args ...interface{}) {
	if err != nil {
		if fmt == "" {
			fmt = "%v"
		} else {
			fmt += ": %v"
		}
		args = append(args, err)
		log.Fatalf(fmt, args...)
	}
}
