# Inklay App Server

A tiny local web server that renders a greyscale PNG for the Inklay **Image URL** app.

- Output: 1872 × 1404 PNG, 8-bit greyscale (256 levels). The Inklay dithers it on the device.
- `http://<your-computer-ip>:8080/dashboard.png` is the URL for the Inklay app.
- `http://localhost:8080/` is a setup page with the URL, a copy icon and a preview (reload the page to refresh it). It opens automatically on start.
- Optional `?w=&h=` to change the size (e.g. portrait `?w=1404&h=1872`).

## Guide
Learn more in [Inklay Docs: Custom Inklay Apps w. Claude Code](https://inklay.gitbook.io/inklay-docs/custom-inklay-apps-w.-claude-code)