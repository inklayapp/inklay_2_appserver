// Inklay App Server: a tiny local web server that renders a greyscale PNG
// for the Inklay "Image URL" app. Double-click to start; a browser page
// opens with the URL to paste into the Inklay app.
package main

import (
	"bytes"
	"fmt"
	"html/template"
	"image/png"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

// Inklay 2: 10.3" panel, 1872 x 1404. Serve full 8-bit greyscale (256 levels);
// the device dithers it down to its 16 panel levels.
const (
	defaultW = 1872
	defaultH = 1404
)

func main() {
	var err error
	loadAssets()

	mux := http.NewServeMux()
	mux.HandleFunc("/dashboard.png", handlePNG)
	mux.HandleFunc("/", handleIndex)
	mux.Handle("/fonts/", http.FileServerFS(files)) // Roboto Mono for the setup page

	// Use port 8080, or the next free one if something else already has it.
	var ln net.Listener
	port := 8080
	for ; port < 8100; port++ {
		if ln, err = net.Listen("tcp", fmt.Sprintf(":%d", port)); err == nil {
			break
		}
	}
	if ln == nil {
		log.Fatalf("no free port between 8080 and 8099: %v", err)
	}

	ip := lanIP()
	imageURL := fmt.Sprintf("http://%s:%d/dashboard.png", ip, port)
	pageURL := fmt.Sprintf("http://localhost:%d/", port)

	fmt.Println("==============================================")
	fmt.Println(" Inklay App Server is running")
	fmt.Println()
	fmt.Println(" Paste this URL into the Inklay Image URL app:")
	fmt.Println("   " + imageURL)
	fmt.Println()
	fmt.Println(" Keep this window open. Close it to stop.")
	fmt.Println("==============================================")

	if os.Getenv("INKLAY_NO_BROWSER") == "" {
		go func() {
			time.Sleep(300 * time.Millisecond)
			openBrowser(pageURL)
		}()
	}
	log.Fatal(http.Serve(ln, mux))
}

// ---------- image ----------

func handlePNG(w http.ResponseWriter, r *http.Request) {
	width := intParam(r, "w", defaultW)
	height := intParam(r, "h", defaultH)

	img := renderDashboard(width, height, time.Now())

	var buf bytes.Buffer
	enc := png.Encoder{CompressionLevel: png.BestCompression}
	if err := enc.Encode(&buf, img); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	w.Write(buf.Bytes())
	log.Printf("served %dx%d PNG (%d KB) to %s", width, height, buf.Len()/1024, r.RemoteAddr)
}

func intParam(r *http.Request, key string, def int) int {
	if v, err := strconv.Atoi(r.URL.Query().Get(key)); err == nil && v >= 100 && v <= 4000 {
		return v
	}
	return def
}

// ---------- setup page ----------

var page = template.Must(template.New("index").Parse(mustRead("page.html")))

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	_, port, _ := net.SplitHostPort(r.Host)
	if port == "" {
		port = "80"
	}
	page.Execute(w, map[string]any{
		"ImageURL": fmt.Sprintf("http://%s:%s/dashboard.png", lanIP(), port),
		"Logo":     logoBlack(),
	})
}

// ---------- helpers ----------

// lanIP finds the address other devices on the Wi-Fi can reach us at.
// Dialing UDP sends no packets; it just asks the OS which interface it would use.
func lanIP() string {
	if c, err := net.Dial("udp", "192.168.0.1:80"); err == nil {
		defer c.Close()
		if a, ok := c.LocalAddr().(*net.UDPAddr); ok && !a.IP.IsLoopback() {
			return a.IP.String()
		}
	}
	addrs, _ := net.InterfaceAddrs()
	for _, a := range addrs {
		if n, ok := a.(*net.IPNet); ok && n.IP.To4() != nil && n.IP.IsPrivate() {
			return n.IP.String()
		}
	}
	return "localhost"
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	_ = cmd.Start()
}
