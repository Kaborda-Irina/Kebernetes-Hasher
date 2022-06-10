#!/usr/bin/env bash

echo "Building for Windows 32-bit"
env GOOS=windows GOARCH=386 go build -o internal/core/services/bin/sha256sum_32bit.exe cmd/main.go

echo "Building for Windows 64-bit"
env GOOS=windows GOARCH=amd64 go build -o  internal/core/services/bin/sha256sum_64bit.exe cmd/main.go

echo "Building for Linux 32-bit"
env GOOS=linux GOARCH=386 go build -o  internal/core/services/bin/sha256sum_32bit cmd/main.go

echo "Building for Linux 64-bit"
env GOOS=linux GOARCH=amd64 go build -o  internal/core/services/bin/sha256sum cmd/main.go

echo "Building for MacOS X 64-bit"
env GOOS=darwin GOARCH=amd64 go build -o  internal/core/services/bin/sha256sum_macos cmd/main.go