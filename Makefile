VERSION=0.1.2
LDFLAGS=-ldflags "-w -s -X main.version=${VERSION}"

all: mackerel-plugin-pdns

.PHONY: mackerel-plugin-pdns linux check lint

mackerel-plugin-pdns: cmd/mackerel-plugin-pdns/*.go
	go build $(LDFLAGS) -o mackerel-plugin-pdns ./cmd/mackerel-plugin-pdns/

linux: cmd/mackerel-plugin-pdns/*.go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o mackerel-plugin-pdns ./cmd/mackerel-plugin-pdns/

check:
	go test -v ./...

lint:
	golangci-lint run --timeout 5m ./...