#!/usr/bin/env bash

BASE_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

PACKAGE_DESCRIPTION="HAproxy healthcheck program for EOS API."
PACKAGE_TMPDIR="pack"

# Default to 1 if no release is set.
if [[ -z $RELEASE ]]; then
  RELEASE="1"
fi

PACKAGE_FULLNAME="${PACKAGE_NAME}_${PACKAGE_VERSION}-${RELEASE}_amd64"

rm -fr ${BASE_DIR}/${PACKAGE_TMPDIR}

# Create debian files.
mkdir -p ${BASE_DIR}/${PACKAGE_TMPDIR}/DEBIAN
echo "Package: ${PACKAGE_NAME}
Version: ${PACKAGE_VERSION}-${RELEASE}
Section: introspection
Priority: optional
Architecture: amd64
Homepage: https://github.com/eosswedenorg/eos-api-healthcheck
Maintainer: Henrik Hautakoski <henrik@eossweden.org>
Description: ${PACKAGE_DESCRIPTION}" &> ${BASE_DIR}/${PACKAGE_TMPDIR}/DEBIAN/control

cat ${BASE_DIR}/${PACKAGE_TMPDIR}/DEBIAN/control

# Create service file
mkdir -p ${BASE_DIR}/${PACKAGE_TMPDIR}/etc/systemd/system
cat ${BASE_DIR}/template.service \
	| sed "s~{{ DESCRIPTION }}~${PACKAGE_DESCRIPTION}~" \
	| sed "s~{{ PROGRAM }}~/${PACKAGE_PREFIX}/bin/${PACKAGE_NAME}~" \
	> ${BASE_DIR}/${PACKAGE_TMPDIR}/etc/systemd/system/${PACKAGE_NAME}.service

# Copy program
mkdir -p ${BASE_DIR}/${PACKAGE_TMPDIR}/${PACKAGE_PREFIX}/bin
cp ${BASE_DIR}/../${PACKAGE_PROGRAM} ${BASE_DIR}/${PACKAGE_TMPDIR}/${PACKAGE_PREFIX}/bin/${PACKAGE_NAME}

fakeroot dpkg-deb --build ${BASE_DIR}/${PACKAGE_TMPDIR} ${BASE_DIR}/${PACKAGE_FULLNAME}.deb
