#!/usr/bin/env bash

BASE_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

PACKAGE_BINDIR=${PACKAGE_PREFIX}/bin
PACKAGE_ETCDIR=etc/${PACKAGE_NAME}
PACKAGE_SYSUNITDIR=etc/systemd/system
PACKAGE_SHAREDIR=${PACKAGE_PREFIX}/share/${PACKAGE_NAME}
PACKAGE_DESCRIPTION="HAproxy healthcheck program for EOSIO API."
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
mkdir -p ${BASE_DIR}/${PACKAGE_TMPDIR}/${PACKAGE_SYSUNITDIR}
cat ${BASE_DIR}/template.service \
	| sed "s~{{ PACKAGE_NAME }}~${PACKAGE_NAME}~" \
	| sed "s~{{ DESCRIPTION }}~${PACKAGE_DESCRIPTION}~" \
	| sed "s~{{ PROGRAM }}~/${PACKAGE_PREFIX}/bin/${PACKAGE_NAME}~" \
	> ${BASE_DIR}/${PACKAGE_TMPDIR}/${PACKAGE_SYSUNITDIR}/${PACKAGE_NAME}.service

# Cerate config file
mkdir -p ${BASE_DIR}/${PACKAGE_TMPDIR}/${PACKAGE_ETCDIR}
cat ${BASE_DIR}/config \
	| sed "s~{{ PACKAGE_NAME }}~${PACKAGE_NAME}~" \
	> ${BASE_DIR}/${PACKAGE_TMPDIR}/${PACKAGE_ETCDIR}/env

# Copy program
mkdir -p ${BASE_DIR}/${PACKAGE_TMPDIR}/${PACKAGE_BINDIR}
cp ${BASE_DIR}/../${PACKAGE_PROGRAM} ${BASE_DIR}/${PACKAGE_TMPDIR}/${PACKAGE_BINDIR}/${PACKAGE_NAME}

# Copy files.
mkdir -p ${BASE_DIR}/${PACKAGE_TMPDIR}/${PACKAGE_SHAREDIR}
cp ${BASE_DIR}/../README.md ${BASE_DIR}/${PACKAGE_TMPDIR}/${PACKAGE_SHAREDIR}

fakeroot dpkg-deb --build ${BASE_DIR}/${PACKAGE_TMPDIR} ${BASE_DIR}/${PACKAGE_FULLNAME}.deb
