package dao

import (
	"context"
	"github.com/ecodeclub/ekit/sqlx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

const (
	asyncStatusWaiting = iota
	asyncStatusFailed
	asyncStatusSuccess
)

var _ AsyncSMSDAO = &GORMAsyncSMSDAO{}

type AsyncSMSDAO interface {
	Insert(ctx context.Context, info SMSAsyncInfo) error
	GetWaitingSMS(ctx context.Context) (SMSAsyncInfo, error)
	MarkSuccess(ctx context.Context, id int64) error
	MarkFailed(ctx context.Context, id int64) error
}

type GORMAsyncSMSDAO struct {
	db *gorm.DB
}

func NewGORMAsyncSMSDAO(db *gorm.DB) *GORMAsyncSMSDAO {
	return &GORMAsyncSMSDAO{
		db: db,
	}
}

func (g *GORMAsyncSMSDAO) Insert(ctx context.Context, info SMSAsyncInfo) error {
	now := time.Now().UnixMilli()
	info.Ctime = now
	info.Utime = now
	return g.db.WithContext(ctx).Create(&info).Error
}

func (g *GORMAsyncSMSDAO) GetWaitingSMS(ctx context.Context) (SMSAsyncInfo, error) {
	var s SMSAsyncInfo
	err := g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UnixMilli()
		endTime := now - time.Minute.Milliseconds()
		// 占有更新锁，防止其他事务修改
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("utime < ? AND status = ?", endTime, asyncStatusWaiting).
			First(&s).Error
		if err != nil {
			return err
		}

		err = tx.Model(&SMSAsyncInfo{}).
			Where("id = ?", s.Id).
			Updates(map[string]any{
				// 使用expr 的方式保证了操作的原子性
				"retry_cnt": gorm.Expr("retry_cnt + 1"),
				"utime":     now,
			}).Error
		return err
	})
	return s, err
}

func (g *GORMAsyncSMSDAO) MarkSuccess(ctx context.Context, id int64) error {
	now := time.Now().UnixMilli()
	return g.db.Model(&SMSAsyncInfo{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status": asyncStatusSuccess,
			"utime":  now,
		}).Error
}

func (g *GORMAsyncSMSDAO) MarkFailed(ctx context.Context, id int64) error {
	now := time.Now().UnixMilli()
	return g.db.Model(&SMSAsyncInfo{}).
		// 达到重试次数后才更新
		Where("id = ? AND `retry_cnt` >= `retry_max`", id).
		Updates(map[string]any{
			"status": asyncStatusFailed,
			"utime":  now,
		}).Error
}

type SMSAsyncInfo struct {
	Id       int64
	Config   sqlx.JsonColumn[SmsConfig]
	RetryCnt int
	RetryMax int
	Status   uint8
	Ctime    int64
	Utime    int64 `gorm:"index"`
}

type SmsConfig struct {
	Tpl     string
	Args    []string
	Numbers []string
}
