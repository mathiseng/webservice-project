Webservice
==========

A Go-based simple web service meant to be the subject for any tutorial
or even the project work.


__Prerequisites:__

* Go toolchain (install via system package manager or [by hand](https://go.dev/doc/install))


__Main interactions:__

1. Install dependencies: `go get -t ./...`
2. Run locally: `go run .`
3. Execute unit tests: `go test -race -v ./...`
4. Build artifact: `go build -o ./artifact.bin ./*.go`

For more details, please refer to the `Makefile`.
