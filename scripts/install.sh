#!/bin/bash

############################
#  Read cli args           #
############################

if [ $# -lt 1 ]; then
    echo "$0 <type>"
    exit 1
fi

SYSTEM_TYPE=$1

if [ $# -gt 1 ]; then
    PREFIX=$2
fi

if [ $# -gt 2 ]; then
    DESTDIR=$3
fi

############################
#  Exported variables.     #
############################

export BASE_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
export TEMPLATE_DIR=${BASE_DIR}/templates
export BUILD_DIR=${BASE_DIR}/../build

# Export info variables
set -o allexport
source ${BASE_DIR}/info.sh
set +o allexport

# Directories.
export DESTDIR
export PREFIX
export BINDIR=${PREFIX}/bin
export ETCDIR=/etc/${PROGRAM_NAME}
export LOGDIR=/var/log
export LOGFILE=${LOGDIR}/${PROGRAM_NAME}.log
export SHAREDIR=${PREFIX}/share/${PROGRAM_NAME}

############################
#  Run script              #
############################

SCRIPT="${BASE_DIR}/install_${SYSTEM_TYPE}.sh"

# Check and call script
if [ ! -x $SCRIPT ]; then
	echo "$SCRIPT not found"
	exit 1
fi

echo -e "[\e[34m::\e[0m] Installing for system: \e[32m${SYSTEM_TYPE}\e[0m"
if [ -n "$DESTDIR" ]; then
    echo -e "[\e[34m::\e[0m] Installing with root:  \e[32m${DESTDIR}\e[0m"
fi

bash $SCRIPT
