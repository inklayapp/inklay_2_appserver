# Inklay App Server

A tiny local web server that renders a greyscale PNG for the Inklay **Image URL** app.

- Output: 1872 × 1404 PNG, 8-bit greyscale (256 levels). The Inklay dithers it on the device.
- `http://<your-computer-ip>:8080/dashboard.png` is the URL for the Inklay app.
- `http://localhost:8080/` is a setup page with the URL, a copy icon and a preview (reload the page to refresh it). It opens automatically on start.
- Optional `?w=&h=` to change the size (e.g. portrait `?w=1404&h=1872`).

## For users
1. Download the file for your computer from `dist/`.
2. Double-click it. A browser tab opens showing your Inklay URL.
3. Paste the URL into the Image URL app. Keep the computer on and on the same Wi-Fi.

First-start warnings (unsigned builds):
- **Windows:** "Windows protected your PC" → *More info* → *Run anyway*. Then allow the firewall prompt on **private networks**, otherwise the Inklay cannot reach it.
- **macOS:** "cannot be opened" → System Settings → Privacy & Security → *Open Anyway*. Mac may also ask to allow incoming connections: allow.

## For developers
```
go run .                      # start locally
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/InklayAppServer-windows.exe .
GOOS=darwin  GOARCH=arm64 go build -ldflags="-s -w" -o dist/InklayAppServer-mac-applesilicon .
```
Edit `renderDashboard` in `dashboard.go` to change what is drawn. Images go in `assets/`, fonts in `fonts/` (default: Roboto Mono). Both are built into the app, so no files are needed at runtime.

## Customize it with Claude
Open this folder in Claude (Claude Code or the Claude desktop app) and describe what you want, for example
"show today's weather for Zurich and my next calendar events". `CLAUDE.md` tells Claude the Inklay display rules
and how to rebuild. After changes, run `./build.sh` (Claude does this for you) and double-click the start file again.
