package gamecoordinator

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Philipp15b/go-steam/v2"
	csgo "github.com/Philipp15b/go-steam/v2/csgo/protocol/protobuf"
	"github.com/Philipp15b/go-steam/v2/protocol/gamecoordinator"
	"github.com/golang/protobuf/proto" //nolint //thinks break if we use the new package
)

// HandlerMap is the map of message types to handler functions
type HandlerMap map[uint32]func(packet *gamecoordinator.GCPacket)

// GC holds the steam client and whether the client is connected to the GameCoordinator
type GC struct {
	client      *steam.Client
	isConnected bool
	handlers    HandlerMap
}

// GCReadyEvent is used to broadcast that the GC is ready
type GCReadyEvent struct{}

// AppID describes the csgo app / steam id.
const AppID = 730

// NewCSGO creates a CS client from a steam client and registers the packet handler
func (s *Service) Connect(client *steam.Client) {
	s.gc = &GC{client: client, isConnected: false}
	s.BuildHandlerMap()
	s.gc.client.GC.RegisterPacketHandler(s)
	s.SetPlaying(true)
	s.ShakeHands()
}

// SetPlaying sets the steam account to play csgo
func (s *Service) SetPlaying(playing bool) {
	if playing {
		s.gc.client.GC.SetGamesPlayed(730)
	} else {
		s.gc.client.GC.SetGamesPlayed()
	}
}

// ShakeHands sends a hello to the GC
func (s *Service) ShakeHands() {
	// Try to avoid not being ready on instant call of connection
	time.Sleep(5 * time.Second)

	s.Write(uint32(csgo.EGCBaseClientMsg_k_EMsgGCClientHello), &csgo.CMsgClientHello{
		Version: proto.Uint32(1),
	})
}

// HandleGCPacket takes incoming packets from the GC and coordinates them to the handler funcs.
func (s *Service) HandleGCPacket(packet *gamecoordinator.GCPacket) {
	if packet.AppId != AppID {
		log.Debug("wrong app id")
		return
	}

	if handler, ok := s.gc.handlers[packet.MsgType]; ok {
		handler(packet)
	}
}

// Write sends a message to the game coordinator.
func (s *Service) Write(messageType uint32, msg proto.Message) {
	s.gc.client.GC.Write(gamecoordinator.NewGCMsgProtobuf(AppID, messageType, msg))
}

// emit emits an event.
func (s *Service) emit(event interface{}) {
	s.gc.client.Emit(event)
}

// registers all csgo message handlers
func (s *Service) BuildHandlerMap() {
	s.gc.handlers = HandlerMap{
		// Welcome
		uint32(csgo.EGCBaseClientMsg_k_EMsgGCClientWelcome): s.HandleClientWelcome,

		// Match Making
		uint32(csgo.ECsgoGCMsg_k_EMsgGCCStrike15_v2_MatchList): s.HandleMatchList,
	}
}
