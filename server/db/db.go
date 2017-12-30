// Package db abstracts all database operations for dexchan
package db

import (
	"dexchan/server/cfg"
	_ "github.com/lib/pq"
	"database/sql"
	"fmt"
	"dexchan/server/model"
)

type D interface {
	GetBoards() ([]*model.Board, error)
	Close() error
}

type db struct {
	openedDB *sql.DB
}

func (d *db) GetBoards() ([]*model.Board, error) {
	q, err := d.openedDB.Query("select id, name, shortname, description, nsfw from board")
	if err != nil {
		return nil, err
	}
	defer q.Close()

	ret := []*model.Board{}
	for q.Next() {
		b := &model.Board{}

		err = q.Scan(&b.ID, &b.Name, &b.ShortCode, &b.Description, &b.IsNotWorksafe)
		if err != nil {
			return nil, err
		}

		ret = append(ret, b)
	}

	return ret, nil
}

func (d *db) Close() error {
	return d.openedDB.Close()
}

func Open(c *cfg.C) (D, error) {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%d",
		c.DbUsername, c.DbPassword, c.DbName, c.DbHost, c.DbPort)

	ret := &db{}
	odb, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return nil, err
	}
	ret.openedDB = odb
	return ret, nil
}