APP=bin/containerd_util
all: build-mac
.PHONY: build-mac
build-mac:
	@go build -o ${APP} main.go

.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=linux go build -o ${APP} main.go

PHONY: tests
tests:
	@go test ./...

.PHONY: run
run:
	@go run ${APP}

.PHONY: clean
clean:
	@go clean