package rdisplay

import "image"

type ScreenGrabber interface {
	Start()
	Frames() <-chan *image.RGBA
	Stop()
	Fps() int
	Screen() *Screen
}

type Screen struct {
	Index  int
	Bounds image.Rectangle
}

type Service interface {
	CreateScreenGrabber(screen Screen, fps int) (ScreenGrabber, error)
	Screens() ([]Screen, error)
}
