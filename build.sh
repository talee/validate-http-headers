#!/bin/sh
rm -rf dist
mkdir dist
echo 'Building for Linux'
env GOOS=linux GOARCH=amd64 go build -o dist/validate-http-headers-linux
echo 'Building for OSX'
env GOOS=darwin GOARCH=amd64 go build -o dist/validate-http-headers-osx
