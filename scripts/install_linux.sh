#!/bin/bash
# Script to install program files on linux systems

source ${BASE_DIR}/functions/log_install.sh

SYSTEMDDIR=${DESTDIR}/lib/systemd/system
SYSTEMDLINKDIR=${DESTDIR}/etc/systemd/system
RSYSLOGDIR=${DESTDIR}/etc/rsyslog.d
LOGROTATEDIR=${DESTDIR}/etc/logrotate.d

# Create service file
log_install ${SYSTEMDDIR}/${PROGRAM_NAME}.service
mkdir -p ${SYSTEMDDIR}
cat ${TEMPLATE_DIR}/sysunit.service \
    | sed "s~{{ PROGRAM_NAME }}~${PROGRAM_NAME}~" \
    | sed "s~{{ PROGRAM_DESCRIPTION }}~${PROGRAM_DESCRIPTION}~" \
    | sed "s~{{ PROGRAM }}~${BINDIR}/${PROGRAM_NAME}~" \
    > ${SYSTEMDDIR}/${PROGRAM_NAME}.service

# Create systemd symlink
log_install ${SYSTEMDLINKDIR}/${PROGRAM_NAME}.service
mkdir -p ${SYSTEMDLINKDIR}
ln -s -T /lib/systemd/system/${PROGRAM_NAME}.service ${SYSTEMDLINKDIR}/${PROGRAM_NAME}.service

# Create systemd/init.d config file
log_install ${DESTDIR}/etc/sysconfig/${PROGRAM_NAME}
mkdir -p ${DESTDIR}/etc/sysconfig
cat ${TEMPLATE_DIR}/config \
    | sed "s~{{ PROGRAM_NAME }}~${PROGRAM_NAME}~" \
    > ${DESTDIR}/etc/sysconfig/${PROGRAM_NAME}

# Create rsyslog file
log_install ${RSYSLOGDIR}/49-${PROGRAM_NAME}.conf
mkdir -p ${RSYSLOGDIR}
cat ${TEMPLATE_DIR}/rsyslog.conf \
    | sed "s~{{ PROGRAM }}~${PROGRAM_NAME}~" \
    | sed "s~{{ LOG_FILE }}~${LOGFILE}~" \
    > ${RSYSLOGDIR}/49-${PROGRAM_NAME}.conf

# Create logrotate file
log_install ${LOGROTATEDIR}/${PROGRAM_NAME}.conf
mkdir -p ${LOGROTATEDIR}
cat ${TEMPLATE_DIR}/logrotate.conf \
    | sed "s~{{ LOG_FILE }}~${LOGFILE}~" \
    > ${LOGROTATEDIR}/${PROGRAM_NAME}.conf
chmod 644 ${LOGROTATEDIR}/${PROGRAM_NAME}.conf

# Copy program
log_install ${DESTDIR}${SHAREDIR}
mkdir -p ${DESTDIR}/${BINDIR}
cp ${BUILD_DIR}/${PROGRAM_NAME} ${DESTDIR}${BINDIR}/${PROGRAM_NAME}

# Copy files.
mkdir -p ${DESTDIR}${SHAREDIR}
cp ${BASE_DIR}/../README.md ${DESTDIR}${SHAREDIR}
