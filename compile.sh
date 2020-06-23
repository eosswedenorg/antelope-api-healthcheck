#!/bin/bash

SYSTEMS=( windows linux freebsd )

ARCHS=( 386 amd64 amd64p32 arm arm64 ppc ppc64 )

function usage() {
	echo "Usage: ${0##*/} [ -h|--help ] [ --target <system> ] [ -a|--arch <arch> ] [ -p|--package ]"
	echo ""
	echo " Valid systems:"
	for i in "${SYSTEMS[@]}"; do
		echo "  * ${i}"
	done
	echo ""
	echo " Valid architectures:"
	for i in "${ARCHS[@]}"; do
		echo "  * ${i}"
	done
	echo ""
	exit 1
}

options=$(getopt -n "${0##*/}" -o "ht:a:p" -l "help,target:,arch:,package" -- "$@")

[ $? -eq 0 ] || usage

eval set -- "$options"

MAKE_TARGET="all"

while true; do

	case $1 in
	-p|--package)
		MAKE_TARGET="package_deb"
		;;
	-t|--target)
		shift
		REGEX=$(echo "${SYSTEMS[@]}" | sed 's/[[:space:]]/|/g')
		[[ ! "$1" =~ ^($REGEX)$ ]] && {
			echo "Incorrect system '$1' provided"
			usage
		}
		export GOOS=$1
		;;
	-a|--arch)
		shift
		REGEX=$(echo "${ARCHS[@]}" | sed 's/[[:space:]]/|/g')
		[[ ! "$1" =~ ^($REGEX)$ ]] && {
			echo "Incorrect architecture '$1' provided"
			usage
		}
		export GOARCH=$1
		;;
	-h|--help)	usage ;;
	--)	shift
		break
		;;
	esac
	shift
done

MESSAGE=""
if [ ! -z "${GOOS}" ]; then
	# Hack to select the right package :)
	if [ "${MAKE_TARGET}" == "package_deb" ] && [ "${GOOS}" == "freebsd" ]; then
		MAKE_TARGET="package_freebsd"
	fi
	MESSAGE="[\e[34m::\e[0m] Crosscompiling for: ${GOOS}"
fi

if [ ! -z "${GOARCH}" ]; then
	MESSAGE="${MESSAGE} (${GOARCH})"
fi


[ ! -z "${MESSAGE}" ] && echo -e "" $MESSAGE

make -B ${MAKE_TARGET}
