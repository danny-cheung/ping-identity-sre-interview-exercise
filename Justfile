default:
    @echo "Available recipes:"
    @just --list

deps:
  go mod tidy

fmt: deps
  go fmt ./...

test: fmt
  go test -v ./...

build: test
  ko build --local