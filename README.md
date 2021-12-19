# CSGO Demo Tools

The CSGO tool suit can automatically detect and download new CSGO official matchmaking demos using the GameCoordinator.
In order to do this, a few API credentials and a **separate** Steam account is needed. The application creates a CSGO game sessions using the separate account
and uses the Steam Web API to check whether a new demo can be fetched. If that's the case, the application sends a full match info request to the game's GameCoordinator.
The GC then returns information about the match which also contain a download link.

## Tools

The toolset currently features the following tools. All tools share one MongoDB database instance.

### Auth

The auth service enables a user to sign in using his / her own Steam account. The generated token can be used with other services at a later point.

### ValveAPI client

The API client consumes Valve's game history API and saves the game share codes in the database.
In order to add a new steam / csgo user, whose demos should be monitored, a user must be manually created in the database.

### Game client

The game client uses the CSGO gamecoordinator to talk to the ingame "API". By doing so, the tool can request match information, history and most importantly the download links for each demo.

### Demo Downloader

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

## REST API

The REST api serves basic match and player stats via the following routes.

| Route | Description |
|---------------------|-------------:|
| `/match`            | Lists all available matches. |
| `/match/:id`        | Serves information and outcome about one specific match. |
| `/player/:id`       | Lists information about one player. |
| `/player/:id/stats` | Calculates and serves average stats for one player. |

## Usage

Get the latest binary and set up your demo location and the config file.

### `config.json`

Copy the `config.json.example` in the `configs` dir and rename it to `config.json` in the same dir.

The `demosDir` setting is the directory, in which the demos should be stored (e.g. `demos/`).
The `debug` parameter can be enabled to receive a few more debug output.

You can also use ENV vars to override single or set all configuration variables. The formatting for the configuration is as with the JSON configuration. The ENV base is `CSGO`. The Steam two factor secret turns into `STEAM_TWOFACTORSECRET`.

### Auth

| Key   |      Value      |  Explanation |
|----------|-------------:|------:|
| `host` |   `http://localhost:8080`   |  The host url for the authentication callback. |
| `secret` |   `http://localhost:8080`   |  The authentication token secret. E.g. `openssl rand -base64 32` |

### Steam

| Key   |      Value      |  Explanation |
|----------|-------------:|------:|
| `apiKey` |   `12345`   | The Steam Web API key. Can be generate [here](https://steamcommunity.com/dev/apikey) |
| `username` |   `user`   |  Steam username |
| `password` |   `totally_secret`   |  Steam password |
| `twoFactorSecret` |   `aGV5IQ==`   | Base64 encoded two factor secret. Can be generated using e.g. the [Steam Desktop Authenticator](https://github.com/Jessecar96/SteamDesktopAuthenticator) |

### Database

| Key   |      Value      |  Explanation |
|----------|-------------:|------:|
| `host` |   `localhost`   |  The database host |
| `port` | `27017` | The database port |
| `username` |   `csgo`   |  Username of the database user |
| `password` |   `b`   |  Secret password of the database user |
| `database` |   `csgo`   | The database name to store the data in |

### Parser

| Key   |      Value      |  Explanation |
|----------|-------------:|------:|
| `workerCount` |   `5`   |  The amount of workers to parellely parse demos |

## Disclaimer

This is my first ever Golang project, thus you might find some bad practice and a few performance issues in the long run.
The project structure is horrible but it works for now. If you have suggestions please go a head and create an issue. This would help me a lot!

This tool is not affiliated with Valve Software or Steam.

## Other projects that helped me a lot

* [go-steam](https://github.com/Philipp15b/go-steam)
* [cs-go](https://github.com/Gacnt/cs-go)
* [go-dota2](https://github.com/paralin/go-dota2)
* [csgo](https://github.com/ValvePython/csgo)
* [Uber style guide](https://github.com/uber-go/guide/blob/master/style.md)
