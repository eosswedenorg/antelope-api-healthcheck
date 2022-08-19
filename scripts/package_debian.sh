#!/bin/bash

if [[ -f /etc/upstream-release/lsb-release ]]; then
    source /etc/upstream-release/lsb-release
elif [[ -f /etc/lsb-release ]]; then
    source /etc/lsb-release
else
    echo "ERROR: could not determine debian release."
    exit 1
fi

DISTRIB_ID=$(echo $DISTRIB_ID | tr '[:upper:]' '[:lower:]')

# Default to 1 if no release is set.
if [[ -z $RELEASE ]]; then
    RELEASE="1"
fi

PACKAGE_FULLNAME="${PROGRAM_NAME}_${PROGRAM_VERSION}-${RELEASE}-${DISTRIB_ID}-${DISTRIB_RELEASE}_amd64"

# Create debian files.
mkdir -p ${PKGROOT}/DEBIAN
echo "Package: ${PROGRAM_NAME}
Version: ${PROGRAM_VERSION}-${RELEASE}
Section: introspection
Priority: optional
Architecture: amd64
Homepage: https://github.com/eosswedenorg/eos-api-healthcheck
Maintainer: Henrik Hautakoski <henrik@eossweden.org>
Description: ${PROGRAM_DESCRIPTION}" &> ${PKGROOT}/DEBIAN/control

cat ${PKGROOT}/DEBIAN/control

fakeroot dpkg-deb --build ${PKGROOT} ${BUILDDIR}/${PACKAGE_FULLNAME}.deb
