#!/bin/bash

cd `dirname $0`
cd ../cmd/sample

OutputLinuxPath="../../bin/linux"
OutputMacPath="../../bin/mac"
OutputWindowsPath="../../bin/win"

mkdir -p $OutputLinuxPath
mkdir -p $OutputMacPath
mkdir -p $OutputWindowsPath

GOOS=linux GOARCH=amd64 go build -o $OutputLinuxPath
GOOS=darwin GOARCH=arm64 go build -o $OutputMacPath
GOOS=windows GOARCH=amd64 go build -o $OutputWindowsPath
