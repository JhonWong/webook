package service

import (
	"context"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/johnwongx/webook/backend/internal/repository"
	repomocks "github.com/johnwongx/webook/backend/internal/repository/mocks"
	"github.com/johnwongx/webook/backend/internal/service/sms/localsms"
	"go.uber.org/mock/gomock"
)

func TestCodeService_Send(t *testing.T) {
	testCase := []struct {
		name     string
		biz      string
		phone    string
		repoFunc func(controller *gomock.Controller) *repomocks.MockCodeRepository
		wantErr  error
	}{
		{
			name:  "success",
			biz:   "login",
			phone: "10086",
			repoFunc: func(ctrl *gomock.Controller) *repomocks.MockCodeRepository {
				repo := repomocks.NewMockCodeRepository(ctrl)
				repo.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any()).Return(nil)
				return repo
			},
			wantErr: nil,
		},
		{
			name:  "repository sotre error",
			biz:   "login",
			phone: "10086",
			repoFunc: func(ctrl *gomock.Controller) *repomocks.MockCodeRepository {
				repo := repomocks.NewMockCodeRepository(ctrl)
				repo.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
					gomock.Any()).Return(repository.ErrCodeSendTooMany)
				return repo
			},
			wantErr: repository.ErrCodeSendTooMany,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			smsSvc := localsms.NewService()
			repo := tc.repoFunc(ctrl)
			cs := NewCodeService(smsSvc, repo)
			err := cs.Send(context.Background(), tc.biz, tc.phone)
			assert.Equal(t, err, tc.wantErr)
		})
	}
}
