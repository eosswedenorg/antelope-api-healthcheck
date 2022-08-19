#!/bin/bash

if [ $# -lt 1 ]; then
    echo "$0 <version>"
    exit 1
fi

sed -i "s:PROGRAM_VERSION=\(.*\):PROGRAM_VERSION=$1:g" scripts/info.sh
sed -i "s~\print(\"Version:\ v\(.*\)\\\n\")~print(\\\"Version:\ v$1\\\n\\\")~g" src/main.go
