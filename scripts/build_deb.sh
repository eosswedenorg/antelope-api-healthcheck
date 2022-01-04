#!/usr/bin/env bash

PACKAGE_SYSUNITDIR=etc/systemd/system
PACKAGE_RSYSLOGDIR=etc/rsyslog.d
PACKAGE_LOGROTATEDIR=etc/logrotate.d

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

rm -fr ${PACKAGE_TMPDIR}

# Create debian files.
mkdir -p ${PACKAGE_TMPDIR}/DEBIAN
echo "Package: ${PACKAGE_NAME}
Version: ${PACKAGE_VERSION}-${RELEASE}
Section: introspection
Priority: optional
Architecture: amd64
Homepage: https://github.com/eosswedenorg/eos-api-healthcheck
Maintainer: Henrik Hautakoski <henrik@eossweden.org>
Description: ${PACKAGE_DESCRIPTION}" &> ${PACKAGE_TMPDIR}/DEBIAN/control

cat ${PACKAGE_TMPDIR}/DEBIAN/control

# Create service file
mkdir -p ${PACKAGE_TMPDIR}/${PACKAGE_SYSUNITDIR}
cat ${TEMPLATE_DIR}/sysunit.service \
    | sed "s~{{ PACKAGE_NAME }}~${PACKAGE_NAME}~" \
    | sed "s~{{ DESCRIPTION }}~${PACKAGE_DESCRIPTION}~" \
    | sed "s~{{ PROGRAM }}~/${PACKAGE_PREFIX}/bin/${PACKAGE_NAME}~" \
    > ${PACKAGE_TMPDIR}/${PACKAGE_SYSUNITDIR}/${PACKAGE_NAME}.service

# Create rsyslog file
mkdir -p ${PACKAGE_TMPDIR}/${PACKAGE_RSYSLOGDIR}
cat ${TEMPLATE_DIR}/rsyslog.conf \
    | sed "s~{{ PROGRAM }}~${PACKAGE_NAME}~" \
    | sed "s~{{ LOG_FILE }}~${PACKAGE_LOGFILE}~" \
    > ${PACKAGE_TMPDIR}/${PACKAGE_RSYSLOGDIR}/49-${PACKAGE_NAME}.conf

# Create logrotate file
mkdir -p ${PACKAGE_TMPDIR}/${PACKAGE_LOGROTATEDIR}
cat ${TEMPLATE_DIR}/logrotate.conf \
    | sed "s~{{ LOG_FILE }}~${PACKAGE_LOGFILE}~" \
    > ${PACKAGE_TMPDIR}/${PACKAGE_LOGROTATEDIR}/${PACKAGE_NAME}.conf
chmod 644 ${PACKAGE_TMPDIR}/${PACKAGE_LOGROTATEDIR}/${PACKAGE_NAME}.conf

# Cerate config file
mkdir -p ${PACKAGE_TMPDIR}/${PACKAGE_ETCDIR}
cat ${TEMPLATE_DIR}/config \
    | sed "s~{{ PACKAGE_NAME }}~${PACKAGE_NAME}~" \
    > ${PACKAGE_TMPDIR}/${PACKAGE_ETCDIR}/env

# Copy program
mkdir -p ${PACKAGE_TMPDIR}/${PACKAGE_BINDIR}
cp ${BASE_DIR}/../${PACKAGE_PROGRAM} ${PACKAGE_TMPDIR}/${PACKAGE_BINDIR}/${PACKAGE_NAME}

# Copy files.
mkdir -p ${PACKAGE_TMPDIR}/${PACKAGE_SHAREDIR}
cp ${BASE_DIR}/../README.md ${PACKAGE_TMPDIR}/${PACKAGE_SHAREDIR}

fakeroot dpkg-deb --build ${PACKAGE_TMPDIR} ${BUILD_DIR}/${PACKAGE_FULLNAME}.deb
