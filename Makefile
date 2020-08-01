CONFIG=github.com/ConfusedPolarBear/lifeguard/pkg/config
ARCHIVE=lifeguard.tar

TIME := $(shell date)
VERSION := $(shell go version)
COMMIT := $(shell git rev-list -1 HEAD)

# returns " (modified)" if local changes have been made to the git repository, "" otherwise
MODIFIED := $(shell git diff --no-ext-diff --quiet || echo " (modified)")

FLAGS="-X '$(CONFIG).Commit=$(COMMIT)'\
	-X '$(CONFIG).BuildTime=$(TIME)'\
	-X '$(CONFIG).GoVersion=$(VERSION)'\
	-X '$(CONFIG).Modified=$(MODIFIED)'"

dev: test server browser web
prod: test server browser-prod web-prod
backend: server browser

server:
	go build -v -ldflags $(FLAGS)

browser-prod:
	make -C cmd/browser prod
	mv cmd/browser/browser ./

browser:
	make -C cmd/browser
	mv cmd/browser/browser ./

web-prod:
	npm audit
	npm run prod

web:
	npm run dev

archive:
	tar -cf $(ARCHIVE) lifeguard browser web/ assets/ scripts/ config/
	tar -f $(ARCHIVE) --delete web/src web/dist/*.map assets/Insomnia*
	gzip -f $(ARCHIVE)

test:
	go test

install:
	scripts/install.sh

clean:
	go clean
	rm -rf web/dist/

# Bug fix: make often incorrectly asserts that some targets are 'up to date' when they are not
.PHONY: web web-prod browser browser-prod
