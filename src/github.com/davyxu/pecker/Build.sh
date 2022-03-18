#!/usr/bin/env bash
CURR_DIR=`pwd`
cd ../../../..
export GOPATH=`pwd`
cd ${CURR_DIR}

go build -v -o pecker.exe github.com/davyxu/pecker

export GOARCH=amd64
export GOOS=linux
go build -v -o pecker github.com/davyxu/pecker