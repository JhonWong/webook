package service

import (
	"github.com/JhonWong/webook/backend/internal/domain"
	"github.com/JhonWong/webook/backend/internal/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicateEmail = repository.ErrUserDuplicateEmail
)

type UserService struct {
	r *repository.UserRepository
}

func NewUserService(r *repository.UserRepository) *UserService {
	return &UserService{
		r: r,
	}
}

func (svc *UserService) SignUp(ctx *gin.Context, u domain.User) error {
	//加密密码
	hash, err := bcrypt.GenerateFromPassword(u.PassWord, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PassWord = hash
	return svc.r.Create(ctx, u)
}
