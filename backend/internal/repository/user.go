package repository

import (
	"database/sql"
	"github.com/JhonWong/webook/backend/internal/domain"
	"github.com/JhonWong/webook/backend/internal/repository/cache"
	"github.com/JhonWong/webook/backend/internal/repository/dao"
	"github.com/gin-gonic/gin"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
	ErrUserDuplicate      = dao.ErrUserDuplicate
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDAO, c *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: c,
	}
}

func (r *UserRepository) Create(ctx *gin.Context, u domain.User) error {
	return r.dao.Insert(ctx, r.domainToEntity(u))
}

func (r *UserRepository) FindByEmail(ctx *gin.Context, email string) (domain.User, error) {
	// SELECT * FROM `users` WHERE `email`=?
	user, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}

	return r.entityToDomain(user), nil
}

func (r *UserRepository) FindByPhone(ctx *gin.Context, phone string) (domain.User, error) {
	// SELECT * FROM `users` WHERE `phone`=?
	user, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}

	return r.entityToDomain(user), nil
}

func (r *UserRepository) FindById(ctx *gin.Context, id int64) (domain.User, error) {
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

func (r *UserRepository) Edit(ctx *gin.Context, u domain.User) error {
	return r.dao.Update(ctx, r.domainToEntity(u))
}

func (r *UserRepository) domainToEntity(u domain.User) dao.User {
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
		Password:         u.PassWord,
		CTime:            u.CTime,
		NickName:         u.NickName,
		Birthday:         u.Birthday,
		SelfIntroduction: u.SelfIntroduction,
	}
}

func (r *UserRepository) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:               u.Id,
		Email:            u.Email.String,
		Phone:            u.Phone.String,
		PassWord:         u.Password,
		CTime:            u.CTime,
		NickName:         u.NickName,
		Birthday:         u.Birthday,
		SelfIntroduction: u.SelfIntroduction,
	}
}
