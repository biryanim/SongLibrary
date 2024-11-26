package http

import (
	"encoding/json"
	"fmt"
	"github.com/biryanim/SongLibrary/internal/entities"
	"net/http"
)

func externalAPIRequest(group, songName string) (*entities.SongDetails, int, error) {
	url := fmt.Sprintf("http://external-api/info?group=%s&song=%s", group, songName)

	resp, err := http.Get(url)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to make request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			return nil, resp.StatusCode, fmt.Errorf("incorrect request: %v", resp.StatusCode)
		} else {
			return nil, resp.StatusCode, fmt.Errorf("unexpected response code: %v", resp.StatusCode)
		}
	}

	var song entities.SongDetails
	if err = json.NewDecoder(resp.Body).Decode(&song); err != nil {
		return nil, 0, err
	}
	return &song, resp.StatusCode, nil
}
