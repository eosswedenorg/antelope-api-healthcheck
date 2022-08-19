#!/bin/bash

PROGRAM_ARCH=$(go env GOARCH)
TAR_FILENAME="${PROGRAM_NAME}-${PROGRAM_VERSION}-freebsd-${PROGRAM_ARCH}.tar.gz"

echo "Create archive: ${BUILDDIR}/${TAR_FILENAME}"
tar -C ${PKGROOT} --owner root --group root -zcvf ${BUILDDIR}/${TAR_FILENAME} .
