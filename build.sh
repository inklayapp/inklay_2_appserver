#!/bin/bash
# Rebuilds the double-click apps for every system into dist/.
# Needs Go (https://go.dev/dl/ or `brew install go`).
set -e
cd "$(dirname "$0")"
if ! command -v go >/dev/null; then
  echo "Go is not installed. Install it from https://go.dev/dl/ (or: brew install go) and run again."
  exit 1
fi
mkdir -p dist
build() { echo "Building $3"; GOOS=$1 GOARCH=$2 go build -ldflags="-s -w" -o "dist/$3" .; }
build darwin  arm64 InklayAppServer-mac-applesilicon
build darwin  amd64 InklayAppServer-mac-intel
build windows amd64 InklayAppServer-windows.exe
build linux   amd64 InklayAppServer-linux-amd64
build linux   arm64 InklayAppServer-linux-arm64
chmod +x dist/*
echo "Done. Double-click 'Start Inklay App Server.command' (Mac) or dist/InklayAppServer-windows.exe (Windows)."
