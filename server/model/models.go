package model

import "time"

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
	MaxAttachmentSize int64  `json:"maxAttachmentSize"`
	AttachmentTypes   string `json:"allowedAttachmentTypes"`
}

type Thread struct {
	ID          int64     `json:"id"`
	BoardID     int64     `json:"boardID"`
	Description string    `json:"description"`
	Subject     string    `json:"subject"`
	PostedByID  int64     `json:"postedByID"`
	CreatedAt   time.Time `json:"createdAt"`
	Attachment
}

type Post struct {
	ID         int64     `json:"ID"`
	ThreadID   int64     `json:"threadID"`
	Content    string    `json:"content"`
	PostedAt   time.Time `json:"postedAt"`
	IsHidden   bool      `json:"isHidden"`
	PostedByID int64     `json:"postedByID"`
	Attachment

	// template helper
	HasAttachment bool
}

type Attachment struct {
	OriginalFilename  string `json:"attachmentOriginalFilename"`
	Location          string `json:"attachmentLocation"`
	ThumbnailLocation string `json:"attachmentTnLocation"`
	Size              int64  `json:"attachmentSize"`
}
