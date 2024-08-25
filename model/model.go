package model

import "gorm.io/gorm"

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
	TalkRecord    []TalkRecord   `gorm:"foreignKey:GreatManId"`
}

type GreatManInfo struct {
	gorm.Model
	GreatManId uint   `gorm:"index:idx_d,priority:2"`
	Language   string `gorm:"index:idx_d,priority:1"`
	Name       string
	Comment    *string
}

type TalkRecord struct {
	gorm.Model
	Type         string        `gorm:"index:idx_t"`
	GreatManId   uint          `gorm:"index:idx_d"`
	TalkContents []TalkContent `gorm:"foreignKey:TalkRecordId"`
	Attachments  []Attachment  `gorm:"foreignKey:TalkRecordId"`
}

type TalkContent struct {
	gorm.Model
	TalkRecordId uint   `gorm:"index:idx_c,priority:2"`
	Language     string `gorm:"index:idx_c,priority:1"`
	Title        string
	MainBody     string
	Interviewer  *string
	Source       *string
}

type Attachment struct {
	gorm.Model
	TalkRecordId uint `gorm:"index:idx_a"`
	Type         AttachmentType
	Value        string
}
