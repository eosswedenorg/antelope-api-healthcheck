#!/bin/bash
# Simple script to create a tar archive for FreeBSD

PACKAGE_TMPDIR="${PACKAGE_TMPDIR}/freebsd"
PACKAGE_RCDIR=/etc/rc.d
PACKAGE_NEWSYSLOGDIR=etc/newsyslog.conf.d

# Common variables
PID_FILE=/var/run/${PACKAGE_NAME}.pid

############################
#  Create rc file          #
############################

# rc does not like "-" in the filename.
RC_NAME=$(echo ${PACKAGE_NAME} | sed "s~-~_~g")

mkdir -p ${PACKAGE_TMPDIR}/${PACKAGE_RCDIR}
cat ${TEMPLATE_DIR}/rc.conf \
	| sed "s~{{ RC_NAME }}~${RC_NAME}~g" \
	| sed "s~{{ PID_FILE }}~${PID_FILE}~g" \
	| sed "s~{{ LOG_FILE }}~${PACKAGE_LOGFILE}~" \
	| sed "s~{{ DESCRIPTION }}~${PACKAGE_DESCRIPTION}~" \
	| sed "s~{{ PROGRAM }}~/${PACKAGE_BINDIR}/${PACKAGE_NAME}~" \
	> ${PACKAGE_TMPDIR}/${PACKAGE_RCDIR}/${RC_NAME}

# Must be executable.
chmod 755 ${PACKAGE_TMPDIR}/${PACKAGE_RCDIR}/${RC_NAME}

############################
#  Create newsyslog config #
############################

mkdir -p ${PACKAGE_TMPDIR}/${PACKAGE_NEWSYSLOGDIR}
cat ${TEMPLATE_DIR}/newsyslog.conf \
	| sed "s~{{ LOG_FILE }}~${PACKAGE_LOGFILE}~" \
	| sed "s~{{ PID_FILE }}~${PID_FILE}~g" \
	> ${PACKAGE_TMPDIR}/${PACKAGE_NEWSYSLOGDIR}/${PACKAGE_NAME}.conf


############################
#  Copy binary             #
############################

mkdir -p ${PACKAGE_TMPDIR}/${PACKAGE_BINDIR}
cp ${BASE_DIR}/../${PACKAGE_PROGRAM} ${PACKAGE_TMPDIR}/${PACKAGE_BINDIR}

############################
#  Create archive          #
############################

TAR_FILENAME="${PACKAGE_NAME}-${PACKAGE_VERSION}-freebsd.tar.gz"

tar -C ${PACKAGE_TMPDIR} --owner root --group root -zcvf ${BUILD_DIR}/${TAR_FILENAME} .
