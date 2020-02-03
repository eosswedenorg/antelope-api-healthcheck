
GO			= go
GOCCFLAGS 	= -v
GOLDFLAGS   =
PREFIX		= /usr/local

PROGRAM_NAME=eosio-api-healthcheck
SOURCES=server.go
DEPENDANCIES= github.com/firstrow/tcp_server \
	github.com/liamylian/jsontime/v2 \
	github.com/imroc/req

all: build
build: build/$(PROGRAM_NAME)

build/$(PROGRAM_NAME) : $(SOURCES)
	$(GO) build -o $@ $(GOCCFLAGS) $(GOLDFLAGS) $<

deps:
	$(GO) get $(DEPENDANCIES)

package_deb: build
	export PACKAGE_NAME="$(PROGRAM_NAME)" \
	export PACKAGE_VERSION="0.1.0" \
	export PACKAGE_PREFIX=$(PREFIX:/%=%) \
	export PACKAGE_PROGRAM="build/$(PROGRAM_NAME)" \
	&& ./scripts/build_deb.sh

clean:
	$(GO) clean
	$(RM) -rf build/
