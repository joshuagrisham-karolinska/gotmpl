# Go Template Rendering CLI and HTTP Server

An insanely simple CLI and HTTP Server which support rendering of [Go Text Templates](https://pkg.go.dev/text/template).

[Sprig v3](https://github.com/Masterminds/sprig) functions are also available for the templates (apart from that, there is nothing else other than plain `text/template`!).

> **Note:** In order to ensure consistency across different platforms, it is recommended to always build with the tag `timetzdata` so that the Timezone Data table will be the same (built-in table from Go) no matter where the tool is executed.

## CLI

CLI tool which renders a given `--template` using the given `--data` and writes the output to stdout.

```sh
# Build and run the CLI
go build -tags timetzdata -o bin/ ./cmd/gotmpl
./bin/gotmpl
./bin/gotmpl -t test.tmpl -d test.json
./bin/gotmpl -t test.tmpl -d test.xml
./bin/gotmpl -t test.tmpl -d test.yaml

# Run the CLI package without building
go run cmd/gotmpl/main.go
go run cmd/gotmpl/main.go -t test.tmpl -d test.json
go run cmd/gotmpl/main.go -t test.tmpl -d test.xml
go run cmd/gotmpl/main.go -t test.tmpl -d test.yaml

# Test with invalid data
go run cmd/gotmpl/main.go -t test.tmpl -d test-bad.json
```

## HTTP Server

```sh
# Build and run the server
go build -tags timetzdata -o bin/ ./cmd/gotmplserver
./bin/gotmplserver

# Run the server without building
go run cmd/gotmplserver/main.go

# Post a template and data to the API
curl -F "template=<test.tmpl" -F "data=<test.json" http://localhost:10000/gotmpl
curl -F "template=<test.tmpl" -F "data=<test.xml" http://localhost:10000/gotmpl
curl -F "template=<test.tmpl" -F "data=<test.yaml" http://localhost:10000/gotmpl

# Post with invalid data
curl -F "template=<test.tmpl" -F "data=<test-bad.json" http://localhost:10000/gotmpl
```

## Build specific version for multiple platforms

```sh
rm -rf bin
export GOTMPL_PKG=github.com/joshuagrisham-karolinska/gotmpl
export GOTMPL_VERSION=v0.0.1-alpha.3

export GOOS=windows
export GOARCH=amd64
go build -tags timetzdata -o "bin/gotmpl_${GOOS}_${GOARCH}.exe" -ldflags "-X ${GOTMPL_PKG}.Version=${GOTMPL_VERSION}" ./cmd/gotmpl
go build -tags timetzdata -o "bin/gotmplserver_${GOOS}_${GOARCH}.exe" -ldflags "-X ${GOTMPL_PKG}.Version=${GOTMPL_VERSION}" ./cmd/gotmplserver

export GOOS=linux
export GOARCH=amd64
go build -tags timetzdata -o "bin/gotmpl_${GOOS}_${GOARCH}" -ldflags "-X ${GOTMPL_PKG}.Version=${GOTMPL_VERSION}" ./cmd/gotmpl
go build -tags timetzdata -o "bin/gotmplserver_${GOOS}_${GOARCH}" -ldflags "-X ${GOTMPL_PKG}.Version=${GOTMPL_VERSION}" ./cmd/gotmplserver
```

## Build and run a GUI using WebAssembly

```sh
# Build the WebAssembly (including the embedded tzdata database so it will work in the browser; see: https://pkg.go.dev/time/tzdata)
GOOS=js GOARCH=wasm go build -tags timetzdata -o wasm/template.wasm ./wasm/main.go
# Copy Go's standard wasm_exec.js
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" wasm/
```

And run a simple web server to host an example of it:

```sh
# install goexec: go install github.com/shurcooL/goexec
# goexec 'http.ListenAndServe(`:8080`, http.FileServer(http.Dir(`.`)))'

# Run a tiny container using docker instead
docker run --name staticwebsite --rm -p 3000:3000 -v ${PWD}/wasm/:/home/static/:ro lipanski/docker-static-website:latest

# Then go to: http://localhost:3000

# Kill the staticwebsite container once you are finished
docker kill staticwebsite
```
