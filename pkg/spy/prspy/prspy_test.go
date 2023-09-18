package prspy_test

import (
	"testing"

	"github.com/sboon-gg/sboon-bot/pkg/spy/prspy"
	"github.com/stretchr/testify/assert"
)

func TestFetchData(t *testing.T) {
	data, err := prspy.FetchData()
	assert.NoError(t, err)
	assert.NotNil(t, data)

	players := prspy.GetAllPlayers(data)
	assert.NotNil(t, players)
}
