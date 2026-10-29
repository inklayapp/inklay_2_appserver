#!/bin/bash
# Double-click on a Mac to start the Inklay App Server.
cd "$(dirname "$0")"
if [ "$(uname -m)" = "arm64" ]; then BIN=dist/InklayAppServer-mac-applesilicon; else BIN=dist/InklayAppServer-mac-intel; fi
chmod +x "$BIN"
exec "./$BIN"
