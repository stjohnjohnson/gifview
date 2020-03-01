package gifview

import (
	"fmt"
	"image/gif"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Omnikron13/pixelview"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// GifView is a box which displays animated gifs via Omnikron13's pixelview
// dynamic color rendering.  It automatically draws the right frame based on
// time elapsed since creation.  You can trigger re-drawing by executing
// Animate(tview.Application) in a goroutine.
type GifView struct {
	sync.Mutex
	*tview.Box

	// Timing for the exection of frames
	delay         []time.Duration
	frames        []string
	startTime     time.Time
	totalDuration time.Duration
}

// NewGifView returns a new GifView.
func NewGifView() *GifView {
	return &GifView{
		Box:       tview.NewBox(),
		startTime: time.Now(),
	}
}

// FromImage creates a new GifView from a GIF Image
func FromImage(image *gif.GIF) (*GifView, error) {
	g := NewGifView()
	return g.SetImage(image)
}

// FromImagePath creates a new GifView from a file on disk
func FromImagePath(imagePath string) (*GifView, error) {
	g := NewGifView()
	return g.SetImagePath(imagePath)
}

// SetImagePath sets the image to a given GIF path
func (g *GifView) SetImagePath(imagePath string) (*GifView, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return g, fmt.Errorf("Unable to open file: %v", err)
	}
	defer file.Close()

	image, err := gif.DecodeAll(file)
	if err != nil {
		return g, fmt.Errorf("Unable to decode GIF: %v", err)
	}

	return g.SetImage(image)
}

// SetImage sets the content to a given gif.GIF
func (g *GifView) SetImage(image *gif.GIF) (*GifView, error) {
	g.Lock()
	defer g.Unlock()

	// Store delay in milliseconds
	g.totalDuration = time.Duration(0)
	for _, i := range image.Delay {
		d := time.Duration(i*10) * time.Millisecond
		g.delay = append(g.delay, d)
		g.totalDuration += d
	}

	// Set height,width of the box
	g.SetRect(0, 0, image.Config.Width, image.Config.Height)

	// Convert images to text
	frames := []string{}
	for i, img := range image.Image {
		parsed, err := pixelview.FromImage(img)
		if err != nil {
			return g, fmt.Errorf("Unable to convert frame %d: %v", i, err)
		}
		frames = append(frames, parsed)
	}

	// Store the output
	g.frames = frames

	return g, nil
}

// GetCurrentFrame returns the current frame the GIF is on
func (g *GifView) GetCurrentFrame() int {
	dur := time.Since(g.startTime) % g.totalDuration
	for i, d := range g.delay {
		dur -= d
		if dur < 0 {
			return i
		}
	}
	return 0
}

// Draw renders the current frame of the GIF
func (g *GifView) Draw(screen tcell.Screen) {
	g.Lock()
	defer g.Unlock()

	currentFrame := g.GetCurrentFrame()

	frame := strings.Split(g.frames[currentFrame], "\n")
	x, y, w, _ := g.GetInnerRect()

	for i, line := range frame {
		tview.Print(screen, line, x, y+i, w, tview.AlignLeft, tcell.ColorWhite)
	}
}

var globalAnimationMutex = &sync.Mutex{}

// Animate triggers the application to redraw every 50ms
func Animate(app *tview.Application) {
	globalAnimationMutex.Lock()
	defer globalAnimationMutex.Unlock()

	for {
		app.QueueUpdateDraw(func() {})
		time.Sleep(50 * time.Millisecond)
	}
}
