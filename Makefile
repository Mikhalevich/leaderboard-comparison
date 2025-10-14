MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
ROOT := $(dir $(MKFILE_PATH))
GOBIN ?= $(ROOT)/tools/bin
ENV_PATH = PATH=$(GOBIN):$(PATH)
BIN_PATH ?= $(ROOT)/bin
LINTER_NAME := golangci-lint
LINTER_VERSION := v2.5.0

.PHONY: all build test vendor compose-up compose-down install-linter lint fmt install-migrate create-migration ab-bench

all: build

build:
	go build -mod=vendor -o $(BIN_PATH)/api ./cmd/api/main.go

test:
	go test ./...

vendor:
	go mod tidy
	go mod vendor

compose-up:
	docker compose -f ./script/docker/docker-compose.yml up --build

compose-down:
	docker compose -f ./script/docker/docker-compose.yml down

install-linter:
	if [ ! -f $(GOBIN)/$(LINTER_VERSION)/$(LINTER_NAME) ]; then \
		echo INSTALLING $(GOBIN)/$(LINTER_VERSION)/$(LINTER_NAME) $(LINTER_VERSION) ; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN)/$(LINTER_VERSION) $(LINTER_VERSION) ; \
		echo DONE ; \
	fi

lint: install-linter
	$(GOBIN)/$(LINTER_VERSION)/$(LINTER_NAME) run --config .golangci.yml

fmt: install-linter
	$(GOBIN)/$(LINTER_VERSION)/$(LINTER_NAME) fmt --config .golangci.yml

install-migrate:
	if [ ! -f $(GOBIN)/migrate ]; then \
		echo "Installing migrate"; \
		GOBIN=$(GOBIN) go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.19.0; \
	fi

create-migration: install-migrate
	$(ENV_PATH) migrate create -ext sql -dir script/db/postgres/migrations -seq $(NAME)

ab-bench:
	docker run --rm --network=host jordi/ab -l -n 1000 -c 10 http://localhost:8080/leaderboard/random
