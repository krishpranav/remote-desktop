package rdisplay

import (
	"image"
	"time"

	"github.com/kbinani/screenshot"
)

type XVideoProvider struct{}

type XScreenGrabber struct {
	fps    int
	screen Screen
	frames chan *image.RGBA
	stop   chan struct{}
}

func (*XVideoProvider) CreateScreenGrabber(screen Screen, fps int) (ScreenGrabber, error) {
	return &XScreenGrabber{
		screen: screen,
		fps:    fps,
		frames: make(chan *image.RGBA),
		stop:   make(chan struct{}),
	}, nil
}

func (x *XVideoProvider) Screens() ([]Screen, error) {
	numScreens := screenshot.NumActiveDisplays()
	screens := make([]Screen, numScreens)
	for i := 0; i < numScreens; i++ {
		screens[i] = Screen{
			Index:  i,
			Bounds: screenshot.GetDisplayBounds(i),
		}
	}
	return screens, nil
}

func (g *XScreenGrabber) Frames() <-chan *image.RGBA {
	return g.frames
}

func (g *XScreenGrabber) Start() {
	delta := time.Duration(1000/g.fps) * time.Millisecond
	go func() {
		for {
			startedAt := time.Now()
			select {
			case <-g.stop:
				close(g.frames)
				return
			default:
				img, err := screenshot.CaptureRect(g.screen.Bounds)
				if err != nil {
					return
				}
				g.frames <- img
				ellapsed := time.Now().Sub(startedAt)
				sleepDuration := delta - ellapsed
				if sleepDuration > 0 {
					time.Sleep(sleepDuration)
				}
			}
		}
	}()
}

func (g *XScreenGrabber) Stop() {
	close(g.stop)
}

func (g *XScreenGrabber) Screen() *Screen {
	return &g.screen
}

func (g *XScreenGrabber) Fps() int {
	return g.fps
}

func NewVideoProvider() (Service, error) {
	return &XVideoProvider{}, nil
}
