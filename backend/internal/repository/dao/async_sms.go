package dao

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
)

var _ AsyncSMSDAO = &GORMAsyncSMSDAO{}

type AsyncSMSDAO interface {
	Store(ctx context.Context, info SMSAsyncInfo) error
	Load(ctx context.Context) (SMSAsyncInfo, error)
	UpdateResult(ctx context.Context, id int64, res bool) error
}

type GORMAsyncSMSDAO struct {
	db *gorm.DB
}

func NewGORMAsyncSMSDAO(db *gorm.DB) *GORMAsyncSMSDAO {
	return &GORMAsyncSMSDAO{
		db: db,
	}
}

func (g *GORMAsyncSMSDAO) Store(ctx context.Context, info SMSAsyncInfo) error {
	return g.db.WithContext(ctx).Create(&info).Error
}

func (g *GORMAsyncSMSDAO) Load(ctx context.Context) (SMSAsyncInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GORMAsyncSMSDAO) UpdateResult(ctx context.Context, id int64, res bool) error {
	//TODO implement me
	panic("implement me")
}

// TODO
type SMSAsyncInfo struct {
	Id            int64 `gorm:"primaryKey,autoIncrement"`
	Tpl           sql.NullString
	Args          []string
	Numbers       []string
	MaxRetryCount int
}
