
PROGRAM_NAME 	= eosio-api-healthcheck
export PROGRAM_VERSION = 1.2.3

GO			= go
PREFIX 		= /usr/local
export GOOS	= $(shell go env GOOS)
export GOARCH = $(shell go env GOARCH)
GOBUILDFLAGS  = -v -ldflags='-v -s -w -X main.VersionString=$(PROGRAM_VERSION)'

DPKG_BUILDPACKAGE = dpkg-buildpackage
DPKG_BUILDPACKAGE_FLAGS = -b -uc

SOURCES=src/main.go src/server.go src/parse_request.go

.PHONY: all build/$(PROGRAM_NAME) clean package_debian
all: build
build: build/$(PROGRAM_NAME)

build/$(PROGRAM_NAME) : $(SOURCES)
	$(GO) build -o $@ $(GOBUILDFLAGS) $^

test:
	$(GO) test -v ./...

coverage:
	$(GO) test -cover -v ./...

install: build
	PREFIX=$(PREFIX) DESTDIR=$(DESTDIR) scripts/install.sh $(GOOS)

package:
	PKGROOT=$(DESTDIR) BUILDDIR=$(realpath build) scripts/package.sh $(PKGTYPE)

package_debian:
	$(DPKG_BUILDPACKAGE) $(DPKG_BUILDPACKAGE_FLAGS)

package_freebsd: PKGTYPE = freebsd
package_freebsd: GOOS = freebsd
package_freebsd: DESTDIR = build/freebsdroot
package_freebsd: install package

clean:
	$(GO) clean
	$(RM) -rf build/
