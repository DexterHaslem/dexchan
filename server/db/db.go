// Package db abstracts all database operations for dexchan
package db

import (
	"dexchan/server/cfg"
	_ "github.com/lib/pq"
	"database/sql"
	"fmt"
	"dexchan/server/model"
	"time"
)

type D interface {
	GetUserID(ip string) (int64, error)
	GetBoards() ([]*model.Board, error)
	CreatePost(p *model.Post, ip string) (int64, error)
	CreateThread(t *model.Thread, ip string) (int64, error)
	Close() error
}

type db struct {
	openedDB *sql.DB
}

func (d *db) GetUserID(ip string) (int64, error) {
	var userID int64
	err := d.openedDB.QueryRow("SELECT get_auser($1);", ip).Scan(&userID)
	return userID, err
}

func (d *db) GetBoards() ([]*model.Board, error) {
	q, err := d.openedDB.Query("SELECT id, name, shortname, description, nsfw, max_attachment_size, allowed_attachment_exts FROM board")
	if err != nil {
		return nil, err
	}
	defer q.Close()

	ret := []*model.Board{}
	for q.Next() {
		b := &model.Board{}

		err = q.Scan(&b.ID, &b.Name, &b.ShortCode, &b.Description,
			&b.IsNotWorksafe, &b.MaxAttachmentSize, &b.AttachmentTypes)

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

func (d *db) CreatePost(p *model.Post, ip string) (int64, error) {
	postedByID, err := d.GetUserID(ip)
	if err != nil {
		return 0, err
	}

	var createdID int64
	// post can have no attachment, we want to store nulls in this case,
	// so handle this. kind of a pain because go is too cool to ternary
	if p.Attachment != nil {
		err = d.openedDB.QueryRow("INSERT INTO post (thread_id, content, posted_at, posted_by_id, hidden, attachment_orig_name, attachment_tn_loc, attachment_loc, attachment_size) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id;",
			p.ThreadID, p.Content, time.Now(), postedByID, false,
			p.OriginalFilename, p.ThumbnailLocation, p.Location, p.Size).Scan(&createdID)
	} else {
		err = d.openedDB.QueryRow("INSERT INTO post (thread_id, content, posted_at, posted_by_id, hidden) VALUES ($1, $2, $3, $4, $5) RETURNING id;",
			p.ThreadID, p.Content, time.Now(), postedByID, false).Scan(&createdID)
	}
	return createdID, err
}

func (d *db) CreateThread(t *model.Thread, ip string) (int64, error) {
	postedByID, err := d.GetUserID(ip)
	if err != nil {
		return 0, err
	}
	var createdID int64
	err = d.openedDB.QueryRow("INSERT INTO thread (board_id, subject, description, created_at, created_by_id, hidden, attachment_orig_name, attachment_tn_loc, attachment_loc, attachment_size) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id;",
		t.BoardID, t.Subject, t.Description, time.Now(), postedByID, false,
		t.OriginalFilename, t.ThumbnailLocation, t.Location, t.Size).Scan(&createdID)
	return createdID, err
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
