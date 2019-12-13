NAME ?= opw
VERSION = $(shell cat VERSION)
FILE ?= $(NAME)-$(VERSION)
GITUSER ?= shawncatz
GITREPO ?= opw

all: build

build: test build-all zip-all

build-all: build-linux build-darwin

build-darwin:
	@GOOS=darwin go build -o $(NAME)-darwin

build-linux:
	@GOOS=linux go build -o $(NAME)-linux

zip-darwin:
	zip $(FILE) ./$(NAME)-darwin

zip-linux:
	zip $(FILE) ./$(NAME)-linux

release: clean build
	git tag -f v$(VERSION)
	git push --tags

	github-release release \
        --user $(GITUSER) \
        --repo $(GITREPO) \
        --tag v$(VERSION) \
        --name "$(NAME)-v$(VERSION)"

	github-release upload \
        --user $(GITUSER) \
        --repo $(GITREPO) \
        --tag v$(VERSION) \
        --name "$(FILE)" \
        --file "$(FILE)"

release-delete:
	github-release delete -u $(GITUSER) -r $(GITREPO) -t v$(VERSION)

release-info:
	github-release info -u $(GITUSER) -r $(GITREPO)

clean:
	rm -f $(NAME) $(NAME)*.zip $(FILE)

test:
	go test -v ./...

deps:
	go get -u github.com/aktau/github-release
	go mod tidy
