# Overview

[![Go Reference][badge-svg]][badge-link]
[![Build and Tests][tests-svg]][tests-link]
[![Code Coverage][codecov-svg]][codecov-link]

`go-improvmx` is a golang wrapper around the [ImprovMX][1] [API][2]. It was
written primarily for a [terraform][3] provider, however others might get use
out of it.

## Installation

Simply use `go get` to add `go-improvmx` to your `go.mod`

```console
$ go get occult.work/improvmx@v1.0.0
```

## Development

`go-improvmx` uses [`task`][4] to run the most common operations. These tasks
are then duplicated with extra flags within the GitHub Actions workflow.

Task has [instructions on installation](https://taskfile.dev/#/installation).
Once installed, simply run `task` from the project directory.

## Dependencies

The current list of third party libraries are

 - [resty](https://github.com/go-resty/resty)

The following libraries are used for *testing only*

 - [testify](https://github.com/stretchr/testify)
 - [mux](https://github.com/gorilla/mux)

[1]: https://improvmx.com/
[2]: https://improvmx.com/api/
[3]: https://www.terraform.io/
[4]: https://github.com/go-task/task

[codecov-svg]: https://codecov.io/gh/slurps-mad-rips/go-improvmx/branch/main/graph/badge.svg?token=4ngB0iw6qf
[tests-svg]: https://github.com/slurps-mad-rips/go-improvmx/actions/workflows/tests.yml/badge.svg
[badge-svg]: https://pkg.go.dev/badge/occult.work/improvmx.svg

[codecov-link]: https://codecov.io/gh/slurps-mad-rips/go-improvmx
[tests-link]: https://github.com/slurps-mad-rips/go-improvmx/actions/workflows/tests.yml
[badge-link]: https://pkg.go.dev/occult.work/improvmx
