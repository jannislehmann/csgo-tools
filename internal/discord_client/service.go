package discord_client

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	dg *discordgo.Session
}

func NewService(apiKey string) *Service {
	service := &Service{}
	service.connect(apiKey)
	return service
}

func (s *Service) connect(apiKey string) {

	if dg, err := discordgo.New("Bot " + apiKey); err != nil {
		log.Warn("discord: error creating Discord session", err)
		return
	} else {
		s.dg = dg
	}

	log.Println("discord: bot is now running.")
}

func (s *Service) SendMessage(message, channelID string) {
	if _, err := s.dg.ChannelMessageSend(channelID, message); err != nil {
		log.Warn("discord: an error occurred when sending: ", err)
	} else {
		log.Debugf("discord: send message: %s", message)
	}
}

func (s *Service) Close() {
	s.dg.Close()
}
