#!/bin/bash
GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -o binary -ldflags '-w' -mod=vendor
upx binary -1
cp binary ../../../f-localhot-binary/bin/binary

