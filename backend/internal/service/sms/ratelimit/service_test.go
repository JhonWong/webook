package ratelimit

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/johnwongx/webook/backend/internal/service/sms"
	smsmocks "github.com/johnwongx/webook/backend/internal/service/sms/mocks"
	"github.com/johnwongx/webook/backend/pkg/ratelimit"
	limitmocks "github.com/johnwongx/webook/backend/pkg/ratelimit/mocks"
	"go.uber.org/mock/gomock"
)

func TestServiceSMSRateLimiter_Send(t *testing.T) {
	testCases := []struct {
		name      string
		limitMock func(ctrl *gomock.Controller) ratelimit.Limiter
		smsMock   func(ctrl *gomock.Controller) sms.Service
		args      []string
		number    []string
		wantErr   error
	}{
		{
			name: "发送成功",
			limitMock: func(ctrl *gomock.Controller) ratelimit.Limiter {
				limiter := limitmocks.NewMockLimiter(ctrl)
				limiter.EXPECT().Limit(gomock.Any(), gomock.Any()).
					Return(false, nil)
				return limiter
			},
			smsMock: func(ctrl *gomock.Controller) sms.Service {
				smsSvc := smsmocks.NewMockService(ctrl)
				smsSvc.EXPECT().Send(gomock.Any(), gomock.Any(), []string{"123456"},
					[]string{"175123456"}).Return(nil)
				return smsSvc
			},
			args:    []string{"123456"},
			number:  []string{"175123456"},
			wantErr: nil,
		},
		{
			name: "触发限流",
			limitMock: func(ctrl *gomock.Controller) ratelimit.Limiter {
				limiter := limitmocks.NewMockLimiter(ctrl)
				limiter.EXPECT().Limit(gomock.Any(), gomock.Any()).
					Return(true, nil)
				return limiter
			},
			smsMock: func(ctrl *gomock.Controller) sms.Service {
				return nil
			},
			args:    []string{"123456"},
			number:  []string{"175123456"},
			wantErr: errLimited,
		},
		{
			name: "限流错误",
			limitMock: func(ctrl *gomock.Controller) ratelimit.Limiter {
				limiter := limitmocks.NewMockLimiter(ctrl)
				limiter.EXPECT().Limit(gomock.Any(), gomock.Any()).
					Return(false, errors.New("系统错误"))
				return limiter
			},
			smsMock: func(ctrl *gomock.Controller) sms.Service {
				return nil
			},
			args:    []string{"123456"},
			number:  []string{"175123456"},
			wantErr: fmt.Errorf("短信服务判断出错，%w", errors.New("系统错误")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			smsSvc := tc.smsMock(ctrl)
			limiter := tc.limitMock(ctrl)
			svc := NewServiceSMSRateLimiter(smsSvc, limiter)
			err := svc.Send(context.Background(), "fake_template",
				tc.args, tc.number...)
			assert.Equal(t, err, tc.wantErr)
		})
	}
}
