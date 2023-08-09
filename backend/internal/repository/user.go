package repository

import (
	"github.com/JhonWong/webook/backend/internal/domain"
	"github.com/JhonWong/webook/backend/internal/repository/dao"
	"github.com/gin-gonic/gin"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
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
