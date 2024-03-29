#!/bin/bash
# Simple script to install program files on FreeBSD systems

source ${BASE_DIR}/functions/log_install.sh

RCDIR=${DESTDIR}/etc/rc.d
NEWSYSLOGDIR=${DESTDIR}/etc/newsyslog.conf.d

# Common variables
PID_FILE=/var/run/${PROGRAM_NAME}.pid
# rc does not like "-" in the filename.
RC_NAME=$(echo ${PROGRAM_NAME} | sed "s~-~_~g")

############################
#  Create rc file          #
############################

log_install ${RCDIR}/${RC_NAME}
mkdir -p ${RCDIR}
cat ${TEMPLATE_DIR}/rc.conf \
	| sed "s~{{ RC_NAME }}~${RC_NAME}~g" \
	| sed "s~{{ PID_FILE }}~${PID_FILE}~g" \
	| sed "s~{{ LOG_FILE }}~${LOGFILE}~" \
	| sed "s~{{ PROGRAM_DESCRIPTION }}~${PROGRAM_DESCRIPTION}~" \
	| sed "s~{{ PROGRAM }}~${BINDIR}/${PROGRAM_NAME}~" \
	> ${RCDIR}/${RC_NAME}

# Must be executable.
chmod 755 ${RCDIR}/${RC_NAME}

############################
#  Create newsyslog config #
############################

log_install ${NEWSYSLOGDIR}/${PROGRAM_NAME}.conf
mkdir -p ${NEWSYSLOGDIR}
cat ${TEMPLATE_DIR}/newsyslog.conf \
	| sed "s~{{ LOG_FILE }}~${LOGFILE}~" \
	| sed "s~{{ PID_FILE }}~${PID_FILE}~g" \
	> ${NEWSYSLOGDIR}/${PROGRAM_NAME}.conf


############################
#  Copy binary             #
############################

log_install ${DESTDIR}${BINDIR}/${PROGRAM_NAME}
mkdir -p ${DESTDIR}/${BINDIR}
cp ${BUILD_DIR}/${PROGRAM_NAME} ${DESTDIR}${BINDIR}/${PROGRAM_NAME}
