package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

type AuthorArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
}

type GORMAuthorArticleDAO struct {
	db *gorm.DB
}

func NewGORMAuthorArticleDAO(db *gorm.DB) AuthorArticleDAO {
	return &GORMAuthorArticleDAO{
		db: db,
	}
}

func (g *GORMAuthorArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := g.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

func (g *GORMAuthorArticleDAO) UpdateById(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	art.Utime = now
	res := g.db.Model(&Article{}).WithContext(ctx).
		Where("id=? AND author_id=?", art.Id, art.AuthorId).
		Updates(map[string]any{
			"tittle":  art.Tittle,
			"content": art.Content,
			"utime":   art.Utime,
		})
	err := res.Error
	if err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return errors.New("更新数据失败")
	}
	return nil
}
