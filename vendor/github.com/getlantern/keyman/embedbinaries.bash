#!/bin/bash

###############################################################################
#
# This script regenerates the source files that embed the platform-specific
# executables.
#
###############################################################################

function die() {
  echo $*
  exit 1
}

if [ -z "$BNS_CERT" ] || [ -z "$BNS_CERT_PASS" ]
then
	die "$0: Please set BNS_CERT and BNS_CERT_PASS to the bns_cert.p12 signing key and the password for that key"
fi

osslsigncode sign -pkcs12 "$BNS_CERT" -pass "$BNS_CERT_PASS" -in certimporter/Release/certimporter.exe -out binaries/windows/certimporter.exe || die "Could not sign windows"
go-bindata -nomemcopy -nocompress -prefix binaries/windows -o ./certimporter/certimporter.go -pkg certimporter binaries/windows
