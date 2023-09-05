package service

import (
	"context"
	"errors"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/johnwongx/webook/backend/internal/repository"
	repomocks "github.com/johnwongx/webook/backend/internal/repository/mocks"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService_Login(t *testing.T) {
	testCase := []struct {
		name     string
		email    string
		passWord string
		daoFunc  func(ctrl *gomock.Controller) repository.UserRepository
		wantUser domain.User
		wantErr  error
	}{
		{
			name:     "success",
			email:    "123@qq.com",
			passWord: "123",
			daoFunc: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				hashPasswod, _ := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{
					Email:    "123@qq.com",
					PassWord: string(hashPasswod),
				}, nil)
				return repo
			},
			wantErr:  nil,
			wantUser: domain.User{Email: "123@qq.com"},
		},
		{
			name:     "not found",
			email:    "123@qq.com",
			passWord: "123",
			daoFunc: func(ctrl *gomock.Controller) *repomocks.MockUserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				hashPasswod, _ := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{
					Email:    "123@qq.com",
					PassWord: string(hashPasswod),
				}, repository.ErrUserNotFound)
				return repo
			},
			wantErr:  ErrInvalidUserOrPassword,
			wantUser: domain.User{Email: "123@qq.com"},
		},
		{
			name:     "system error",
			email:    "123@qq.com",
			passWord: "123",
			daoFunc: func(ctrl *gomock.Controller) *repomocks.MockUserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				hashPasswod, _ := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{
					Email:    "123@qq.com",
					PassWord: string(hashPasswod),
				}, errors.New("system error"))
				return repo
			},
			wantErr:  errors.New("system error"),
			wantUser: domain.User{Email: "123@qq.com"},
		},
		{
			name:     "success",
			email:    "123@qq.com",
			passWord: "123",
			daoFunc: func(ctrl *gomock.Controller) *repomocks.MockUserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				hashPasswod, _ := bcrypt.GenerateFromPassword([]byte("1234"), bcrypt.DefaultCost)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{
					Email:    "123@qq.com",
					PassWord: string(hashPasswod),
				}, nil)
				return repo
			},
			wantErr:  ErrInvalidUserOrPassword,
			wantUser: domain.User{Email: "123@qq.com"},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := tc.daoFunc(ctrl)
			us := NewUserService(repo)
			user, err := us.Login(context.Background(), tc.email, tc.passWord)
			assert.Equal(t, err, tc.wantErr)
			if err != nil {
				return
			}

			assert.Equal(t, user.Email, tc.wantUser.Email)
		})
	}
}
