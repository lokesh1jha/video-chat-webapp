package webrtc

import (
	"sync"
	"videochat/pkg/webrtc"
)

type Room struct {
	Peers *Peers
	Hub   *chat.Hub
}

type Peers struct {
	ListLock    sync.RWMutex
	Connection  []PeerConnectionState
	TrackLocals map[string]*webrtc.TrackLocalStaticRTP
}

func (p *Peers) DispatchKeyFrame() {

}
