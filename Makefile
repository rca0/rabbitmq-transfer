build:
	CGO_ENABLED=0 CGO111MODULES=on go build -ldflags="-s -w"
.PHONY: build

fmt:
	gofmt -w .
.PHONY: fmt

.DEFAULT_GOAL := build
