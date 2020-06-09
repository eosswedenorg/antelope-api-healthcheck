
GO			= go
GOCCFLAGS 	= -v
GOLDFLAGS   =
PREFIX		= /usr/local

PROGRAM_NAME=eosio-api-healthcheck
SOURCES=src/server.go
DEPENDANCIES= github.com/firstrow/tcp_server \
	github.com/liamylian/jsontime/v2 \
	github.com/imroc/req

all: build
build: build/$(PROGRAM_NAME)

build/$(PROGRAM_NAME) : $(SOURCES)
	$(GO) build -o $@ $(GOCCFLAGS) $(GOLDFLAGS) $<

deps:
	$(GO) get $(DEPENDANCIES)

package_info : scripts/info
	echo PACKAGE_NAME=\"$(PROGRAM_NAME)\" "\n"\
	PACKAGE_DESCRIPTION=\"HAproxy healthcheck program for EOSIO API.\" "\n"\
	PACKAGE_VERSION=\"0.3.1\" "\n"\
	PACKAGE_PREFIX=\"$(PREFIX:/%=%)\" "\n"\
	PACKAGE_PROGRAM=\"build/$(PROGRAM_NAME)\" > $<

package : package_info build

package_deb: package
	./scripts/build.sh deb

clean:
	$(GO) clean
	$(RM) -rf build/
