#!/bin/bash
set -e

banner() {
echo "================================================================"
echo "$1"
echo "================================================================"
}

PROJECT="terraform-provider-infoblox"

banner "Building docker image..."
docker build -t $PROJECT .
banner "Copying the binary..."
docker run -v ${PWD}/bin:/out:rw $PROJECT
