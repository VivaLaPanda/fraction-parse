#!/bin/bash

mkdir -p tmp/src/
ln -s $PWD tmp/src/fraction-parse
export GOPATH=$PWD/tmp
go build -o exe fraction-parse
rm -rf tmp

