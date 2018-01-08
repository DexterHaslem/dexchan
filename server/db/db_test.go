// dexChan copyright Dexter Haslem <dmh@fastmail.com> 2018
// This file is part of dexChan
//
// dexChan is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// dexChan is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with dexChan. If not, see <http://www.gnu.org/licenses/>.


package db_test

import (
	"testing"
	"dexchan/server/db"
	"dexchan/server/cfg"
	"github.com/stretchr/testify/assert"
	"dexchan/server/model"
)

func TestDb(t *testing.T) {
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
		BoardID:     boards[0].ID,
		Subject:     "hello world",
		Description: "hopefully a thread created",
		Attachment: model.Attachment{
			Size:              1024,
			ThumbnailLocation: "Thumbnaillocation",
			OriginalFilename:  "Orig filename.webm",
			Location:          "cdn location",
		},
	}

	tid, err := d.CreateThread(newThread, "192.168.1.1")
	if assert.NoError(t, err) {
		assert.NotEqual(t, 0, tid)

		p := &model.Post{
			ThreadID: tid,
			Content: "shitposting something with no attachment",
		}
		pid, err := d.CreatePost(p, "192.168.1.2")
		assert.NoError(t, err)
		assert.NotEqual(t, 0, pid)

		p2 := &model.Post{
			ThreadID: tid,
			Content: "shitposting something with attachment",
			Attachment: model.Attachment{
				Size:              1024,
				ThumbnailLocation: "Thumbnaillocation",
				OriginalFilename:  "Orig filename.webm",
				Location:          "cdn location",
			},
		}

		pid2, err := d.CreatePost(p2, "192.168.1.2")
		if assert.NoError(t, err) {
			assert.NotEqual(t, 0, pid2)

			posts, err := d.GetPosts(tid)
			assert.NoError(t, err)
			assert.True(t, len(posts) == 2)
		}

		threads, err := d.GetThreads(boards[0].ID)
		assert.NoError(t, err)
		assert.True(t, len(threads) >= 1)
	}
	assert.NoError(t, d.Close())
}
