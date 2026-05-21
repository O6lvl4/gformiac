VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS  = -X github.com/O6lvl4/gformiac/cmd.Version=$(VERSION)

.PHONY: build test lint clean

build:
	go build -ldflags "$(LDFLAGS)" -o gformiac .

test:
	go test ./... -v

lint:
	go vet ./...

clean:
	rm -f gformiac
