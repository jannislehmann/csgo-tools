package valveapi

import (
	"compress/bzip2"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// MatchResponse contains information about the latest match
type MatchResponse struct {
	Result struct {
		Nextcode string `json:"nextcode"`
	} `json:"result"`
}

// InvalidDownloadURLError is return when the url to be download is invalid or malicious.
type InvalidDownloadURLError struct{}

func (e *InvalidDownloadURLError) Error() string {
	return "Invalid download url"
}

// DemoNotFoundError is used when a valid matchid / demo is not found or can no longer be downloaded.
type DemoNotFoundError struct {
	URL string
}

func (e *DemoNotFoundError) Error() string {
	return fmt.Sprintf("Demo no longer downloadable: %s", e.URL)
}

// InvalidMatchHistoryCredentials is used to notify when the supplied credentials are not valid / cannot be used with the api.
type InvalidMatchHistoryCredentials struct {
	SteamID string
}

func (e *InvalidMatchHistoryCredentials) Error() string {
	return fmt.Sprintf("Invalid match history credentials for steam id %v", e.SteamID)
}

// GetNextMatch returns the next match's share code.
// It uses the saved share codes as the current one.
func GetNextMatch(steamAPIKey string, steamID uint64, historyAuthenticationCode string, lastShareCode string) (string, error) {
	// Get latest match
	u, err := url.Parse("https://api.steampowered.com/ICSGOPlayers_730/GetNextMatchSharingCode/v1")
	if err != nil {
		log.Error(err)
	}

	steamIDString := strconv.FormatUint(steamID, 10)

	// Build query
	q := u.Query()
	q.Set("key", steamAPIKey)
	q.Set("steamid", steamIDString)
	q.Set("steamidkey", historyAuthenticationCode)
	q.Set("knowncode", lastShareCode)
	u.RawQuery = q.Encode()

	matchResponse := &MatchResponse{}

	// Request match code
	r, err := http.Get(u.String())
	if err != nil {
		log.Error(err)
		return "", err
	}

	// Forbidden = wrong api keys
	// Precondition Failed = Know match code or steam id wrong
	if r.StatusCode == http.StatusForbidden || r.StatusCode == http.StatusPreconditionFailed {
		r.Body.Close()
		return "", &InvalidMatchHistoryCredentials{SteamID: steamIDString}
	}

	// Accepted means that there is no recent match code available
	if r.StatusCode == http.StatusAccepted {
		r.Body.Close()
		return "", nil
	}

	err = json.NewDecoder(r.Body).Decode(matchResponse)

	if err != nil {
		r.Body.Close()
		log.Error(err)
		return "", err
	}

	defer r.Body.Close()

	return matchResponse.Result.Nextcode, nil
}

// DownloadDemo will download a demo from an url and decompress and store it in local filepath.
// It writes to the destination file as it downloads it, without
// loading the entire file into memory.
func DownloadDemo(url string, demoDir string, lastModified time.Time) error {
	// Validate the url
	re := regexp.MustCompile(`^http:\/\/replay[\d]{3}\.valve\.net\/730\/[\d]{21}_([\d]*)\.dem\.bz2$`)

	if !re.MatchString(url) {
		return &InvalidDownloadURLError{}
	}

	// Get file name
	fileName := strings.Split(path.Base(url), ".")[0] + ".dem"
	filePath := path.Join(demoDir, fileName)

	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		log.Error(err)
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url) //nolint // We have to take dynamic replay urls in order to download them. URL is validated before.
	if err != nil || resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return &DemoNotFoundError{URL: url}
	}

	// Decompress and write to file
	cr := bzip2.NewReader(resp.Body)
	_, err = io.Copy(out, cr)

	defer resp.Body.Close()

	if err != nil {
		return err
	}

	// Update file modified information
	err = os.Chtimes(filePath, lastModified, lastModified)
	if err != nil {
		log.Warnf("unable to set correct last modified date for demo %v", fileName)
		log.Error(err)
	}

	log.Infof("downloaded demo %v", fileName)

	return nil
}
