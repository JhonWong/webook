package dao

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
	Sync(ctx context.Context, art Article) (int64, error)
	Upsert(ctx context.Context, art PublishArticle) error
	SyncStatus(ctx context.Context, id, usrId int64, status uint8) error
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
			"status":  art.Status,
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
	err = g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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
		art.Id = id
		return g.Upsert(ctx, PublishArticle{Article: art})
	})
	return id, err
}

func (g *GORMArticleDAO) Upsert(ctx context.Context, art PublishArticle) error {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	return g.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"tittle":  art.Tittle,
				"content": art.Content,
				"utime":   art.Utime,
				"status":  art.Status,
			}),
		}).Create(&art).Error
}

func (g *GORMArticleDAO) SyncStatus(ctx context.Context, id, usrId int64, status uint8) error {
	now := time.Now().UnixMilli()
	err := g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&Article{}).
			Where("id = ? AND author_id = ?", id, usrId).
			Updates(map[string]any{
				"status": status,
				"utime":  now,
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return fmt.Errorf("可能有人在攻击系统，误操作非自己的文章, Uid:%d, authorId:", id, usrId)
		}

		return tx.Model(&PublishArticle{}).
			Where("id = ?", id).
			Updates(map[string]any{
				"status": status,
				"utime":  now,
			}).Error
	})
	return err
}

type PublishArticle struct {
	Article
}

type Article struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Tittle   string `gorm:"type=varchar(4096)"`
	Content  string `gorm:"type=BLOB"`
	AuthorId int64  `gorm:"index"`
	Status   uint8
	Ctime    int64
	Utime    int64
}
