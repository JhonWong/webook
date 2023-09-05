package cache

import (
	"context"
	"errors"
	"github.com/go-playground/assert/v2"
	"github.com/johnwongx/webook/backend/internal/repository/cache/redismocks"
	"github.com/redis/go-redis/v9"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestRedisCodeCache_Set(t *testing.T) {
	testCase := []struct {
		name       string
		mock       func(ctrl *gomock.Controller) redis.Cmdable
		biz        string
		phone      string
		code       string
		experation time.Duration
		wantErr    error
	}{
		{
			name: "success",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				r := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(0))
				r.EXPECT().Eval(gomock.Any(), luaSetCode,
					[]string{"phone_code:login:152"},
					[]any{"123456"},
				).Return(res)
				return r
			},
			biz:     "login",
			phone:   "152",
			code:    "123456",
			wantErr: nil,
		},
		{
			name: "redis error",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				r := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetErr(errors.New("redis error"))
				r.EXPECT().Eval(gomock.Any(), luaSetCode,
					[]string{"phone_code:login:152"},
					[]any{"123456"},
				).Return(res)
				return r
			},
			biz:     "login",
			phone:   "152",
			code:    "123456",
			wantErr: errors.New("redis error"),
		},
		{
			name: "send too many",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				r := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(-1))
				r.EXPECT().Eval(gomock.Any(), luaSetCode,
					[]string{"phone_code:login:152"},
					[]any{"123456"},
				).Return(res)
				return r
			},
			biz:     "login",
			phone:   "152",
			code:    "123456",
			wantErr: ErrCodeSendTooMany,
		},
		{
			name: "other error",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				r := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(-2))
				r.EXPECT().Eval(gomock.Any(), luaSetCode,
					[]string{"phone_code:login:152"},
					[]any{"123456"},
				).Return(res)
				return r
			},
			biz:     "login",
			phone:   "152",
			code:    "123456",
			wantErr: errors.New("系统错误"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			r := tc.mock(ctrl)
			rc := NewRedisCodeCache(r)
			err := rc.Set(context.Background(), tc.biz, tc.phone, tc.code, tc.experation)
			assert.Equal(t, err, tc.wantErr)
		})
	}
}
