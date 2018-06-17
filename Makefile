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
deps: generate
	dep ensure -v

.PHONY: setup
setup:
	go get ${u} github.com/golang/dep/cmd/dep
	go get github.com/rakyll/statik

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
	go build $(BUILD_PATH)

.PHONY: install
install: deps
	go install $(BUILD_PATH)

.PHONY: circleci
circleci:
	circleci build -e GITHUB_TOKEN=$GITHUB_TOKEN
