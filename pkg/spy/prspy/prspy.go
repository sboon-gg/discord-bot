package prspy

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	PRSPY_URL = "https://servers.realitymod.com/api/ServerInfo"
)

type player struct {
	TagName string `json:"name"`
	IsAI    int    `json:"isAI"`
}

func (p player) IGN() string {
	return strings.Split(p.TagName, " ")[1]
}

func (p player) Tag() string {
	return strings.Split(p.TagName, " ")[0]
}

type server struct {
	Players []player `json:"players"`
}

type prspyData struct {
	Servers []server `json:"servers"`
}

func FetchData() (*prspyData, error) {
	resp, err := http.Get(PRSPY_URL)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Invalid status code: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data prspyData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func FetchAllPlayers() (map[string]player, error) {
	players := make(map[string]player)

	data, err := FetchData()
	if err != nil {
		return players, err
	}

	for _, sv := range data.Servers {
		for _, p := range sv.Players {
			if p.IsAI == 0 {
				players[p.IGN()] = p
			}
		}
	}

	return players, nil
}

// goroutine: fetch PRSpy data, give roles based on activity
// remove all ingame roles on inactivity
