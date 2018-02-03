#!/bin/bash

mkdir -p tmp/src/
ln -s $PWD tmp/src/fraction-parse
GOPATH=$PWD/tmp
go build fraction-parse
rm -rf tmp

