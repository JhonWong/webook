package async

import (
	"context"
	"errors"
	"github.com/go-playground/assert/v2"
	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/johnwongx/webook/backend/internal/repository"
	repomocks "github.com/johnwongx/webook/backend/internal/repository/mocks"
	"github.com/johnwongx/webook/backend/internal/service/sms"
	"github.com/johnwongx/webook/backend/internal/service/sms/async/serviceprobe"
	serviceprobemocks "github.com/johnwongx/webook/backend/internal/service/sms/async/serviceprobe/mocks"
	smsmocks "github.com/johnwongx/webook/backend/internal/service/sms/mocks"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestService_Send(t *testing.T) {
	testCases := []struct {
		name string

		svcMock   func(ctrl *gomock.Controller) sms.Service
		probeMock func(ctrl *gomock.Controller) serviceprobe.ServiceProbe
		repoMock  func(ctrl *gomock.Controller) repository.SMSRepository
		inter     time.Duration

		tpl     string
		args    []string
		numbers []string

		wantErr  error
		waitTime time.Duration
	}{
		{
			name: "success",
			svcMock: func(ctrl *gomock.Controller) sms.Service {
				svc := smsmocks.NewMockService(ctrl)
				svc.EXPECT().Send(gomock.Any(), "123", []string{"123"}, []string{"1751234567"}).
					Return(nil)
				return svc
			},
			probeMock: func(ctrl *gomock.Controller) serviceprobe.ServiceProbe {
				p := serviceprobemocks.NewMockServiceProbe(ctrl)
				p.EXPECT().Add(gomock.Any(), nil).Return(true)
				p.EXPECT().IsCrashed(gomock.Any()).Return(false)
				return p
			},
			repoMock: func(ctrl *gomock.Controller) repository.SMSRepository {
				r := repomocks.NewMockSMSRepository(ctrl)
				r.EXPECT().IsEmpty(gomock.Any()).Return(true)
				return r
			},
			inter:   time.Second,
			tpl:     "123",
			args:    []string{"123"},
			numbers: []string{"1751234567"},
		},
		{
			name: "send success, triger resend ",
			svcMock: func(ctrl *gomock.Controller) sms.Service {
				svc := smsmocks.NewMockService(ctrl)
				svc.EXPECT().Send(gomock.Any(), "123", []string{"123"}, []string{"1751234567"}).
					Return(nil)
				return svc
			},
			probeMock: func(ctrl *gomock.Controller) serviceprobe.ServiceProbe {
				p := serviceprobemocks.NewMockServiceProbe(ctrl)
				p.EXPECT().Add(gomock.Any(), nil).Return(true)
				p.EXPECT().IsCrashed(gomock.Any()).Return(false)
				return p
			},
			repoMock: func(ctrl *gomock.Controller) repository.SMSRepository {
				r := repomocks.NewMockSMSRepository(ctrl)
				r.EXPECT().IsEmpty(gomock.Any()).Return(false)
				return r
			},
			inter:   time.Second * 10,
			tpl:     "123",
			args:    []string{"123"},
			numbers: []string{"1751234567"},
		},
		{
			name: "send failed, but not crashed",
			svcMock: func(ctrl *gomock.Controller) sms.Service {
				svc := smsmocks.NewMockService(ctrl)
				svc.EXPECT().Send(gomock.Any(), "123", []string{"123"}, []string{"1751234567"}).
					Return(errors.New("Send failed"))
				return svc
			},
			probeMock: func(ctrl *gomock.Controller) serviceprobe.ServiceProbe {
				p := serviceprobemocks.NewMockServiceProbe(ctrl)
				p.EXPECT().Add(gomock.Any(), errors.New("Send failed")).Return(true)
				p.EXPECT().IsCrashed(gomock.Any()).Return(false)
				return p
			},
			repoMock: func(ctrl *gomock.Controller) repository.SMSRepository {
				return nil
			},
			inter:   time.Second * 10,
			tpl:     "123",
			args:    []string{"123"},
			numbers: []string{"1751234567"},
			wantErr: errors.New("Send failed"),
		},
		{
			name: "send failed, and crashed",
			svcMock: func(ctrl *gomock.Controller) sms.Service {
				svc := smsmocks.NewMockService(ctrl)
				svc.EXPECT().Send(gomock.Any(), "123", []string{"123"}, []string{"1751234567"}).
					Return(errors.New("Send failed"))
				return svc
			},
			probeMock: func(ctrl *gomock.Controller) serviceprobe.ServiceProbe {
				p := serviceprobemocks.NewMockServiceProbe(ctrl)
				p.EXPECT().Add(gomock.Any(), errors.New("Send failed")).Return(true)
				p.EXPECT().IsCrashed(gomock.Any()).Return(true)
				return p
			},
			repoMock: func(ctrl *gomock.Controller) repository.SMSRepository {
				r := repomocks.NewMockSMSRepository(ctrl)
				r.EXPECT().Put(gomock.Any(), domain.SMSInfo{
					Tpl:     "123",
					Args:    []string{"123"},
					Numbers: []string{"1751234567"},
				}).Return(nil)
				return r
			},
			inter:    time.Second * 10,
			tpl:      "123",
			args:     []string{"123"},
			numbers:  []string{"1751234567"},
			wantErr:  errors.New("Send failed"),
			waitTime: time.Second,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			smsSvc := tc.svcMock(ctrl)
			probe := tc.probeMock(ctrl)
			repo := tc.repoMock(ctrl)
			svc := NewService(smsSvc, probe, repo, tc.inter)
			err := svc.Send(context.Background(), tc.tpl, tc.args, tc.numbers...)
			assert.Equal(t, tc.wantErr, err)

			time.Sleep(tc.waitTime)
		})
	}
}
