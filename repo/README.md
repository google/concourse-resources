# [Repo](https://source.android.com/source/using-repo) Resource for [Concourse](https://concourse.ci/)

[![Build Status](https://travis-ci.org/google/concourse-resources.svg?branch=master)](https://travis-ci.org/google/concourse-resources)

Tracks changes to repo projects.

*This is not an official Google product.*

## Usage

Define a new [resource type](https://concourse.ci/configuring-resource-types.html)
for your pipeline:

``` yaml
resource_types:
- name: repo
  type: docker-image
  source:
    repository: us.gcr.io/concourse-resources/repo-resource
```

## Source Configuration

* `manifest_url`: *Required.* The URL of the repo manifest repository.
  (See: `--manifest-url` in `repo init --help`.)

* `manifest_name`: The name of the manifest file to use.
  (See: `--manifest-name` in `repo init --help`.)

* `manifest_branch`: The name of the manifest repository branch.
  (See: `--manifest-branch` in `repo init --help`.)

* `groups`: An array of manifest group names to use.
  (See: `--groups` in `repo init --help`.)

* `init_options` and `sync_options`: A map of `repo init` or `repo sync` option
  names (full names only, without `--`) to option values. Values may be strings,
  numbers, or booleans (where `true` represents setting a valueless flag option
  and `false` represents unsetting an option). These options will override any
  other option, including hardcoded defaults.
  *Be careful! This is an advanced feature that can break the resource!*
  (See: `repo init --help` and `repo sync --help`.)

## Behavior

### `check`: Check for new revisions.

`repo init` and `repo sync` are called with the provided options. `repo
manifest` is used to take a snapshot of the state of the repositories, which is
used as the resource version.

### `in`: Clone the repo project repositories at the given revisions.

`repo init` and `repo sync` are called with the provided options, using the
provided version as the manifest. This results in a snapshot of the project
repositories.

### `out`

This resource does not implement `out`.

## Example Pipeline

``` yaml
resource_types:
- name: repo
  type: docker-image
  source:
    repository: us.gcr.io/concourse-resources/repo-resource

resources:
- name: example-repo
  type: repo
  source:
    manifest_url: https://source.example.com/manifest.git
    groups:
      - tests

jobs:
- name: example-ci
  plan:
  - get: example-repo
    trigger: true

  - task: example-ci
    file: example-repo/testing/ci.yml
```
