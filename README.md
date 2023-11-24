Webservice
==========

A Go-based simple web service meant to be the subject of any tutorial
or even used the project work.


#### Prerequisites:

* Go toolchain (install via system package manager or [by hand](https://go.dev/doc/install))


#### Build

1. Install dependencies: `go get -t ./...`
2. Run locally: `go run .`
3. Execute unit tests: `go test -race -v ./...`
4. Build artifact: `go build -o ./artifact.bin ./*.go`

To build for another platform, set `GOOS` and `GOARCH`. To yield a static binary (fully
self-contained, no dynamic linking) set `CGO_ENABLED=0`. For more details, please refer
to the [Makefile](./Makefile).
