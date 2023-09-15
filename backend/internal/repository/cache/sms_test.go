package cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-playground/assert/v2"
	"github.com/johnwongx/webook/backend/internal/repository/cache/redismocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestSmsCache_Add(t *testing.T) {
	testCase := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) redis.Cmdable
		info    SMSInfo
		wantErr error
	}{
		{
			name: "success",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				rdb := redismocks.NewMockCmdable(ctrl)
				info := SMSInfo{
					Tpl:     "123",
					Args:    []string{"123"},
					Numbers: []string{"175123456"},
				}
				jsonData, err := json.Marshal(info)
				require.NoError(t, err)

				res := redis.NewIntCmd(context.Background())
				res.SetVal(1)
				res.SetErr(nil)
				rdb.EXPECT().LPush(gomock.Any(), sms_key, jsonData).Return(res)
				return rdb
			},
			info: SMSInfo{
				Tpl:     "123",
				Args:    []string{"123"},
				Numbers: []string{"175123456"},
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rdb := tc.mock(ctrl)
			svc := NewSMSCache(rdb)
			err := svc.Add(context.Background(), tc.info)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestSmsCache_Take(t *testing.T) {
	testCase := []struct {
		name      string
		mock      func(ctrl *gomock.Controller) redis.Cmdable
		cnt       int
		wantInfos []SMSInfo
		wantErr   error
	}{
		{
			name: "take 1",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				info := SMSInfo{
					Tpl:     "123",
					Args:    []string{"123"},
					Numbers: []string{"175123456"},
				}
				jsonData, err := json.Marshal(info)
				require.NoError(t, err)

				res := redis.NewStringCmd(context.Background())
				res.SetVal(string(jsonData))
				res.SetErr(nil)

				rdb := redismocks.NewMockCmdable(ctrl)
				rdb.EXPECT().RPop(gomock.Any(), sms_key).Return(res)
				return rdb
			},
			cnt: 1,
			wantInfos: []SMSInfo{{
				Tpl:     "123",
				Args:    []string{"123"},
				Numbers: []string{"175123456"},
			}},
		},
		{
			name: "take 2",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				infoArr := []SMSInfo{
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

				jsonData1, err := json.Marshal(infoArr[0])
				require.NoError(t, err)
				res1 := redis.NewStringCmd(context.Background())
				res1.SetVal(string(jsonData1))
				res1.SetErr(nil)
				rdb := redismocks.NewMockCmdable(ctrl)
				rdb.EXPECT().RPop(gomock.Any(), sms_key).Return(res1)

				jsonData2, err := json.Marshal(infoArr[1])
				require.NoError(t, err)
				res2 := redis.NewStringCmd(context.Background())
				res2.SetVal(string(jsonData2))
				res2.SetErr(nil)
				rdb.EXPECT().RPop(gomock.Any(), sms_key).Return(res2)

				return rdb
			},
			cnt: 2,
			wantInfos: []SMSInfo{
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
			name: "take 2 return 1",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				info := SMSInfo{
					Tpl:     "123",
					Args:    []string{"123"},
					Numbers: []string{"175123456"},
				}
				jsonData, err := json.Marshal(info)
				require.NoError(t, err)

				rdb := redismocks.NewMockCmdable(ctrl)

				res := redis.NewStringCmd(context.Background())
				res.SetVal(string(jsonData))
				res.SetErr(nil)
				rdb.EXPECT().RPop(gomock.Any(), sms_key).Return(res)

				res1 := redis.NewStringCmd(context.Background())
				res1.SetErr(redis.Nil)
				rdb.EXPECT().RPop(gomock.Any(), sms_key).Return(res1)
				return rdb
			},
			cnt: 2,
			wantInfos: []SMSInfo{
				{
					Tpl:     "123",
					Args:    []string{"123"},
					Numbers: []string{"175123456"},
				},
			},
		},
		{
			name: "redis error",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				res := redis.NewStringCmd(context.Background())
				res.SetErr(errors.New("redis error"))
				rdb := redismocks.NewMockCmdable(ctrl)
				rdb.EXPECT().RPop(gomock.Any(), sms_key).Return(res)
				return rdb
			},
			cnt:     2,
			wantErr: errors.New("redis error"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rdb := tc.mock(ctrl)
			svc := NewSMSCache(rdb)

			infos, err := svc.Take(context.Background(), tc.cnt)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantInfos, infos)
		})
	}
}

func TestSmsCache_KeyExists(t *testing.T) {
	testCase := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) redis.Cmdable
		wantRet bool
		wantErr error
	}{
		{
			name: "not empty",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				res := redis.NewIntCmd(context.Background())
				res.SetVal(1)
				res.SetErr(nil)

				rdb := redismocks.NewMockCmdable(ctrl)
				rdb.EXPECT().Exists(gomock.Any(), sms_key).Return(res)
				return rdb
			},
			wantRet: true,
		},
		{
			name: "empty",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				res := redis.NewIntCmd(context.Background())
				res.SetVal(0)
				res.SetErr(nil)

				rdb := redismocks.NewMockCmdable(ctrl)
				rdb.EXPECT().Exists(gomock.Any(), sms_key).Return(res)
				return rdb
			},
		},
		{
			name: "redis error",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				res := redis.NewIntCmd(context.Background())
				res.SetErr(errors.New("redis error"))

				rdb := redismocks.NewMockCmdable(ctrl)
				rdb.EXPECT().Exists(gomock.Any(), sms_key).Return(res)
				return rdb
			},
			wantErr: errors.New("redis error"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rdb := tc.mock(ctrl)
			svc := NewSMSCache(rdb)

			res, err := svc.KeyExists(context.Background())
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantRet, res)
		})
	}
}
