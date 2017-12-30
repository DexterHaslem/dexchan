package model

import "time"

type AUser struct {
	ID        int       `json:"id"`
	IPAddress string    `json:"ip"`
	FirstSeen time.Time `json:"firstSeen"`
	LastSeen  time.Time `json:"firstSeen"`
}

type Login struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash"`
	IsEnabled    bool   `json:"isEnabled"`
	IsAdmin      bool   `json:"isAdmin"`
}

type Board struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	ShortCode     string `json:"shortCode"`
	Description   string `json:"description"`
	IsNotWorksafe bool   `json:"isNsfw"`
}

type Thread struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Subject     string `json:"subject"`
	PostedByID  int    `json:"postedByID"`
}

type Post struct {
	ID         int       `json:"ID"`
	Content    string    `json:"content"`
	PostedAt   time.Time `json:"postedAt"`
	IsHidden   bool      `json:"isHidden"`
	PostedByID int       `json:"postedByID"`
}

type Attachment struct {
	ID                int    `json:"id"`
	AttachmentTypeID  int    `json:"attachmentTypeID"`
	OriginalFilename  string `json:"originalFilename"`
	UploadedByID      int    `json:"uploadedByID"`
	Location          string `json:"location"`
	ThumbnailLocation string `json:"tnLocation"`
}

type PostWithAttachment struct {
	Post
	Attachment
}

type AttachmentType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Ext  string `json:"ext"`
	//	MaxSizeInBytes int `json:"maxSizeBytes"`
}
