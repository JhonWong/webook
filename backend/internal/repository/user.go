package repository

import (
	"context"
	"database/sql"
	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/johnwongx/webook/backend/internal/repository/cache"
	"github.com/johnwongx/webook/backend/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
	ErrUserDuplicate      = dao.ErrUserDuplicate
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindByWechat(ctx context.Context, info domain.WechatInfo) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
	Edit(ctx context.Context, u domain.User) error
}

type CachedUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDAO, c cache.UserCache) UserRepository {
	return &CachedUserRepository{
		dao:   dao,
		cache: c,
	}
}

func (r *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, r.domainToEntity(u))
}

func (r *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	// SELECT * FROM `users` WHERE `email`=?
	user, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}

	return r.entityToDomain(user), nil
}

func (r *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	// SELECT * FROM `users` WHERE `phone`=?
	user, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}

	return r.entityToDomain(user), nil
}

func (r *CachedUserRepository) FindByWechat(ctx context.Context, info domain.WechatInfo) (domain.User, error) {
	user, err := r.dao.FindByWechat(ctx, info.OpenID)
	if err != nil {
		return domain.User{}, err
	}

	return r.entityToDomain(user), nil
}

func (r *CachedUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.cache.Get(ctx, id)
	if err == nil {
		return u, err
	}

	//TODO 添加数据库限流，防止崩溃
	// SELECT * FROM `users` WHERE `id`=?
	ue, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	u = r.entityToDomain(ue)

	go func() {
		//更新缓存
		err = r.cache.Set(ctx, u)
		if err != nil {
			//TODO:添加监控日志
		}
	}()

	return u, nil
}

func (r *CachedUserRepository) Edit(ctx context.Context, u domain.User) error {
	return r.dao.Update(ctx, r.domainToEntity(u))
}

func (r *CachedUserRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		WechatUnionID: sql.NullString{
			String: u.WechatInfo.UnionID,
			Valid:  u.WechatInfo.UnionID != "",
		},
		WechatOpenID: sql.NullString{
			String: u.WechatInfo.OpenID,
			Valid:  u.WechatInfo.OpenID != "",
		},
		Password:         u.PassWord,
		CTime:            u.CTime,
		NickName:         u.NickName,
		Birthday:         u.Birthday,
		SelfIntroduction: u.SelfIntroduction,
	}
}

func (r *CachedUserRepository) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:    u.Id,
		Email: u.Email.String,
		Phone: u.Phone.String,
		WechatInfo: domain.WechatInfo{
			UnionID: u.WechatUnionID.String,
			OpenID:  u.WechatOpenID.String,
		},
		PassWord:         u.Password,
		CTime:            u.CTime,
		NickName:         u.NickName,
		Birthday:         u.Birthday,
		SelfIntroduction: u.SelfIntroduction,
	}
}
