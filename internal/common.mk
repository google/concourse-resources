# Copyright 2017 Google Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

TREE_NAME=$(shell git write-tree)
DIRTY_MARK=-dirty-$(shell git rev-parse --short ${TREE_NAME})
BUILD=$(shell git describe --always --dirty=${DIRTY_MARK})

IMAGE_NAME=us.gcr.io/concourse-resources/${NAME}
IMAGE_TAG=${IMAGE_NAME}:${BUILD}

LDFLAGS=-ldflags "-X main.Build=${BUILD}"

# Prevent dynamic linking errors in docker containers
export CGO_ENABLED=0

build: clean
	go build ${LDFLAGS} -o build/${NAME}
	ln -s ${NAME} build/check
	ln -s ${NAME} build/in
	ln -s ${NAME} build/out

clean:
	rm -rf build

image: build
	cp Dockerfile build/
	docker build -t ${IMAGE_TAG} build/

image-push: image
	gcloud docker -- push ${IMAGE_TAG}

image-run: image
	docker run --rm -it ${IMAGE_TAG} /bin/sh
