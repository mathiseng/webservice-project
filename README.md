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


#### Run:

```bash
HOST=0.0.0.0 PORT=8080 ./artifact.bin
```


#### Interact:

##### Landing page 

plain text:
```bash
curl http://localhost:8080
```

HTML:
```bash
curl --header 'Content-Type: text/html; charset=utf-8' http://localhost:8080
# or just open in a browser
```


##### Health check

```bash
curl http://localhost:8080/health
```


##### Server side environment variables

List environment variables visible by the webservice process if environment
is not `production`.

```bash
curl http://localhost:8080/env
```


##### State life cycle

URL slug is used as identifier and the body is the actual *data* being stored.
Please note, when writing (add or change) something, `Content-Type` must be set
in the request header.

Write an entry:
```bash
curl \
  -X PUT \
  --header 'Content-Type: text/plain; charset=utf-8' \
  --data 'foo' \
  http://localhost:8080/state/bar
```

Obtain an entry:
```bash
curl \
  -X GET \
  http://localhost:8080/state/bar
```

Remove an entry:
```bash
curl \
  -X DELETE \
  --verbose \
  http://localhost:8080/state/bar
```

List all existing entries (returns JSON or plain text, depending on the `Accept` header):
```bash
curl \
  -X GET
  --header 'Accept: text/plain'\
  http://localhost:8080/states
```

Upload an entire file:
```bash
curl \
  -X PUT \
  --header 'Content-Type: application/pdf' \
  --upload-file ./example.pdf \
  http://localhost:8080/state/pdf-doc
```

Download a file:
```bash
curl \
  -X GET \
  --output ./example-copy.pdf \
  http://localhost:8080/state/pdf-doc
```
