NAME ?= opw
VERSION = $(shell cat VERSION)
FILE ?= $(NAME)-$(VERSION).zip
GITUSER ?= shawncatz
GITREPO ?= opw

all: build

build: test $(NAME) $(FILE)

$(NAME):
	@GOOS=linux go build -o $(NAME)

$(FILE):
	zip $(FILE) ./$(NAME)

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
