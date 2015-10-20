.PHONY: all test clean build install

GOFLAGS ?= $(GOFLAGS:)

all: install test


build:
	@gox -os="linux darwin" -arch="amd64" \
         -output "pkg/{{.OS}}_{{.Arch}}/hacher" \
		 .

install:
	@go get github.com/mitchellh/gox
	@go get $(GOFLAGS)

test: install
	@go test $(GOFLAGS)

clean:
	@go clean $(GOFLAGS) -i
