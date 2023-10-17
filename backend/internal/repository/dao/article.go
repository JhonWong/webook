package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
	Sync(ctx context.Context, art Article) (int64, error)
	Upsert(ctx context.Context, art Article) error
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
	art.Ctime = now
	art.Utime = now
	err := g.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

func (g *GORMArticleDAO) UpdateById(ctx context.Context, art Article) error {
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

func (g *GORMArticleDAO) Sync(ctx context.Context, art Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)

	// 使用事物保证两张表同时成功或失败
	err = g.db.Transaction(func(tx *gorm.DB) error {
		// 更新制作库，插入或删除
		if id > 0 {
			err = g.UpdateById(ctx, art)
		} else {
			id, err = g.Insert(ctx, art)
		}
		if err != nil {
			return err
		}
		// 更新数据到线上库
		return g.Upsert(ctx, art)
	})
	return id, err
}

func (g *GORMArticleDAO) Upsert(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	return g.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"tittle":  art.Tittle,
				"content": art.Content,
				"utime":   art.Utime,
			}),
		}).Create(&art).Error
}

type Article struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Tittle   string `gorm:"type=varchar(4096)"`
	Content  string `gorm:"type=BLOB"`
	AuthorId int64  `gorm:"index"`
	Ctime    int64
	Utime    int64
}
