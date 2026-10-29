package main

// This file is the dashboard design. Change renderDashboard to change what
// the Inklay shows. Coordinates are in pixels on the 1872 x 1404 canvas.

import (
	"image"
	"time"
)

func renderDashboard(w, h int, now time.Time) *image.Gray {
	img := image.NewGray(image.Rect(0, 0, w, h))
	fillRect(img, 0, 0, w, h, 0) // black background

	const margin = 50

	// Top bar: logo left, time right.
	drawSVG(img, "assets/inklay-logo_black.svg", margin, margin, 200)
	drawTextRight(img, fontRegular, 24, now.Format("15:04"), w-margin, 68, 255)

	// Centre: waving illustration + greeting.
	const illW, illH = 602, 415
	ix, iy := w/2-illW/2+1, h/2-illH/2-46
	drawImage(img, "assets/inklay-welcome-w.png", ix, iy, 0)
	drawText(img, fontRegular, 24, "Hello World", ix+556, iy+116, 255)

	// Bottom left: how to get started.
	steps := []string{
		"1. Open the folder in Claude Code",
		"2. Prompt your ideas",
		"3. Run Start Inklay App Server.command",
		"4. Add the URL to the INKLAY Image URL app.",
	}
	y := h - margin - 32*(len(steps)-1)
	for _, s := range steps {
		drawText(img, fontRegular, 24, s, margin, y, 255)
		y += 32
	}
	return img
}
