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
	dg, err := discordgo.New("Bot " + apiKey)
	if err != nil {
		log.Warn("discord: error creating Discord session", err)
		return
	}

	s.dg = dg

	log.Println("discord: bot is now running.")
}

func (s *Service) SendMessage(message, channelId string) {
	if _, err := s.dg.ChannelMessageSend(channelId, message); err != nil {
		log.Warn("discord: an error occurred when sending: ", err)
	}
}

func (s *Service) Close() {
	s.dg.Close()
}
