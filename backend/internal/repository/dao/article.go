package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
}

type GORMArticleDAO struct {
	db *gorm.DB
}

func NewGORMArticleDAO(db *gorm.DB) ArticleDAO {
	return &GORMArticleDAO{
		db: db,
	}
}

func (g *GORMArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.CTime = now
	art.UTime = now
	err := g.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

type Article struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Tittle   string `gorm:"type=varchar(4096)"`
	Content  string `gorm:"type=BLOB"`
	AuthorId int64  `gorm:"index"`
	CTime    int64
	UTime    int64
}
