package dao

import (
	"context"
	"gorm.io/gorm"
)

type ReaderArticleDAO interface {
	Upsert(ctx context.Context, art Article) error
}

type GORMReaderArticleDAO struct {
	db *gorm.DB
}

func NewGORMReaderArticleDAO(db *gorm.DB) ReaderArticleDAO {
	return &GORMReaderArticleDAO{
		db: db,
	}
}

func (g *GORMReaderArticleDAO) Upsert(ctx context.Context, art Article) error {
	//TODO implement me
	panic("implement me")
	return nil
}
