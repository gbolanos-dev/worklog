VERSION ?= dev
LDFLAGS  = -ldflags "-X github.com/gbolanos-dev/worklog/internal/cli.Version=$(VERSION)"

.PHONY: build clean

build:
	go build $(LDFLAGS) -o bin/worklog ./cmd/worklog
	go build $(LDFLAGS) -o bin/wl ./cmd/wl

clean:
	rm -rf bin/