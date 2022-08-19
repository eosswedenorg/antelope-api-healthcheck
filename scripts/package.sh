#!/bin/bash

############################
#  Read cli args           #
############################

if [ $# -lt 1 ]; then
    echo "$0 <type>"
    exit 1
fi

PKG_TYPE=$1

# Setup vars
export BASE_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

# export
export PKGROOT
export BUILDDIR

############################
#  Run script              #
############################

SCRIPT="${BASE_DIR}/package_${PKG_TYPE}.sh"

# Check and call script
if [ ! -x $SCRIPT ]; then
	echo "$SCRIPT not found"
	exit 1
fi

# Export info variables
set -o allexport
source ${BASE_DIR}/info.sh
set +o allexport

bash $SCRIPT
