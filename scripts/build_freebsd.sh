#!/bin/bash
# Simple script to create a tar archive for FreeBSD

PACKAGE_TMPDIR="${PACKAGE_TMPDIR}/freebsd"
PACKAGE_RCDIR=/etc/rc.d
PACKAGE_NEWSYSLOGDIR=etc/newsyslog.conf.d

############################
#  Create rc file          #
############################

# rc does not like "-" in the filename.
RC_NAME=$(echo ${PACKAGE_NAME} | sed "s~-~_~g")

mkdir -p ${BASE_DIR}/${PACKAGE_TMPDIR}/${PACKAGE_RCDIR}
cat ${BASE_DIR}/rc.template \
	| sed "s~{{ RC_NAME }}~${RC_NAME}~g" \
	| sed "s~{{ DESCRIPTION }}~${PACKAGE_DESCRIPTION}~" \
	| sed "s~{{ PROGRAM }}~/${PACKAGE_BINDIR}/${PACKAGE_NAME}~" \
	> ${BASE_DIR}/${PACKAGE_TMPDIR}/${PACKAGE_RCDIR}/${RC_NAME}

############################
#  Create newsyslog config #
############################

mkdir -p ${BASE_DIR}/${PACKAGE_TMPDIR}/${PACKAGE_NEWSYSLOGDIR}
cat ${TEMPLATE_DIR}/newsyslog.conf \
	| sed "s~{{ LOG_FILE }}~${PACKAGE_LOGFILE}~" \
	> ${BASE_DIR}/${PACKAGE_TMPDIR}/${PACKAGE_NEWSYSLOGDIR}/${PACKAGE_NAME}.conf


############################
#  Copy binary             #
############################

mkdir -p ${BASE_DIR}/${PACKAGE_TMPDIR}/${PACKAGE_BINDIR}
cp ${BASE_DIR}/../${PACKAGE_PROGRAM} ${BASE_DIR}/${PACKAGE_TMPDIR}/${PACKAGE_BINDIR}

############################
#  Create archive          #
############################

TAR_FILENAME="${PACKAGE_NAME}-${PACKAGE_VERSION}-freebsd.tar.gz"

tar -C ${BASE_DIR}/${PACKAGE_TMPDIR} --owner root --group root -zcvf ${BASE_DIR}/${TAR_FILENAME} .
