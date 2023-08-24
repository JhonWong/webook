package repository

import (
	"github.com/JhonWong/webook/backend/internal/domain"
	"github.com/JhonWong/webook/backend/internal/repository/cache"
	"github.com/JhonWong/webook/backend/internal/repository/dao"
	"github.com/gin-gonic/gin"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
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
	return r.dao.Insert(ctx, dao.User{
		Email:    string(u.Email),
		Password: string(u.PassWord),
	})
}

func (r *UserRepository) FindByEmail(ctx *gin.Context, email string) (domain.User, error) {
	// SELECT * FROM `users` WHERE `email`=?
	user, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}

	return convertDaoUser2DomainUser(user), nil
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

	u = convertDaoUser2DomainUser(ue)

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
	return r.dao.Update(ctx, dao.User{
		Id:               u.Id,
		Email:            u.Email,
		Password:         u.PassWord,
		NickName:         u.NickName,
		Birthday:         u.Birthday,
		SelfIntroduction: u.SelfIntroduction,
	})
}

func convertDaoUser2DomainUser(u dao.User) domain.User {
	return domain.User{
		Id:               u.Id,
		Email:            u.Email,
		PassWord:         u.Password,
		CTime:            u.CTime,
		NickName:         u.NickName,
		Birthday:         u.Birthday,
		SelfIntroduction: u.SelfIntroduction,
	}
}
