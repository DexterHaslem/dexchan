// Package db abstracts all database operations for dexchan
package db

import "dexchan/server/cfg"

type D interface {
	GetBoards()
}

func Open(c *cfg.C) {

}
