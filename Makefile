BINARY=wuzzlmoasta
MAIN=./cmd/wuzzlmoasta

.PHONY: default build clean release test format run

default: build

build:
	@go generate ./...
	@go build -o bin/$(BINARY) $(MAIN)

clean:
	@rm -rf bin

release: clean
	@go generate ./...
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_darwin_amd64 $(MAIN)
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_linux_amd64 $(MAIN)
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_windows_amd64 $(MAIN)

test:
	@go test -v ./...

format:
	@go fmt ./...

run:
	@go run ./cmd/wuzzlmoasta