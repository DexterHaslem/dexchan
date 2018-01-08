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


package model

import (
	"time"
	"mime/multipart"
	"path/filepath"
)

type AUser struct {
	ID        int64     `json:"id"`
	IPAddress string    `json:"ip"`
	FirstSeen time.Time `json:"firstSeen"`
	LastSeen  time.Time `json:"firstSeen"`
}

type Login struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash"`
	IsEnabled    bool   `json:"isEnabled"`
	IsAdmin      bool   `json:"isAdmin"`
}

type Board struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	ShortCode         string `json:"shortCode"`
	Description       string `json:"description"`
	IsNotWorksafe     bool   `json:"isNsfw"`
	IsHidden          bool
	MaxAttachmentSize int64  `json:"maxAttachmentSize"`
	AttachmentTypes   string `json:"allowedAttachmentTypes"`
}

type TemplateHelper struct {
	IsVideo       bool
	HasAttachment bool
}

type Thread struct {
	ID          int64     `json:"id"`
	BoardID     int64     `json:"boardID"`
	Description string    `json:"description"`
	Subject     string    `json:"subject"`
	PostedByID  int64     `json:"postedByID"`
	CreatedAt   time.Time `json:"createdAt"`
	Attachment
	TemplateHelper
}

type Post struct {
	ID         int64     `json:"ID"`
	ThreadID   int64     `json:"threadID"`
	Content    string    `json:"content"`
	PostedAt   time.Time `json:"postedAt"`
	IsHidden   bool      `json:"isHidden"`
	PostedByID int64     `json:"postedByID"`
	Attachment
	TemplateHelper
}

type Attachment struct {
	OriginalFilename  string `json:"attachmentOriginalFilename"`
	Location          string `json:"attachmentLocation"`
	ThumbnailLocation string `json:"attachmentTnLocation"`
	Size              int64  `json:"attachmentSize"`
}

type AttachmentEntity interface {
	SetThumbnail(string)
	SetLocation(string)
	GetLocation() string
	GetThumbnail() string
	ParseFromHeader(*multipart.FileHeader)
}

func (a *Attachment) ParseFromHeader(h *multipart.FileHeader) {
	a.OriginalFilename = filepath.Base(h.Filename)
	a.Size = h.Size
}

func (a *Attachment) SetThumbnail(tn string) {
	a.ThumbnailLocation = tn
}
func (a *Attachment) SetLocation(loc string) {
	a.Location = loc
}

func (a *Attachment) GetLocation() string {
	return a.Location
}

func (a *Attachment) GetThumbnail() string {
	return a.ThumbnailLocation
}
