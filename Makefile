VERSION ?= dev
LDFLAGS  = -ldflags "-X github.com/gbolanos-dev/worklog/internal/cli.Version=$(VERSION)"

.PHONY: build clean

build:
	go build $(LDFLAGS) -o bin/worklog ./cmd/worklog

clean:
	rm -rf bin/