#!/bin/bash

echo "Building for Linux"
GOOS="linux" GOOS="linux" go build -o jsonscraper_linux

echo "Building for Darwin"
GOOS="linux" GOOS="darwin" go build -o jsonscraper_darwin

echo "Building for Windows"
GOOS="linux" GOOS="windows" go build -o jsonscraper_windows.exe