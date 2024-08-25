package model

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

type LanguageCode string
type AttachmentType string

const (
	Video AttachmentType = "video"
	Audio AttachmentType = "audio"
	File  AttachmentType = "file"
)

type GreatMan struct {
	gorm.Model
	HeadImgUrl    *string
	GreatManInfos []GreatManInfo `gorm:"foreignKey:GreatManId"`
	TalkRecords   []TalkRecord   `gorm:"foreignKey:GreatManId"`
}

type GreatManInfo struct {
	GreatManId uint   `gorm:"primaryKey;autoIncrement:false"`
	Language   string `gorm:"primaryKey"`
	Name       string
	Comment    *string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  sql.NullTime `gorm:"index"`
}

type TalkRecord struct {
	gorm.Model
	Type         string        `gorm:"index:idx_t"`
	GreatManId   uint          `gorm:"index:idx_d"`
	TalkContents []TalkContent `gorm:"foreignKey:TalkRecordId"`
	Attachments  []Attachment  `gorm:"foreignKey:TalkRecordId"`
}

type TalkContent struct {
	TalkRecordId uint   `gorm:"primaryKey;autoIncrement:false"`
	Language     string `gorm:"primaryKey"`
	Title        string
	MainBody     string
	Interviewer  *string
	Source       *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    sql.NullTime `gorm:"index"`
}

type Attachment struct {
	gorm.Model
	TalkRecordId uint `gorm:"index:idx_a"`
	Type         AttachmentType
	Value        string
}
