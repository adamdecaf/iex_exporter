VERSION := $(shell grep -Eo '(\d\.\d\.\d)(-dev)?' main.go | head -n1)

.PHONY: build deps docker

build:
	go fmt ./...
	go vet ./...
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/iex-exporter-darwin .
	CGO_ENABLED=0 GOOS=linux  GOARCH=amd64 go build -o bin/iex-exporter-linux .

docker: build
	docker build --pull -t adamdecaf/iex_exporter:$(VERSION) -f Dockerfile .

test:
	go test -v ./...
