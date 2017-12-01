# [Gerrit](https://www.gerritcodereview.com/) Resource for [Concourse](https://concourse.ci/)

[![Build Status](https://travis-ci.org/google/concourse-resources.svg?branch=master)](https://travis-ci.org/google/concourse-resources)

Tracks Gerrit change revisions (patch sets).

*This is not an official Google product.*

## Usage

Define a new [resource type](https://concourse.ci/configuring-resource-types.html)
for your pipeline:

``` yaml
resource_types:
- name: gerrit
  type: docker-image
  source:
    repository: us.gcr.io/concourse-resources/gerrit-resource
```

## Source Configuration

* `url`: *Required.* The base URL of the Gerrit REST API.

* `query`: A Gerrit Search query matching desired changes. Defaults to
  `status:open`. You may want to specify a project like:
  `status:open project:my-project`. See Gerrit documentation on
  [Searching Changes](https://gerrit-documentation.storage.googleapis.com/Documentation/2.14.2/user-search.html).

* `cookies`: A string containing cookies in "Netscape cookie file format" (as
  supported by libcurl) to be used when connecting to Gerrit.  Usually used for
  authentication.

* `username`: A username for HTTP Basic authentication to Gerrit.

* `password`: A password for HTTP Basic authentication to Gerrit.

* `digest_auth`: If `true`, use HTTP Digest auth instead of Basic auth.

## Behavior

### `check`: Check for new revisions.

The Gerrit REST API is queried for revisions created since the given version
was created. If no version is given, the latest revision of the most recently
updated change is returned.

### `in`: Clone the git repository at the given revision.

The repository is cloned and the given revision is checked out.

#### Parameters

* `fetch_protocol`: A protocol name used to resolve a fetch URL for the given
  revision. For more information see the `fetch` field in the
  [Gerrit REST API documenation](https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#revision-info).
  Defaults to `http` or `anonymous http` if available.

* `fetch_url`: A URL to the Gerrit git repository where the given revision can
  be found. Overrides `fetch_protocol`.

### `out`

The given revision is updated with the given message and/or label(s).

#### Parameters

* `repository`: *Required.* The directory previously cloned by `in`; usually
  just the resource name.

* `message`: A message to be posted as a comment on the given revision.
  The message can contain build metadata variables. (e.g.: ${BUILD_ID})
  See the [Concourse.CI Metadata Documentation](https://concourse.ci/implementing-resources.html#section_resource-metadata
  for a complete list of variables.

* `message_file`: Path to a file containing a message to be posted as a comment
  on the given revision. This overrides `message` *unless* reading
  `message_file` fails, in which case `message` is used instead. If reading
  `message_file` fails and `message` is not specified then the `put` will fail.
  The message can contain build metadata variables. (e.g.: ${BUILD_ID})
  See the [Concourse.CI Metadata Documentation](https://concourse.ci/implementing-resources.html#section_resource-metadata
  for a complete list of variables.

* `labels`: A map of label names to integers to set on the given revision, e.g.:
  `{Verified: 1}`.

## Example Pipeline

``` yaml
resource_types:
- name: gerrit
  type: docker-image
  source:
    repository: us.gcr.io/concourse-resources/gerrit-resource

resources:
- name: example-gerrit
  type: gerrit
  source:
    url: https://review.example.com
    query: status:open project:example
    cookies: ((gerrit-cookies))

jobs:
- name: example-ci
  plan:
  # Trigger this job for every new patch set
  - get: example-gerrit
    version: every
    trigger: true

  - task: example-ci
    file: example-gerrit/ci.yml

  # After a successfuly build, mark the patch set Verified +1
  - put: example-gerrit
    params:
      repository: example-gerrit
      message: CI passed!
      labels: {Verified: 1}
```
