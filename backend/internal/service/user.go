package service

import (
	"context"
	"errors"

	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/johnwongx/webook/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicateEmail    = repository.ErrUserDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("账号/邮箱或密码不对")
)

type UserService interface {
	SignUp(ctx context.Context, u domain.User) error
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	Login(ctx context.Context, email, password string) (domain.User, error)
	Edit(ctx context.Context, id int64, nickName, birthday, selfIntro string) error
	Profile(ctx context.Context, id int64) (domain.User, error)
}

type userService struct {
	r repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{
		r: r,
	}
}

func (svc *userService) SignUp(ctx context.Context, u domain.User) error {
	//加密密码
	hash, err := bcrypt.GenerateFromPassword([]byte(u.PassWord), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PassWord = string(hash)
	return svc.r.Create(ctx, u)
}

func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	u, err := svc.r.FindByPhone(ctx, phone)
	if err != repository.ErrUserNotFound {
		return u, err
	}

	err = svc.r.Create(ctx, domain.User{
		Phone: phone,
	})
	if err != nil && err != repository.ErrUserDuplicate {
		return u, err
	}

	return svc.r.FindByPhone(ctx, phone)
}

func (svc *userService) Login(ctx context.Context, email, password string) (domain.User, error) {
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

func (svc *userService) Edit(ctx context.Context, id int64, nickName, birthday, selfIntro string) error {
	user, err := svc.r.FindById(ctx, id)
	if err != nil {
		return err
	}

	user.NickName = nickName
	user.Birthday = birthday
	user.SelfIntroduction = selfIntro
	return svc.r.Edit(ctx, user)
}

func (svc *userService) Profile(ctx context.Context, id int64) (domain.User, error) {
	user, err := svc.r.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return user, err
}
