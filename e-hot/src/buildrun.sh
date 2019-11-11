#!/bin/sh
echo Compiling...
GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -o binary -ldflags '-w' -mod=vendor
./binary