package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/johnwongx/webook/backend/internal/repository/cache"
	cachemocks "github.com/johnwongx/webook/backend/internal/repository/cache/mocks"
	"github.com/johnwongx/webook/backend/internal/repository/dao"
	daomocks "github.com/johnwongx/webook/backend/internal/repository/dao/mocks"
	"go.uber.org/mock/gomock"
)

func TestCachedUserRepository_FindById(t *testing.T) {
	testCase := []struct {
		name      string
		id        int64
		cacheMock func(ctrl *gomock.Controller) cache.UserCache
		userMock  func(ctrl *gomock.Controller) dao.UserDAO
		wantErr   error
		wantUser  domain.User
	}{
		{
			name: "success",
			id:   1,
			cacheMock: func(ctrl *gomock.Controller) cache.UserCache {
				cm := cachemocks.NewMockUserCache(ctrl)
				cm.EXPECT().Get(gomock.Any(), int64(1)).Return(domain.User{}, errors.New("user not found"))
				return cm
			},
			userMock: func(ctrl *gomock.Controller) dao.UserDAO {
				um := daomocks.NewMockUserDAO(ctrl)
				um.EXPECT().FindById(gomock.Any(), int64(1)).Return(dao.User{
					Id: 1,
				}, nil)
				return um
			},
			wantErr:  nil,
			wantUser: domain.User{Id: 1},
		},
		{
			name: "find by cache",
			id:   1,
			cacheMock: func(ctrl *gomock.Controller) cache.UserCache {
				cm := cachemocks.NewMockUserCache(ctrl)
				cm.EXPECT().Get(gomock.Any(), int64(1)).Return(domain.User{Id: 1}, nil)
				return cm
			},
			userMock: func(ctrl *gomock.Controller) dao.UserDAO {
				return nil
			},
			wantErr:  nil,
			wantUser: domain.User{Id: 1},
		},
		{
			name: "not found",
			id:   1,
			cacheMock: func(ctrl *gomock.Controller) cache.UserCache {
				cm := cachemocks.NewMockUserCache(ctrl)
				cm.EXPECT().Get(gomock.Any(), int64(1)).Return(domain.User{}, errors.New("user not found"))
				return cm
			},
			userMock: func(ctrl *gomock.Controller) dao.UserDAO {
				um := daomocks.NewMockUserDAO(ctrl)
				um.EXPECT().FindById(gomock.Any(), int64(1)).Return(dao.User{
					Id: 1,
				}, errors.New("user not found"))
				return um
			},
			wantErr: errors.New("user not found"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cm := tc.cacheMock(ctrl)
			um := tc.userMock(ctrl)
			ur := NewUserRepository(um, cm)
			user, err := ur.FindById(context.Background(), tc.id)
			assert.Equal(t, err, tc.wantErr)
			if err != nil {
				return
			}
			assert.Equal(t, user.Id, tc.wantUser.Id)
		})
	}
}
