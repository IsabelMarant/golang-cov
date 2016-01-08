#!/usr/bin/bash
# dir structure
# /bfe
#     /bfe-common
#         /go
#             /output
#         /golang-lib
#     /go-bfe
#         /go-bfe
#             build.sh
#

WORKROOT=$(pwd) 

# prepare tmp dir for go build
TMPDIR="$WORKROOT/tmp"
mkdir $TMPDIR
export TMPDIR=$TMPDIR 

# generate protobuf file
SUBDIRS="src/bfe_modules/mod_access_pb/bfe_access_pb \
        src/bfe_modules/mod_waf_client/waf_pb \
        src/bfe_modules/mod_dict_client/dict_pb"
for i in ${SUBDIRS}
do
    make -C $i
    if [ $? != 0 ] 
    then
        echo "generate gogoprotobuf $i failed"
        exit 1
    fi
done
echo "OK for generate gogoprotobuf"

# restore working dir
cd ${WORKROOT}
# unzip go environment
tar zxvf ../../bfe-common/go/output/go.data >/dev/null
if [ $? -ne 0 ];
then
    echo "fail in extract go"
    exit 1
fi

echo "OK for extract go"

# prepare PATH, GOROOT and GOPATH
export PATH=$(pwd)/go/bin:$PATH
export GOROOT=$(pwd)/go
export GOPATH=$(pwd)

make -C "src/bfe_basic/condition/parser"
if [ $? != 0 ] 
then
    echo "generate parser failed"
    exit 1
fi

# go to /bfe
oldpath=`pwd`
cd ../..
export GOPATH=$(pwd)/bfe-common/golang-lib:$GOPATH

cd $oldpath

# first make deplib
make -C src/bfe_modules/mod_decrypt_uri/enc/c
if [ $? -ne 0 ];
then
    echo "fail in make decrypt lib"
    exit 1
fi
libdir=`pwd`/src/bfe_modules/mod_decrypt_uri/enc/c
echo "libdir $libdir"
export CGO_CFLAGS="-I$libdir"
export CGO_LDFLAGS="-L$libdir -lstdc++ -lcrypto -lenc_for_go"

# set cgo flags for mod_crypto
export CGO_CFLAGS="-I`pwd`/../../bfe-common/openssl/include $CGO_CFLAGS"
export CGO_LDFLAGS="`pwd`/../../bfe-common/openssl/lib/libcrypto.a $CGO_LDFLAGS"

# run go test for all subdirectory
#cd src && go test -cover ./... 
sh coverage.sh 
if [ $? -ne 0 ];
then
    echo "go test failed"
    exit 1
fi
echo "OK for go test"
