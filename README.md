# Google Cloud Platform Pub/Sub

A project for learning and teaching pub/sub using Go in the
[Google Cloud Platform](https://cloud.google.com/)
using [Google's Pubsub package](https://pkg.go.dev/cloud.google.com/go/pubsub?tab=doc)

[![Build Status](https://travis-ci.com/arcanericky/gcppubsubdemo.svg?branch=main)](https://travis-ci.com/arcanericky/gcppubsubdemo)
[![GoDoc](https://img.shields.io/badge/docs-GoDoc-brightgreen.svg)](https://pkg.go.dev/github.com/arcanericky/gcppubsubdemo?tab=doc)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](http://makeapullrequest.com)

---

## Purpose

`TL;DR`: If you are not a software developer, find another project.

There is no normal end user functionality provided. It's a project
implemented using Go to build demonstration executables so the user
can inspect the source code and run the executables to learn how to
implement _simple_ pub/sub functionality in Google Cloud Platform
using [Google's Pub/Sub package](https://pkg.go.dev/cloud.google.com/go/pubsub?tab=doc).

## Building

```
$ go build -o gcppubsubdemo ./cmd/...
```

## Pub/Sub Emulator

This project is meant to be used with the GCP Pub/Sub emulator. See
these [installation and execution instructions](https://cloud.google.com/pubsub/docs/emulator)
to get it installed and running.

In summary, install with

```
$ gcloud components install pubsub-emulator
$ gcloud components update
```

Execute with

```
$ gcloud beta emulators pubsub start --project=PUBSUB_PROJECT_ID
```

or

```
$ gcloud beta emulators pubsub start --host-port 127.0.0.1:9999
```

## Executing

### Quick Start for `gcppubsubdemo`

In another terminal, execute the publisher with

```
$ $(gcloud beta emulators pubsub env-init)
$ gcppubsubdemo publish
```

And in another terminal, execute the subscriber with

```
$ $(gcloud beta emulators pubsub env-init)
$ gcppubsubdemo subscribe
```

The default for both the `publish` and `subscribe` commands are to
publish and receive data continuously. Use Ctrl-C to stop execution
noting that in this case, the subscription ID and topics will remain
created in the pub/sub emulator. When using the `--once` option these
will be removed upon exit.

The `--verbose` option is particularly useful when learning how the
code flows. It enables DEBUG and TRACE level logging useful for
following the code logic.

### Help

There are command line options for `gcppubsubdemo`. Use the `--help`
flag to find them.

Examples

```
$ gcppubsubdemo --help
$ gcppubsubdemo subscribe --help
$ gcppubsubdemo publish --help
```

### API and Design

The design of `gcppubsubcdemo` consists of command-line user
interface code in `cmd/` driven with [Cobra](https://github.com/spf13/cobra).
This code in turn drives the `gcppubsubdemo` package (API) in
the project root. This is mainly to illustrate good architecture.
Because the core pub/sub functionality is implemented as an API
it can allow for re-use in other projects. However, I don't plan
on this package ever evolving past the `0.0.x` version because
this repository is for learning and this pub/sub code should be
implemented natively by other developers. You have been warned
when I break the API.

### Documentation

This project is documented with this `README.md`, scattered code
comments, and full GoDoc documentation as if it was a viable API.
This is done to show good practices.

## Contributing

Contributions and bug fixes (and there are bugs because I also
used this code to learn) are welcome. When contributing new
functionality, try to keep the code straightforward because the
purpose of this project is to help fellow developers learn. Try
to use functions rather than methods. Keep dependency injection
to a minimum, etc.

Some areas for consideration:

* Defect remedies (bug fixes)
* Pub/Sub API usage improvements
* Unit Tests (as long as the code remains readable by novice Gophers)
* Additional command line options to expose more pub/sub functionality
* Documentation corrections or additional information

## Inspiration

I had to learn this pubsub package for my day job. For me, learning
means writing and experimenting with code snippets. After I learned
the minimal amount to get the job done, I spent my evenings pulling
these code snippets together into self-contained meaningful project
to help others. This project is the result.