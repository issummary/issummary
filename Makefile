REPO_OWNER = mpppk
REPO_NAME = issummary
BUILD_PATH = .
VERSION_PATH = cmd/
SHELL = /bin/bash

ifdef update
  u=-u
endif

.PHONY: deps
deps:
	dep ensure

.PHONY: setup
setup:
	go get ${u} github.com/golang/dep/cmd/dep

.PHONY: lint
lint: deps
	gometalinter

.PHONY: test
test: deps
	go test ./...

.PHONY: coverage
coverage: deps
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: codecov
codecov: deps coverage
	bash <(curl -s https://codecov.io/bash)

.PHONY: build
build: deps
	cd static && npm i && npm run build:prod
	go generate
	go build $(BUILD_PATH)

.PHONY: install
install: deps
	go install $(BUILD_PATH)

.PHONY: install
circleci:
	circleci build -e GITHUB_TOKEN=$GITHUB_TOKEN
