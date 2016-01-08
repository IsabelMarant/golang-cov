#!/usr/bin/bash

WORKROOT=$(pwd) 

# prepare PATH, GOROOT and GOPATH
#export PATH=$(pwd)/go/bin:$PATH
#export GOROOT=$(pwd)/go
export GOPATH=$(pwd)

# go to /bfe
oldpath=`pwd`
cd ../../../
export GOPATH=$(pwd)/bfe-common/golang-lib:$GOPATH

cd $oldpath

# build 
cd src/main
go build -o coverage
if [ $? -ne 0 ];
then
    echo "fail to go build coverage.go"
    exit 1
fi

echo "OK for go build go_bfe.go"

cd ../..
# create directory for output
if [ -d "./output" ]
then
    rm -rf output
fi

mkdir output

# copy main conf 
cp -r conf output/
cp -r cover_output output/

mkdir output/bin
cp src/main/coverage output/bin 

echo "OK for build coverage"
