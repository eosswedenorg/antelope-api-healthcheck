#!/bin/bash

BUILD_INFO=build/.buildinfo

if [ ! -f "${BUILD_INFO}" ]; then
    echo "Could not find '${BUILD_INFO}' file, You need to compile first."
    exit 1
fi

TYPE=$(cat "${BUILD_INFO}" | sed -n 's/^GOOS=\"\(.*\)\"/\1/p')
if [ "$TYPE" == "freebsd" ]; then
    MAKE_TARGET="package_freebsd"
else
    MAKE_TARGET="package_deb"
fi

make -B ${MAKE_TARGET}
