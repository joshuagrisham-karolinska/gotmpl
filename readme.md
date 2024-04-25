# Go Template Rendering CLI and HTTP Server

An insanely simple CLI and HTTP Server which support rendering of [Go Text Templates](https://pkg.go.dev/text/template).

[Sprig v3](https://github.com/Masterminds/sprig) functions are also available for the templates (apart from that, there is nothing else other than plain `text/template`!).

## CLI

```sh
# Build and run the CLI
go build ./cmd/gotmpl
./gotmpl
./gotmpl -t test.tmpl -d test.json

# Run the CLI package without building
go run cmd/gotmpl/main.go
go run cmd/gotmpl/main.go -t test.tmpl -d test.json

# Test with invalid data
go run cmd/gotmpl/main.go -t test.tmpl -d test-bad.json

# Build with a specific version
go build -ldflags "-X github.com/joshuagrisham-karolinska/gotmpl.Version=v0.0.1-alpha.0" ./cmd/gotmpl
```

## Server

```sh
# Build and run the server
go build ./cmd/gotmpl-server
./gotmpl-server

# Run the server without building
go run cmd/gotmpl-server/main.go

# Post a template and data to the API
curl -F "template=<test.tmpl" -F "data=<test.json" http://localhost:10000/gotmpl

# Post with invalid data
curl -F "template=<test.tmpl" -F "data=<test-bad.json" http://localhost:10000/gotmpl
```
