package failover

import (
	"context"
	"errors"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/johnwongx/webook/backend/internal/service/sms"
	smsmocks "github.com/johnwongx/webook/backend/internal/service/sms/mocks"
	"go.uber.org/mock/gomock"
)

func TestFailover_Send(t *testing.T) {
	testCases := []struct {
		name string
		mock func(*gomock.Controller) []sms.Service

		idx       uint32
		cnt       uint32
		threshold uint32

		wantErr error
		wantIdx uint32
		wantCnt uint32
	}{
		{
			name: "超时，未连续超时",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any()).Return(context.DeadlineExceeded)

				return []sms.Service{svc1}
			},
			threshold: 3,
			wantErr:   context.DeadlineExceeded,
			wantIdx:   uint32(0),
			wantCnt:   uint32(1),
		},
		{
			name: "超时，切换后成功",
			cnt:  uint32(4),
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc1 := smsmocks.NewMockService(ctrl)
				svc2 := smsmocks.NewMockService(ctrl)
				svc2.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any()).Return(nil)

				return []sms.Service{svc1, svc2}
			},
			threshold: 3,
			wantErr:   nil,
			wantIdx:   uint32(1),
			wantCnt:   uint32(0),
		},
		{
			name: "切换后发生错误",
			cnt:  uint32(4),
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc1 := smsmocks.NewMockService(ctrl)
				svc2 := smsmocks.NewMockService(ctrl)
				svc2.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any()).Return(errors.New("发送错误"))

				return []sms.Service{svc1, svc2}
			},
			threshold: 3,
			wantErr:   errors.New("发送错误"),
			wantIdx:   uint32(1),
			wantCnt:   uint32(0),
		},
		{
			name: "切换后超时",
			cnt:  uint32(4),
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc1 := smsmocks.NewMockService(ctrl)
				svc2 := smsmocks.NewMockService(ctrl)
				svc2.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any()).Return(context.DeadlineExceeded)

				return []sms.Service{svc1, svc2}
			},
			threshold: 3,
			wantErr:   context.DeadlineExceeded,
			wantIdx:   uint32(1),
			wantCnt:   uint32(1),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svcs := tc.mock(ctrl)
			fs := NewTimeoutFailoverSMSService(svcs, tc.threshold)
			fs.idx = tc.idx
			fs.cnt = tc.cnt
			err := fs.Send(context.Background(), "fake_template", []string{"123"},
				"123")

			assert.Equal(t, err, tc.wantErr)
			assert.Equal(t, fs.idx, tc.wantIdx)
			assert.Equal(t, fs.cnt, tc.wantCnt)
		})
	}
}
