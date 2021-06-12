package gamecoordinator

import (
	"github.com/Cludch/csgo-tools/pkg/share_code"
	"github.com/Philipp15b/go-steam/v2"
	"github.com/Philipp15b/go-steam/v2/protocol/gamecoordinator"
	"github.com/golang/protobuf/proto" //nolint //thinks break if we use the new package
)

type UseCase interface {
	Connect(client *steam.Client)
	BuildHandlerMap()
	ShakeHands()
	SetPlaying(playing bool)

	GetRecentGames()
	RequestMatch(*share_code.ShareCodeData)

	Write(messageType uint32, msg proto.Message)

	HandleGCPacket(packet *gamecoordinator.GCPacket)
	HandleMatchList(packet *gamecoordinator.GCPacket)
	HandleGCReady(e *GCReadyEvent)
	HandleClientWelcome(packet *gamecoordinator.GCPacket)
}
