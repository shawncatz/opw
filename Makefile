NAME ?= opw
LINUX=build/linux/$(NAME)
DARWIN=build/darwin/$(NAME)
VERSION = $(shell cat VERSION)
GITUSER ?= shawncatz
GITREPO ?= opw
LDFLAGS = -s -w -X github.com/shawncatz/opw/cmd.Version=v$(VERSION)

all: test ## Build and run tests

test: ## Run tests
	go test -v ./...

clean: ## Remove previous build
	rm -rf build

build: test linux darwin ## Build binaries
	@echo version: $(VERSION)

linux: $(LINUX).zip ## Build for Linux

darwin: $(DARWIN).zip ## Build for Darwin (macOS)

$(LINUX):
	env GOOS=linux GOARCH=amd64 go build -i -v -o $(LINUX) -ldflags="$(LDFLAGS)"

$(DARWIN):
	env GOOS=darwin GOARCH=amd64 go build -i -v -o $(DARWIN) -ldflags="$(LDFLAGS)"

$(LINUX).zip: $(LINUX)
	cd build/linux/ && zip $(NAME).zip $(NAME)

$(DARWIN).zip: $(DARWIN)
	cd build/darwin/ && zip $(NAME).zip $(NAME)

release: clean build release-create release-darwin release-linux ## Create Github Release

release-create:
	git tag -f v$(VERSION)
	git push --tags
	github-release release \
        --user $(GITUSER) \
        --repo $(GITREPO) \
        --tag v$(VERSION) \
        --name "$(NAME)-v$(VERSION)"

release-darwin:
	github-release upload \
        --user $(GITUSER) \
        --repo $(GITREPO) \
        --tag v$(VERSION) \
        --name "$(NAME)-linux-amd64-$(VERSION).zip" \
        --file "build/linux/$(NAME).zip"

release-linux:
	github-release upload \
        --user $(GITUSER) \
        --repo $(GITREPO) \
        --tag v$(VERSION) \
        --name "$(NAME)-darwin-amd64-$(VERSION).zip" \
        --file "build/darwin/$(NAME).zip"

release-delete:
	github-release delete -u $(GITUSER) -r $(GITREPO) -t v$(VERSION)

release-info:
	github-release info -u $(GITUSER) -r $(GITREPO)

deps: ## Install Dependencies
	go get -u github.com/aktau/github-release
	go mod tidy

.PHONY: all test clean

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
