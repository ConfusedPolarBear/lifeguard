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

dev: test server web
prod: test server web-prod

server:
	go build -v -ldflags $(FLAGS)

web-prod:
	npm run prod

web:
	npm run dev

archive:
	tar -cf $(ARCHIVE) lifeguard web/ assets/ scripts/ example_config.ini
	tar -f $(ARCHIVE) --delete web/src
	gzip -f $(ARCHIVE)

test:
	go test

install:
	scripts/install.sh

clean:
	go clean
	rm -rf web/dist/
