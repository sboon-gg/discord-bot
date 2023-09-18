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

type Player struct {
	TagName string `json:"name"`
	IsAI    int    `json:"isAI"`
}

func (p Player) IGN() string {
	return strings.Split(p.TagName, " ")[1]
}

func (p Player) Tag() string {
	return strings.Split(p.TagName, " ")[0]
}

type Server struct {
	Players []Player `json:"players"`
}

type PRSpyData struct {
	Servers []Server `json:"servers"`
}

func FetchData() (*PRSpyData, error) {
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

	var data PRSpyData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func GetAllPlayers(data *PRSpyData) map[string]Player {
	players := make(map[string]Player)

	for _, sv := range data.Servers {
		for _, p := range sv.Players {
			if p.IsAI == 0 {
				players[p.IGN()] = p
			}
		}
	}

	return players
}
