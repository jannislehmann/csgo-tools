package discord_client

type UseCase interface {
	SendMessage(message, ChannelID string)
	Close()
}
