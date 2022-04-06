package faceitapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// MatchResponse contains information about the latest match.
type MatchDetailResponse struct {
	DemoUrl []string `json:"demo_url"`
}

// GetPlayerMatchHistory returns the match history for a given player.
func GetMatchDetails(faceitAPIKey string, matchId string) (*MatchDetailResponse, error) {
	playerResponse := &MatchDetailResponse{}

	// Request match details.
	client := &http.Client{}
	request, _ := http.NewRequest("GET", fmt.Sprintf("https://open.faceit.com/data/v4/matches/%s", matchId), nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", faceitAPIKey))
	r, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if r.StatusCode == http.StatusTooManyRequests || r.StatusCode == http.StatusServiceUnavailable {
		r.Body.Close()
		return nil, &FaceitApiConnectionIssues{}
	}

	if r.StatusCode == http.StatusUnauthorized || r.StatusCode == http.StatusForbidden {
		r.Body.Close()
		return nil, &InvalidFaceitApiCredentials{}
	}

	if err = json.NewDecoder(r.Body).Decode(playerResponse); err != nil {
		r.Body.Close()
		return nil, err
	}

	defer r.Body.Close()

	return playerResponse, nil
}
