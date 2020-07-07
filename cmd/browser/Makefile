LIFEGUARD=../../lifeguard

SHA256 := $(shell sha256sum $(LIFEGUARD) | cut -f 1 -d " ")
SHA512 := $(shell sha512sum $(LIFEGUARD) | cut -f 1 -d " ")

FLAGS="-X 'main.SHA256=$(SHA256)' -X 'main.SHA512=$(SHA512)'"

BASE=go build -v -ldflags $(FLAGS)

dev:
	# The last invocation of ldflags takes priority over all previous ones
	$(BASE) -ldflags "-X main.DisableSecurity=yes"

prod:
	$(BASE)