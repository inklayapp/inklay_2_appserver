# Inklay App Server

A small Go web server that draws a dashboard as a greyscale PNG for the **Inklay 2** e-ink display.
The Inklay's **Image URL** app fetches `http://<computer-ip>:8080/dashboard.png` on each wake-up and shows it.
Users of this project are usually not developers. Keep explanations simple and do the technical steps for them.

## Where things are
- `dashboard.go`: **the design.** `renderDashboard(w, h, now)` is what gets drawn. Customizations normally only touch this file.
- `assets.go`: drawing helpers: `drawText`, `drawTextRight`, `drawImage` (PNG/JPEG), `drawSVG`, `fillRect`, `strokeRect`.
- `assets/`: logo, illustrations and other images. `fonts/`: font files. Both are built into the app, so **rebuild after adding or changing a file there.**
- `page.html`: the setup page in the browser (white Inklay style, Roboto Mono, responsive down to phone width). Built into the app.
- `main.go`: web server, serves `page.html` at `/` (URL + preview), `/dashboard.png?w=&h=` overrides size. Rarely needs changes.
- `build.sh`: rebuilds the double-click apps into `dist/`.
- `Start Inklay App Server.command`: Mac double-click launcher.

## Fonts and brand
- **Default font: Roboto Mono** (`fontRegular`, `fontMedium`, `fontBold`, from `fonts/`, SIL Open Font License). Use it for all text unless the user asks for another font. For another font, put its `.ttf` in `fonts/` and load it in `loadAssets` the same way.
- The Inklay look: black background, white Roboto Mono text, the INKLAY logo (`assets/inklay-logo_black.svg`, white artwork) top left, plenty of empty space. The hello-world screen in `dashboard.go` is the reference.
- `drawSVG` handles simple shapes and paths. For complex SVGs (filters, masks, embedded images) export a PNG into `assets/` and use `drawImage`.

## Display rules (follow these)
- Canvas **1872 × 1404** landscape (portrait: 1404 × 1872). Lay out relative to the size, as the existing code does with the `s` scale factor.
- **Full 8-bit greyscale (256 levels) is fine.** The Inklay dithers the image to its panel on the device, so soft greys, gradients, shading and anti-aliased text all look good. Don't force everything to black and white, and don't quantize to 16 levels.
- Pick font sizes and contrast to fit the design the user wants. Small text is fine when the user asks for a dense layout. Don't blow everything up by default.
- Avoid animation (the image is static) and very fine detail like 1 px hairlines, which dithering can break up.
- The image refreshes only every 15 min to 24 h (user setting). Seconds are pointless. Showing "updated HH:MM" helps make stale data obvious.
- Keep the output a **PNG, 8-bit greyscale** (`image.Gray`). Keep it fast (< 1 s) and small (< 500 KB); the ESP32-S3 on the device decodes it.

## Getting data
- **Render only on request.** The image is drawn fresh each time `/dashboard.png` is fetched (by the Inklay on wake-up, or when the setup page loads). Don't add background timers, refresh loops or auto-reloading previews.
- Fetch external data (weather, calendar, solar, etc.) inside the request handler with a short timeout (≈ 5 s). If it fails, still draw the dashboard with a clear "no data" note. Never return an error image.
- Cache API responses for a few minutes if the API has rate limits.
- Put API keys or settings in a `config.json` next to the app (create it if needed, and add it to `.gitignore`). Never hard-code a user's secrets in the code.
- Prefer standard library only. If a library is truly needed, `go get` it and keep `go.mod` tidy.

## Test every change
1. `go vet . && go run .` (or `INKLAY_NO_BROWSER=1 go run .` to skip opening a browser).
2. `curl -s -o /tmp/d.png localhost:8080/dashboard.png` and **look at the image** before saying it's done.
3. Check it is still 1872×1404 greyscale: `file /tmp/d.png` should say `8-bit grayscale`.
4. Stop the server, then run `./build.sh` so the double-click apps include the change.

If Go is not installed: on Mac, `brew install go` if Homebrew exists, otherwise ask the user to install it from https://go.dev/dl/.

## Example requests users might make
- "Show today's weather for Zurich" (Open-Meteo needs no API key: `https://api.open-meteo.com/v1/forecast?latitude=..&longitude=..&current=temperature_2m,weather_code&daily=temperature_2m_max,temperature_2m_min&timezone=auto`).
- "Show my next 5 calendar events" (a public/secret iCal `.ics` URL from Google Calendar is easiest).
- "Add a daily quote", "Show a countdown to our holiday", "Make it portrait", "Use German date format".
