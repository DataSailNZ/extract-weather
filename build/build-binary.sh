#!/bin/sh

set -x
cd ../api

go get -t -d -v .

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o "../build/extract-weather"
