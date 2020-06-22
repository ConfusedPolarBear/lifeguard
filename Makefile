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

all: test build

build:
	go build -v -ldflags $(FLAGS)
	npm run dev

archive:
	tar -cf $(ARCHIVE) lifeguard web/ example_config.ini
	tar -f $(ARCHIVE) --delete web/src
	gzip -f $(ARCHIVE)

test:
	go test

clean:
	go clean
	rm -rf web/dist/
