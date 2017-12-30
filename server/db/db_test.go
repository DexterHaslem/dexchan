package db_test

import (
	"testing"
	"dexchan/server/db"
	"dexchan/server/cfg"
	"github.com/stretchr/testify/assert"
	"dexchan/server/model"
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

	newThread := &model.Thread{
		Subject: "hello world",
		Description: "hopefully a thread created",
	}

	tid, err := d.CreateThread(boards[0].ID, newThread, "192.168.1.1");
	assert.NoError(t, err)
	assert.NotEqual(t, 0, tid)
	assert.NoError(t, d.Close())
}
