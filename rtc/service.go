package rtc

import (
	"io"
)

type videoStreamer interface {
	start()
	close()
}

type RemoteScreenConnection interface {
	io.Closer
	ProcessOffer(offer string) (string, error)
}

type Service interface {
	CreateRemoteScreenConnection(screenIx int, fps int) (RemoteScreenConnection, error)
}
