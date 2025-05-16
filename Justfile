kind_cluster_name := "local"
kind_context := "kind-" + kind_cluster_name

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
  go build main.go

kind-test: test
  ko version || brew install ko
  docker --version || brew install docker
  kind --version || brew install kind
  kubectl version || brew install kubectl

  kind get clusters | grep "{{kind_cluster_name}}" || kind create cluster --name "{{kind_cluster_name}}"
  KO_DOCKER_REPO=kind.local KIND_CLUSTER_NAME="{{kind_cluster_name}}" ko build
  kubectl --context "{{kind_context}}" apply -f kubernetes.yaml