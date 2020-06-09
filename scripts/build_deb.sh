#!/usr/bin/env bash

PACKAGE_TMPDIR="${PACKAGE_TMPDIR}/debian"

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

PACKAGE_FULLNAME="${PACKAGE_NAME}_${PACKAGE_VERSION}-${RELEASE}-${DISTRIB_ID}-${DISTRIB_RELEASE}_amd64"

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

# Create rsyslog file
mkdir -p ${BASE_DIR}/${PACKAGE_TMPDIR}/${PACKAGE_RSYSLOGDIR}
cat ${BASE_DIR}/rsyslog-template.conf \
	| sed "s~{{ PROGRAM }}~${PACKAGE_NAME}~" \
	| sed "s~{{ LOG_FILE }}~${PACKAGE_LOGFILE}~" \
	> ${BASE_DIR}/${PACKAGE_TMPDIR}/${PACKAGE_RSYSLOGDIR}/49-${PACKAGE_NAME}.conf

# Create logrotate file
mkdir -p ${BASE_DIR}/${PACKAGE_TMPDIR}/${PACKAGE_LOGROTATEDIR}
cat ${BASE_DIR}/logrotate-template.conf \
	| sed "s~{{ LOG_FILE }}~${PACKAGE_LOGFILE}~" \
	> ${BASE_DIR}/${PACKAGE_TMPDIR}/${PACKAGE_LOGROTATEDIR}/${PACKAGE_NAME}.conf
chmod 644 ${BASE_DIR}/${PACKAGE_TMPDIR}/${PACKAGE_LOGROTATEDIR}/${PACKAGE_NAME}.conf

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
