package dao

import (
	"github.com/johnwongx/webook/backend/internal/repository/dao/article"
	"gorm.io/gorm"
)

func InitTable(db *gorm.DB) error {
	return db.AutoMigrate(&User{},
		&article.Article{},
		&article.PublishArticle{},
		&SMSAsyncInfo{})
}
