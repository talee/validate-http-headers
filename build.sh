#!/bin/sh
env GOOS=linux GOARCH=amd64 go build -o validate-http-headers-linux
env GOOS=darwin GOARCH=amd64 go build -o validate-http-headers-osx
