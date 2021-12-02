#!/bin/bash

BINARY=build/eosio-api-healthcheck

if [ ! -f "${BINARY}" ]; then
	echo "Could not find '${BINARY}', You need to compile first."
	exit 1
fi

# Bit of a hack to figure out if we need to package for FreeBSD or not.
if [ -n "$(file $BINARY | grep 'FreeBSD')" ]; then
	MAKE_TARGET="package_freebsd"
else
	MAKE_TARGET="package_deb"
fi

make -B ${MAKE_TARGET}
