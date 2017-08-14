NAME=gerrit-resource

TREE_NAME=$(shell git write-tree)
DIRTY_MARK=-dirty-$(shell git rev-parse --short ${TREE_NAME})
BUILD=$(shell git describe --always --dirty=${DIRTY_MARK})

IMAGE_NAME=us.gcr.io/concourse-gerrit/resource
IMAGE_TAG=${IMAGE_NAME}:${BUILD}

LDFLAGS=-ldflags "-X main.Build=${BUILD}"

export CGO_ENABLED=0

build: clean
	git submodule update --init --recursive
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
	docker run -it ${IMAGE_TAG} /bin/bash

test:
	git submodule update --init --recursive
	go test .
