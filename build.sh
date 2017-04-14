#!/bin/bash

echo "Building for Linux"
GOARCH="amd64" GOOS="linux" go build -o jsonscraper_linux

echo "Building for Darwin"
GOARCH="amd64" GOOS="darwin" go build -o jsonscraper_darwin

echo "Building for Windows"
GOARCH="amd64" GOOS="windows" go build -o jsonscraper_windows.exe