
GO			= go
GOCCFLAGS 	= -v
GOLDFLAGS   = -ldflags="-s -w"
PREFIX		= /usr/local

PROGRAM_NAME=eosio-api-healthcheck
SOURCES=src/main.go src/server.go
DEPENDANCIES= github.com/firstrow/tcp_server \
	github.com/liamylian/jsontime/v2 \
	github.com/imroc/req \
	github.com/pborman/getopt/v2

.PHONY: all build/$(PROGRAM_NAME) clean
all: build
build: build/$(PROGRAM_NAME)

build/$(PROGRAM_NAME) : $(SOURCES)
	$(GO) build -o $@ $(GOCCFLAGS) $(GOLDFLAGS) $^
	$(GO) env > build/.buildinfo

deps:
	$(GO) get $(DEPENDANCIES)

info-file :
	echo PACKAGE_NAME=\"$(PROGRAM_NAME)\" "\n"\
	PACKAGE_DESCRIPTION=\"HAproxy healthcheck program for EOSIO API.\" "\n"\
	PACKAGE_VERSION=\"1.2\" "\n"\
	PACKAGE_PREFIX=\"$(PREFIX:/%=%)\" "\n"\
	PACKAGE_PROGRAM=\"build/$(PROGRAM_NAME)\" > scripts/pkg_info

package_deb: info-file
	./scripts/build.sh deb $(realpath build)

package_freebsd: info-file
	./scripts/build.sh freebsd $(realpath build)

clean:
	$(GO) clean
	$(RM) -rf build/
