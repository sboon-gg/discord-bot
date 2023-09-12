package prspy

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

const (
	PRSPY_URL = "https://servers.realitymod.com/api/ServerInfo"
)

type player struct {
	TagName string `json:"name"`
	IsAI    bool   `json:"isAI"`
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
		return nil, errors.New("Invalid status code")
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

func FetchAllPlayers() map[string]player {
	players := make(map[string]player)

	data, err := FetchData()
	if err != nil {
		return players
	}

	for _, sv := range data.Servers {
		for _, p := range sv.Players {
			players[p.IGN()] = p
		}
	}

	return players
}

// goroutine: fetch PRSpy data, give roles based on activity
// remove all ingame roles on inactivity
