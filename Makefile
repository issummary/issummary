SHELL = /bin/bash

.PHONY: install-front
install-front:
	cd static && npm install

.PHONY: test-front
test-front: install-front
	cd static && npm test

.PHONY: build-front
build-front: install-front
	cd static && npm run build:prod

.PHONY: generate
generate:
	go generate

.PHONY: deps
deps:
	dep ensure -v

.PHONY: setup
setup:
	go get github.com/golang/dep/cmd/dep
	go get github.com/rakyll/statik

.PHONY: lint
lint: deps generate
	gometalinter

.PHONY: test
test: deps generate
	go test ./...

.PHONY: coverage
coverage: deps generate
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: codecov
codecov: deps coverage
	bash <(curl -s https://codecov.io/bash)

.PHONY: build
build: deps generate
	go build $(BUILD_PATH)

.PHONY: cross-build-snapshot
cross-build: deps generate
	goreleaser --rm-dist --snapshot

.PHONY: install
install: deps generate
	go install $(BUILD_PATH)

.PHONY: circleci
circleci:
	circleci build -e GITHUB_TOKEN=$GITHUB_TOKEN
