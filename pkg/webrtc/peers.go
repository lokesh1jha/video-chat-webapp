package webrtc

import (
	"log"
	"sync"
	"videochat/pkg/chat"
	"videochat/pkg/webrtc/v2"

	"github.com/gofiber/websocket"
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

type PeerConnectionState struct {
	PeerConnection *webrtc.PeerConnection
	websocket      *ThreadSafeWriter
}

type ThreadSafeWriter struct {
	Conn  *websocket.Conn
	Mutex sync.Mutex
}

func (t *ThreadSafeWriter) WriteJSON(v interface{}) error {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	return t.Conn.WriteJSON(v)
}

func (p *Peers) AddTrack(t *webrts.TrackRemote) *webrtc.TrackLocalStaticRTP {
	p.ListLock.Lock()
	defer func() {
		p.ListLock.Unlock()
		p.SignalPeerConnection()
	}()

	trackLocal, err := webrtc.NewTrackLocalStaticRTP(t.Codec().RTPCodecCapability, "video", t.ID(), t.StreamID())

	if err != nil {
		log.Println(err.Error())
		return nil
	}

	p.TrackLocals[t.ID()] = trackLocal
	return trackLocal
}

func (p *Peers) RemoveTrack(t *webrtc.TrackLocalStaticRTP) {
	p.ListLock.Lock()
	defer func() {
		p.ListLock.Unlock()
		p.SignalPeerConnection()
	}()

	delete(p.TrackLocals, t.ID())
}

func (p *Peers) SignalPeerConnection() {
	p.ListLock.Lock()
	defer func() {
		p.ListLock.Unlock()
		p.DispatchKeyFrame()
	}()
	attemptSync := func() (tryAgain bool) {
		for i := range p.Connections {
			if p.Connections[i].PeerConnection.ConnectionState() == webrtc.PeerConnectionStateClosed {
				p.Connections = append(p.Connections[:i], p.Connections[i+1:]...)
				log.Printf("a", p.Connections)
				return false
			}

			existingSender := map[string]bool{}
			for _, sender := range p.Connections[i].PeerConnection.GetSenders() {
				if sender.Track() == nil {
					continue
				}

				existingSender[sender.Track().ID()] = true

				if _, ok := p.TrackLocals[sender.Track().ID()]; !ok {
					if err := p.Connection[i].PeerConnection.RemoveTrack(sender); err != nil {
						return true
					}
				}
			}

			for _, receiver := range p.Connection[i].PeerConnection.GetReceivers() {
				if receiver.Track() == nil {
					continue
				}
				existingSender[receiver.Track().ID()] = true
			}

			for trackID := range p.TrackLocals {
				if _, ok := existingSender[trackID]; !ok {
					if _, err := p.Connection[i].PeerConnection.AddTrack(p.TrackLocals[trackID]); err != nil {
						return true
					}
				}
			}
			return false
		}

		return true
	}
}

func (p *Peers) DispatchKeyFrame() {

}

type websocketMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}
