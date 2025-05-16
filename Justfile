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

run: test
  go run main.go

kind-test: test
  @ko version || brew install ko
  @docker --version || brew install docker
  @kind --version || brew install kind
  @kubectl version || brew install kubectl

  # Setup the KinD cluster
  @kind get clusters | grep "{{kind_cluster_name}}" || kind create cluster --name "{{kind_cluster_name}}" --config kind-config.yaml
  # Install NginX ingress controller
  @kubectl --context "{{kind_context}}" apply --filename https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/kind/deploy.yaml
  @sleep 1
  @kubectl --context "{{kind_context}}" wait --namespace ingress-nginx --for=condition=ready pod --selector=app.kubernetes.io/component=controller --timeout=180s

  # Build the binary, build the image, and push image to Docker Hub
  @KO_DOCKER_REPO=dann7387/ping-identity-sre-interview-exercise ko build --bare

  # Deploy the app to Kubernetes
  kubectl --context "{{kind_context}}" apply -f kubernetes.yaml

  @echo
  @echo "==========================================="
  @echo "App is accessible on http://localhost:30000"
  @echo "==========================================="
