# OUI Lookup

OUI Lookup is used to locally store and query a Organizational Unit Identifier (OUI) database. Its intended audience are network administrators that do not want to rely on an always on Internet connectivity for retrieving this kind of information.

This tool has to be considered as **WORK IN PROGRESS** for now.

## Usage

Commands are not stable right now. Use `ouilookup -h` to obtain the most current usage information.

### Exit Codes

Exit codes are not stable right now. Do not rely on them, except that anything that is not 0 is some kind of error.

## Dependencies

This tool uses Go modules to handle dependencies.

## Running / Compiling

Use `go run ./...` to run the tool directly or `go build -o ouilookup ./...` to compile a binary.

Alternatively, Docker can be used to compile binaries by running `docker run --rm -v $PWD:/go/src -w /go/src golang:1.15 go build -o ouilookup ./...`. By passing the `GOOS` and `GOARCH` environment variables (via `-e`) this also enables cross compiling using Docker.

Tested with [go1.15](https://golang.org/doc/go1.15).

## Source

The original project is [hosted at GitLab](https://gitlab.com/rbrt-weiler/ouilookup), with a [copy over at GitHub](https://github.com/rbrt-weiler/ouilookup) for the folks over there.
