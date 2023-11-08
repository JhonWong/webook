package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

var ErrDataNotFound = gorm.ErrRecordNotFound

type InteractiveDAO interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	IncrLike(ctx context.Context, id int64, biz string, uid int64) error
	DecrLike(ctx context.Context, id int64, biz string, uid int64) error
	InsertCollectionBiz(ctx context.Context, id int64, biz string, cid int64, uid int64) error
	Get(ctx context.Context, biz string, bizId int64) (Interactive, error)
	GetLikeInfo(ctx context.Context, biz string, id int64, uid int64) (UserLikeBiz, error)
	GetCollectInfo(ctx context.Context, biz string, id int64, uid int64) (UserCollectBiz, error)
}

type GORMInteractiveDAO struct {
	db *gorm.DB
}

func NewGORMInteractiveDAO(db *gorm.DB) InteractiveDAO {
	return &GORMInteractiveDAO{
		db: db,
	}
}

func (g *GORMInteractiveDAO) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"read_cnt": gorm.Expr("read_cnt+1"),
				"utime":    now,
			}),
		}).Create(&Interactive{
		BizId:   bizId,
		Biz:     biz,
		ReadCnt: 1,
		Utime:   now,
		Ctime:   now,
	}).Error
}

func (g *GORMInteractiveDAO) Get(ctx context.Context, biz string, bizId int64) (Interactive, error) {
	var res Interactive
	err := g.db.WithContext(ctx).
		Where("biz_id = ? AND biz = ?", bizId, biz).
		First(&res).Error
	return res, err
}

func (g *GORMInteractiveDAO) IncrLike(ctx context.Context, id int64, biz string, uid int64) error {
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UnixMilli()
		err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"status": 1,
				"utime":  now,
			}),
		}).Create(&UserLikeBiz{
			BizId:  id,
			Biz:    biz,
			Status: 1,
			UserId: uid,
			Ctime:  now,
			Utime:  now,
		}).Error
		if err != nil {
			return err
		}

		return tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"like_cnt": gorm.Expr("`like_cnt`+1"),
				"utime":    now,
			}),
		}).Create(&Interactive{
			BizId:   id,
			Biz:     biz,
			LikeCnt: 1,
			Ctime:   now,
			Utime:   now,
		}).Error
	})
}

func (g *GORMInteractiveDAO) DecrLike(ctx context.Context, id int64, biz string, uid int64) error {
	return g.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now().UnixMilli()
		err := tx.Model(&UserLikeBiz{}).
			Where("biz_id = ? AND biz = ? AND user_id = ?", id, biz, uid).
			Updates(map[string]any{
				"status": 0,
				"utime":  now,
			}).Error
		if err != nil {
			return err
		}

		return tx.Model(&Interactive{}).
			Where("biz_id = ? AND biz = ?", id, biz).
			Updates(map[string]any{
				"like_cnt": gorm.Expr("`like_cnt`-1"),
				"utime":    now,
			}).Error
	})
}

func (g *GORMInteractiveDAO) InsertCollectionBiz(ctx context.Context, id int64, biz string, cid int64, uid int64) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&UserCollectBiz{
			BizId:  id,
			Biz:    biz,
			Cid:    cid,
			UserId: uid,
			Ctime:  now,
			Utime:  now,
		}).Error
		if err != nil {
			return err
		}

		return tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"collect_cnt": gorm.Expr("`collect_cnt`+1"),
				"utime":       now,
			}),
		}).Create(&Interactive{
			BizId:      id,
			Biz:        biz,
			CollectCnt: 1,
			Ctime:      now,
			Utime:      now,
		}).Error
	})
}

func (g *GORMInteractiveDAO) GetLikeInfo(ctx context.Context, biz string, id int64, uid int64) (UserLikeBiz, error) {
	var info UserLikeBiz
	err := g.db.WithContext(ctx).
		Where("biz_id = ? AND biz = ? AND user_id = ?", id, biz, uid).
		First(&info).Error
	return info, err
}

func (g *GORMInteractiveDAO) GetCollectInfo(ctx context.Context, biz string, id int64, uid int64) (UserCollectBiz, error) {
	var info UserCollectBiz
	err := g.db.WithContext(ctx).
		Where("biz_id = ? AND biz = ? AND user_id = ?", id, biz, uid).
		First(&info).Error
	return info, err
}

type UserCollectBiz struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`

	Cid    int64  `gorm:"index"`
	BizId  int64  `gorm:"uniqueIndex:biz_id_uid"`
	Biz    string `gorm:"uniqueIndex:biz_id_uid;type:varchar(128)"`
	UserId int64  `gorm:"uniqueIndex:biz_id_uid"`

	Utime int64
	Ctime int64
}

type UserLikeBiz struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`

	BizId  int64  `gorm:"uniqueIndex:biz_id_type"`
	Biz    string `gorm:"uniqueIndex:biz_id_type;type:varchar(128)"`
	UserId int64  `gorm:"uniqueIndex:biz_id_type"`

	Status int64
	Utime  int64
	Ctime  int64
}

type Collection struct {
	Id     int64  `gorm:"primaryKey,autoIncrement"`
	Name   string `gorm:"type:varchar(1024)"`
	UserId int64

	Utime int64
	Ctime int64
}

type Interactive struct {
	Id         int64  `gorm:"primaryKey,autoIncrement"`
	BizId      int64  `gorm:"uniqueIndex:biz_id_type"`
	Biz        string `gorm:"uniqueIndex:biz_id_type;type:varchar(128)"`
	ReadCnt    int64
	LikeCnt    int64
	CollectCnt int64
	Utime      int64
	Ctime      int64
}
