package common

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"greatmanbackend/model"
	"os"
)

var DB *gorm.DB

func InitDB() {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: os.Getenv("DATABASE_URL"),
	}))

	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&model.GreatMan{}, &model.GreatManInfo{}, &model.TalkRecord{}, &model.TalkContent{}, &model.Attachment{}); err != nil {
		panic(err)
	}

	DB = db
}

func GetDB() *gorm.DB {
	return DB
}
