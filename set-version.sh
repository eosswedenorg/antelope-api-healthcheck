#!/bin/bash
# Simple script to make it easy to update the version number for the program.
#
# Debian
# ----------------------------
# For releasing debian packages, there must be a name and email associated with the version.
# You can pass "-n|--name" and "-e|--email" as parameters to this script.
#
# You can if you want set the following enviroment variables in your shell to have your name and email be inserted without cli flags.
#  DEB_MAINT_NAME
#  DEB_MAINT_EMAIL

function usage() {
    echo "Usage: ${0##*/} [ -h|--help ] [ -n|--name <value> ] [ -e|--email <value> ] [ --nodebchanges ] <version>"
    exit 1
}

eval set -- "$(getopt -n "${0##*/}" -o "hn:e:" -l "help,name:,email:,nodebchanges" -- "$@")"

WRITE_DEBCHANGES=1
while true; do

    case $1 in
    -n|--name)
        shift
        DEB_MAINT_NAME=$1
        ;;
    -e|--email)
        shift
        DEB_MAINT_EMAIL=$1
        ;;
    --nodebchanges)
        WRITE_DEBCHANGES=0
        ;;
    -h|--help) usage ;;
    --) shift
        break
        ;;
    esac
    shift
done

[ $# -gt 0 ] || [ $? -eq 0 ] || usage

VERSION=$@

if [ ${WRITE_DEBCHANGES} -ne 0 ]; then
    # Update debian changelog
    ex debian/changelog <<EOF
1 insert
eosio-api-healthcheck (${VERSION}) unstable; urgency=medium

  *

 -- ${DEB_MAINT_NAME} <${DEB_MAINT_EMAIL}>  $(date -R)

.
xit
EOF
    echo -e "[\e[34m::\e[0m] Inserted template in \e[1mdebian/changelog\e[0m. \e[33mMake sure you edit this file with the actual changes!\e[0m"
else :
    echo -e "[\e[33m::\e[0m] Skipping \e[1mdebian/changelog\e[0m."
fi


# Update Makefile
sed -i "s:PROGRAM_VERSION\(\s*\)=\(\s*\)\(.*\):PROGRAM_VERSION\1=\2$VERSION:g" Makefile
echo -e "[\e[34m::\e[0m] Set PROGRAM_VERSION=\e[34m$VERSION\e[0m in \e[1mMakefile\e[0m"
