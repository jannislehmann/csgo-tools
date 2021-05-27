# CSGO Demo Tools

The CSGO tool suit can automatically detect and download new CSGO official matchmaking demos using the GameCoordinator.
In order to do this, a few API credentials and a **separate** Steam account is needed. The application creates a CSGO game sessions using the separate account
and uses the Steam web API to check whether a new demo can be fetched. If that is the case, the application sends a full match info request to the game's GameCoordinator.
The GC then returns information about the match which also contain a download link.

The tool saves all the match ids from the demos in the `demos` directory. This is used to prevent downloading a demo every few minutes.

## Tools

The toolset currently features the following tools.

### valveapiclient

The API client consumes Valve's game history API and saves the game share codes in the database.

### Gameclient

The game client uses the CSGO gamecoordinator to talk to the ingame "API". By doing so, the tool can request match information, history and most importantly the download links for each demo.

### Demodownloader

The demo downloader takes the demo urls from the database and downloads them if they are missing.

## Demoparser

The demo parser parses the previously downloaded demo files and calculates the following statistics for each player:
* Kills, Deaths, (Flash) Assists, Headshots and percentage
  * Information about the kill such as wallbang, flashed, through smoke
  * Also on per weapon basis
* Player MVPs
* Map
* Team Scores

When a new demoparser version gets released, the tool will automatically reparse demos, which were parsed with
an older version. Therefore, new statistiscs will also be added for older demos. This, however, requires the demos to be permanently persisted.

## Usage

Get the latest binary and set up your demo location and the config file.

### `config.json`

Copy the `config.json.example` in the `configs` dir and rename it to `config.json` in the same dir.

The `demosDir` setting is the directory, in which the demos should be stored (e.g. `demos/`).
The `debug` parameter can be enabled to receive a few more debug output.

### Steam

| Key   |      Value      |  Explanation |
|----------|-------------:|------:|
| `apiKey` |   `12345`   | The Steam Web API key. Can be generate [here](https://steamcommunity.com/dev/apikey) |
| `username` |   `user`   |  Steam username |
| `password` |   `totally_secret`   |  Steam password |
| `twoFactorSecret` |   `aGV5IQ==`   | Base64 encoded two factor secret. Can be generated using e.g. the [Steam Desktop Authenticator](https://github.com/Jessecar96/SteamDesktopAuthenticator) |

### CSGO

The csgo array contains mulitple account information about accounts to watch.

| Key   |      Value      |  Explanation |
|----------|-------------:|------:|
| `matchHistoryAuthenticationCode` |  `1234-ABCDE-5678`  | The match history authentication code can be generated [here](https://help.steampowered.com/en/wizard/HelpWithGameIssue/?appid=730&issueid=128) |
| `knownMatchCode` | `CSGO-abcde-efghi-jklmn-opqrs-tuvwx` |  A share code from one of your latest matches. Can be received via the Game -> Matches |
| `steamId` |   `76561198185324675`   |  The SteamID64 of the account to watch |

### Database

| Key   |      Value      |  Explanation |
|----------|-------------:|------:|
| `host` |   `localhost`   |  The database host |
| `port` | `5432` | The database port |
| `username` |   `csgo`   |  Username of the database user |
| `password` |   `b`   |  Secret password of the database user |
| `database` |   `csgo`   | The database name to store the data in |

## Disclaimer

This is my first ever Golang project, thus you might find some bad practice and a few performance issues in the long run.
The project structure is horrible but it works for now. If you have suggestions please go a head and create an issue. This would help me a lot!

This tool is not affiliated with Valve Software or Steam.

## Other projects that helped me a lot

* [go-steam](https://github.com/Philipp15b/go-steam)
* [cs-go](https://github.com/Gacnt/cs-go)
* [go-dota2](https://github.com/paralin/go-dota2)
* [csgo](https://github.com/ValvePython/csgo)
