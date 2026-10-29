package main

// Fonts, images and drawing helpers. Everything in fonts/ and assets/ is
// built into the app, so users only need the single file.

import (
	"bytes"
	"embed"
	"html/template"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"strings"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

//go:embed fonts/*.ttf assets page.html
var files embed.FS

// Default font family: Roboto Mono.
var (
	fontRegular *opentype.Font
	fontMedium  *opentype.Font
	fontBold    *opentype.Font
)

func loadAssets() {
	fontRegular = mustFont("fonts/RobotoMono-Regular.ttf")
	fontMedium = mustFont("fonts/RobotoMono-Medium.ttf")
	fontBold = mustFont("fonts/RobotoMono-Bold.ttf")
}

func mustRead(path string) string {
	b, err := files.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

// logoBlack returns the INKLAY logo SVG recoloured black for the white setup page.
func logoBlack() template.HTML {
	svg := strings.ReplaceAll(mustRead("assets/inklay-logo_black.svg"), `fill="white"`, `fill="#000"`)
	return template.HTML(strings.Replace(svg, "<svg ", `<svg class="logo" role="img" aria-label="INKLAY" `, 1))
}

func mustFont(path string) *opentype.Font {
	b, err := files.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	f, err := opentype.Parse(b)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func face(f *opentype.Font, size float64) font.Face {
	fc, err := opentype.NewFace(f, &opentype.FaceOptions{Size: size, DPI: 72, Hinting: font.HintingFull})
	if err != nil {
		log.Fatal(err)
	}
	return fc
}

// drawText draws text with its baseline at y, starting at x.
func drawText(img *image.Gray, f *opentype.Font, size float64, text string, x, y int, grey uint8) {
	fc := face(f, size)
	defer fc.Close()
	d := &font.Drawer{Dst: img, Src: image.NewUniform(color.Gray{grey}), Face: fc, Dot: fixed.P(x, y)}
	d.DrawString(text)
}

// drawTextRight draws text so it ends at x (right-aligned).
func drawTextRight(img *image.Gray, f *opentype.Font, size float64, text string, x, y int, grey uint8) {
	fc := face(f, size)
	defer fc.Close()
	w := font.MeasureString(fc, text).Round()
	d := &font.Drawer{Dst: img, Src: image.NewUniform(color.Gray{grey}), Face: fc, Dot: fixed.P(x-w, y)}
	d.DrawString(text)
}

// drawImage draws a PNG or JPEG from assets/ with its top-left corner at x, y.
// width 0 keeps the original size; otherwise it is scaled to that width.
func drawImage(img *image.Gray, path string, x, y, width int) {
	b, err := files.ReadFile(path)
	if err != nil {
		log.Printf("image %s: %v", path, err)
		return
	}
	src, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		log.Printf("image %s: %v", path, err)
		return
	}
	sb := src.Bounds()
	if width <= 0 || width == sb.Dx() {
		draw.Draw(img, image.Rect(x, y, x+sb.Dx(), y+sb.Dy()), src, sb.Min, draw.Over)
		return
	}
	h := sb.Dy() * width / sb.Dx()
	scaleInto(img, src, image.Rect(x, y, x+width, y+h))
}

// drawSVG draws an SVG from assets/ scaled to the given width, top-left at x, y.
// Simple shapes and paths work; filters, masks and embedded images do not.
func drawSVG(img *image.Gray, path string, x, y, width int) {
	b, err := files.ReadFile(path)
	if err != nil {
		log.Printf("svg %s: %v", path, err)
		return
	}
	icon, err := oksvg.ReadIconStream(bytes.NewReader(b), oksvg.IgnoreErrorMode)
	if err != nil {
		log.Printf("svg %s: %v", path, err)
		return
	}
	h := int(float64(width) * icon.ViewBox.H / icon.ViewBox.W)
	icon.SetTarget(0, 0, float64(width), float64(h))
	rgba := image.NewRGBA(image.Rect(0, 0, width, h))
	icon.Draw(rasterx.NewDasher(width, h, rasterx.NewScannerGV(width, h, rgba, rgba.Bounds())), 1)
	draw.Draw(img, image.Rect(x, y, x+width, y+h), rgba, image.Point{}, draw.Over)
}

// scaleInto does a simple box-filter resize of src into dst rectangle r.
func scaleInto(dst *image.Gray, src image.Image, r image.Rectangle) {
	sb := src.Bounds()
	for dy := r.Min.Y; dy < r.Max.Y; dy++ {
		sy0 := sb.Min.Y + (dy-r.Min.Y)*sb.Dy()/r.Dy()
		sy1 := max(sy0+1, sb.Min.Y+(dy-r.Min.Y+1)*sb.Dy()/r.Dy())
		for dx := r.Min.X; dx < r.Max.X; dx++ {
			sx0 := sb.Min.X + (dx-r.Min.X)*sb.Dx()/r.Dx()
			sx1 := max(sx0+1, sb.Min.X+(dx-r.Min.X+1)*sb.Dx()/r.Dx())
			var sum, a, n uint32
			for sy := sy0; sy < sy1; sy++ {
				for sx := sx0; sx < sx1; sx++ {
					c := color.GrayModel.Convert(src.At(sx, sy)).(color.Gray)
					_, _, _, al := src.At(sx, sy).RGBA()
					sum += uint32(c.Y) * (al >> 8)
					a += al >> 8
					n++
				}
			}
			if a == 0 {
				continue
			}
			fg := sum / a
			alpha := a / n
			bg := uint32(dst.GrayAt(dx, dy).Y)
			dst.SetGray(dx, dy, color.Gray{uint8((fg*alpha + bg*(255-alpha)) / 255)})
		}
	}
}

func fillRect(img *image.Gray, x0, y0, x1, y1 int, grey uint8) {
	draw.Draw(img, image.Rect(x0, y0, x1, y1), image.NewUniform(color.Gray{grey}), image.Point{}, draw.Src)
}

func strokeRect(img *image.Gray, x0, y0, x1, y1 int, grey uint8, t int) {
	fillRect(img, x0, y0, x1, y0+t, grey)
	fillRect(img, x0, y1-t, x1, y1, grey)
	fillRect(img, x0, y0, x0+t, y1, grey)
	fillRect(img, x1-t, y0, x1, y1, grey)
}
