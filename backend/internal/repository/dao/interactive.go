package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type InteractiveDAO interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	Get(ctx context.Context, biz string, bizId int64) (Interactive, error)
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
				"read_cnt": "read_cnt+1",
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

type Interactive struct {
	Id      int64  `gorm:"primaryKey,autoIncrement"`
	BizId   int64  `gorm:"uniqueIndex:biz_id_type"`
	Biz     string `gorm:"uniqueIndex:biz_id_type;type:varchar(128)"`
	ReadCnt int64
	Utime   int64
	Ctime   int64
}
