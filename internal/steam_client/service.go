package steam_client

import (
	"github.com/Cludch/csgo-tools/internal/gamecoordinator"
	"github.com/Philipp15b/go-steam/v2"
	"github.com/Philipp15b/go-steam/v2/protocol/steamlang"
	"github.com/Philipp15b/go-steam/v2/totp"

	log "github.com/sirupsen/logrus"
)

type Service struct {
	gamecoordinatorService gamecoordinator.UseCase
}

func NewService(g gamecoordinator.UseCase) *Service {
	return &Service{
		gamecoordinatorService: g,
	}
}

func (s *Service) Connect(username, password, twoFactorSecret string) {
	totpInstance := totp.NewTotp(twoFactorSecret)

	myLoginInfo := new(steam.LogOnDetails)
	myLoginInfo.Username = username
	myLoginInfo.Password = password
	twoFactorCode, err := totpInstance.GenerateCode()

	if err != nil {
		log.Error(err)
	}

	myLoginInfo.TwoFactorCode = twoFactorCode

	client := steam.NewClient()
	_, connectErr := client.Connect()
	if connectErr != nil {
		log.Panic(connectErr)
	}

	for event := range client.Events() {
		switch e := event.(type) {
		case *steam.ConnectedEvent:
			log.Info("connected to steam. Logging in...")
			client.Auth.LogOn(myLoginInfo)
		case *steam.LoggedOnEvent:
			log.Info("logged on")
			client.Social.SetPersonaState(steamlang.EPersonaState_Invisible)

			s.gamecoordinatorService.Connect(client)
		case *gamecoordinator.GCReadyEvent:
			s.gamecoordinatorService.HandleGCReady(e)
		case steam.DisconnectedEvent:
			log.Panic("steam_client: disconnected")
		case steam.FatalErrorEvent:
			log.Fatal(e)
		}
	}
}
