package steam_client

type UseCase interface {
	Connect(username, password, twoFactorSecret string)
}
