.ONESHELL:
SHELL = /bin/bash
.PHONY: build clean update migrate-agent compile update clean get build clean-cache tidy

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOFILES=$(wildcard *.go)
SOURCE=$(GOBASE)/cmd/
BINNAME="gpwd"
CGO_ENABLED=1

# For tests
DATABASE_DSN = "postgres://dbuser:dbpass@localhost:5432/secrets?sslmode=disable"

## compile: Compile the binary.
build:
	$(MAKE) -s compile

## clean: Clean build files. Runs `go clean` internally.
clean:
	$(MAKE) go-clean

## update: Update modules
update:
	$(MAKE) go-update

migration-tools:
	@echo "  >  Install migration tools..."
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-server:
	@echo "  >  Do DB migrations..."
	@migrate -path db/migrations/server -database ${DATABASE_DSN} up

migrate-agent:
	@echo "  >  Do DB migrations..."
	@go-bindata -prefix "db/migrations/agent" -pkg migrations -o internal/agent/migrations/bindata.go db/migrations/agent

test: go-test go-vet

compile: go-clean go-get migrate-agent build-bin

go-update: go-clean go-clean-cache go-tidy go-download

go-clean:
	@echo "  >  Cleaning build cache"
	@GOBIN=$(GOBIN) go clean
	@rm -rf $(GOBIN)

go-get:
	@echo "  >  Checking if there is any missing dependencies..."
	@cd $(SOURCE); GOBIN=$(GOBIN) go get $(get)

build-bin:
	@echo "  >  Building binaries..."
	@cd $(SOURCE); go build --tags=sqlite_userauth -o $(GOBIN)/$(BINNAME) $(GOFILES)

go-clean-cache:
	@echo "  >  Clean modules cache..."
	@go clean -modcache

go-tidy:
	@echo "  >  Update modules..."
	@go mod tidy

go-download:
	@echo "  >  Download modules..."
	@go mod download

go-test:
	@echo "  >  Test project..."
	@go test ./...

go-vet:
	@echo "  >  Vet project..."
	@go vet ./...

go-proto:
	@echo "  >  Generate protobufs..."
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/proto/gpwd.proto
