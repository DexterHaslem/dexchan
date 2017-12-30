package db_test

import (
	"testing"
	"dexchan/server/db"
	"dexchan/server/cfg"
	"github.com/stretchr/testify/assert"
)

func TestDb_GetBoards(t *testing.T) {
	c := cfg.C{
		DbName:     "dexchan",
		DbHost:     "localhost",
		DbPort:     5432,
		DbPassword: "dexchan",
		DbUsername: "dexchan",
		Port:       1234,
	}

	d, err := db.Open(&c)
	if assert.NoError(t, err) {
		assert.NotNil(t, d)
	}

	boards, err := d.GetBoards()
	assert.NoError(t, err)
	assert.NotNil(t, boards)

	for _, b := range boards {
		assert.NotEqual(t, b.Name, "")
	}

	assert.NoError(t, d.Close())
}
