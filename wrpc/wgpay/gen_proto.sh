#!/usr/bin/env bash

# Copyright (c) 2019-2029 Dunyu All Rights Reserved.
#
# Author      : yangping
# Email       : ping.yang@wengold.net
# Version     : 1.0.0
#
# Prismy.No | Date       | Modified by. | Description
# -------------------------------------------------------------------
# 00001       2022/10/213  yangping       New version
# -------------------------------------------------------------------

bin=`dirname "$0"`
bin=`cd "$bin"; pwd`

cd ${bin}
rm -rf ./proto/*.pb.go
protoc --go_out=. --go-grpc_out=. ./proto/wgpay.proto

echo "Finished generat wgpay proto files:"
ls ./proto/*

