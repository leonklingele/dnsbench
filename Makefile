.PHONY: all
all: build

.PHONY: build
build:
	go build ./...

.PHONY: install
install: build
	go install ./...
