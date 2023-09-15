package repository

import (
	"context"
	"errors"
	"github.com/go-playground/assert/v2"
	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/johnwongx/webook/backend/internal/repository/cache"
	cachemocks "github.com/johnwongx/webook/backend/internal/repository/cache/mocks"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestSmsRepository_Put(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) cache.SMSCache
		info    domain.SMSInfo
		wantErr error
	}{
		{
			name: "success",
			mock: func(ctrl *gomock.Controller) cache.SMSCache {
				c := cachemocks.NewMockSMSCache(ctrl)
				c.EXPECT().Add(gomock.Any(), cache.SMSInfo{
					Tpl:     "123",
					Args:    []string{"123"},
					Numbers: []string{"1751234567"},
				}).Return(nil)
				return c
			},
			info: domain.SMSInfo{
				Tpl:     "123",
				Args:    []string{"123"},
				Numbers: []string{"1751234567"},
			},
		},
		{
			name: "failed",
			mock: func(ctrl *gomock.Controller) cache.SMSCache {
				c := cachemocks.NewMockSMSCache(ctrl)
				c.EXPECT().Add(gomock.Any(), cache.SMSInfo{
					Tpl:     "123",
					Args:    []string{"123"},
					Numbers: []string{"1751234567"},
				}).Return(errors.New("cache error"))
				return c
			},
			info: domain.SMSInfo{
				Tpl:     "123",
				Args:    []string{"123"},
				Numbers: []string{"1751234567"},
			},
			wantErr: errors.New("cache error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			c := tc.mock(ctrl)

			repo := NewSMSRepository(c)
			err := repo.Put(context.Background(), tc.info)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestSmsRepository_Get(t *testing.T) {
	testCases := []struct {
		name      string
		mock      func(ctrl *gomock.Controller) cache.SMSCache
		cnt       int
		wantInfos []domain.SMSInfo
		wantErr   error
	}{
		{
			name: "get 2 object",
			mock: func(ctrl *gomock.Controller) cache.SMSCache {
				c := cachemocks.NewMockSMSCache(ctrl)
				infoArr := []cache.SMSInfo{
					{
						Tpl:     "1234",
						Args:    []string{"1234"},
						Numbers: []string{"1751234567"},
					},
					{
						Tpl:     "123",
						Args:    []string{"123"},
						Numbers: []string{"175123456"},
					},
				}
				c.EXPECT().Take(gomock.Any(), 2).Return(infoArr, nil)
				return c
			},
			cnt: 2,
			wantInfos: []domain.SMSInfo{
				{
					Tpl:     "1234",
					Args:    []string{"1234"},
					Numbers: []string{"1751234567"},
				},
				{
					Tpl:     "123",
					Args:    []string{"123"},
					Numbers: []string{"175123456"},
				},
			},
		},
		{
			name: "cache error",
			mock: func(ctrl *gomock.Controller) cache.SMSCache {
				c := cachemocks.NewMockSMSCache(ctrl)
				c.EXPECT().Take(gomock.Any(), 2).Return(nil, errors.New("cache error"))
				return c
			},
			cnt:     2,
			wantErr: errors.New("cache error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			c := tc.mock(ctrl)

			repo := NewSMSRepository(c)
			res, err := repo.Get(context.Background(), tc.cnt)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantInfos, res)
		})
	}
}
