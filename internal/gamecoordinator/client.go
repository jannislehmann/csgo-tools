package gamecoordinator

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Cludch/csgo-tools/internal/gamecoordinator/protocol"
	"github.com/Philipp15b/go-steam"
	"github.com/Philipp15b/go-steam/protocol/gamecoordinator"
	"github.com/golang/protobuf/proto" //nolint //thinks break if we use the new package
)

// HandlerMap is the map of message types to handler functions
type HandlerMap map[uint32]func(packet *gamecoordinator.GCPacket) error

// CS holds the steam client and whether the client is connected to the GameCoordinator
type CS struct {
	client      *steam.Client
	isConnected bool
	handlers    HandlerMap
}

// GCReadyEvent is used to broadcast that the GC is ready
type GCReadyEvent struct{}

// AppID describes the csgo app / steam id.
const AppID = 730

// NewCSGO creates a CS client from a steam client and registers the packet handler
func NewCSGO(client *steam.Client) *CS {
	c := &CS{client: client, isConnected: false}

	client.GC.RegisterPacketHandler(c)
	c.buildHandlerMap()

	return c
}

// SetPlaying sets the steam account to play csgo
func (c *CS) SetPlaying(playing bool) {
	if playing {
		c.client.GC.SetGamesPlayed(730)
	} else {
		c.client.GC.SetGamesPlayed()
	}
}

// ShakeHands sends a hello to the GC
func (c *CS) ShakeHands() {
	// Try to avoid not being ready on instant call of connection
	time.Sleep(5 * time.Second)

	c.Write(uint32(protocol.EGCBaseClientMsg_k_EMsgGCClientHello), &protocol.CMsgClientHello{
		Version: proto.Uint32(1),
	})
}

// HandleGCPacket takes incoming packets from the GC and coordinates them to the handler funcs.
func (c *CS) HandleGCPacket(packet *gamecoordinator.GCPacket) {
	if packet.AppId != AppID {
		log.Debug("wrong app id")
		return
	}

	handler, ok := c.handlers[packet.MsgType]
	if ok && handler != nil {
		if err := handler(packet); err != nil {
			log.Debugf("Error handling packet %d", packet.MsgType)
			log.Error(err)
			ok = false //nolint
		}
	}
}

// Write sends a message to the game coordinator.
func (c *CS) Write(messageType uint32, msg proto.Message) {
	c.client.GC.Write(gamecoordinator.NewGCMsgProtobuf(AppID, messageType, msg))
}

// emit emits an event.
func (c *CS) emit(event interface{}) {
	c.client.Emit(event)
}

// registers all csgo message handlers
func (c *CS) buildHandlerMap() {
	c.handlers = HandlerMap{
		// Welcome
		uint32(protocol.EGCBaseClientMsg_k_EMsgGCClientWelcome): c.HandleClientWelcome,

		// Match Making
		uint32(protocol.ECsgoGCMsg_k_EMsgGCCStrike15_v2_MatchList): c.HandleMatchList,
	}
}
