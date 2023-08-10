package repository

import (
	"github.com/JhonWong/webook/backend/internal/domain"
	"github.com/JhonWong/webook/backend/internal/repository/dao"
	"github.com/gin-gonic/gin"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
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

	return domain.User{
		Id:       user.Id,
		Email:    user.Email,
		PassWord: user.Password,
		CTime:    user.CTime,
	}, nil
}

func (r *UserRepository) FindById(ctx *gin.Context, id int64) (domain.User, error) {
	// SELECT * FROM `users` WHERE `id`=?
	user, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		Id:       user.Id,
		Email:    user.Email,
		PassWord: user.Password,
		CTime:    user.CTime,
	}, nil
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
