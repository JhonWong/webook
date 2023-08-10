package service

import (
	"errors"

	"github.com/JhonWong/webook/backend/internal/domain"
	"github.com/JhonWong/webook/backend/internal/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicateEmail    = repository.ErrUserDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("账号/邮箱或密码不对")
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
	hash, err := bcrypt.GenerateFromPassword([]byte(u.PassWord), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PassWord = string(hash)
	return svc.r.Create(ctx, u)
}

func (svc *UserService) Login(ctx *gin.Context, email, password string) (domain.User, error) {
	user, err := svc.r.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PassWord), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}

	return user, nil
}

func (svc *UserService) Edit(ctx *gin.Context, id int64, nickName, birthday, selfIntro string) error {
	user, err := svc.r.FindById(ctx, id)
	if err != nil {
		return err
	}

	user.NickName = nickName
	user.Birthday = birthday
	user.SelfIntroduction = selfIntro
	return svc.r.Edit(ctx, user)
}

func (svc *UserService) Profile(ctx *gin.Context, id int64) (domain.User, error) {
	user, err := svc.r.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return user, err
}
