#!/usr/bin/env bash

XC_OS=$(go env GOOS)
XC_ARCH=$(go env GOARCH)
DEST_BIN=terraform-provider-infoblox

echo "Compiling for OS: $XC_OS and ARCH: $XC_ARCH"

gox -os="${XC_OS}" -arch="${XC_ARCH}"

if [ $? != 0 ] ; then
    echo "Failed to compile, bailing."
    exit 1
fi

echo "Looking for Terraform install"

TERRAFORM=$(which terraform)

[ $TERRAFORM ] && TERRAFORM_LOC=$(dirname ${TERRAFORM})

if [ $TERRAFORM_LOC  ] ; then
    DEST_PATH=$TERRAFORM_LOC
else
    DEST_PATH=$GOPATH/bin
fi

echo ""
echo "Moving terraform-provider-infoblox_${XC_OS}_${XC_ARCH} to $DEST_PATH/$DEST_BIN"
echo ""

mv terraform-provider-infoblox_${XC_OS}_${XC_ARCH} $DEST_PATH/$DEST_BIN

echo "Resulting binary: "
echo ""
echo $(ls -la $DEST_PATH/$DEST_BIN)
