package web

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/johnwongx/webook/backend/internal/domain"
	svcmocks "github.com/johnwongx/webook/backend/internal/service/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_SignUps(t *testing.T) {
	testCase := []struct {
		name string
	}{
		{},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			server := gin.Default()
			h := NewUserHandler(nil, nil)
			h.RegisterRoutes(server)

			req, err := http.NewRequest(http.MethodPost,
				"users/signup", bytes.NewBuffer([]byte("")))
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)

			assert.Equal(t, resp.Code, nil)
			assert.Equal(t, resp.Body.String(), nil)
		})
	}
}

func TestMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	us := svcmocks.NewMockUserService(ctrl)
	us.EXPECT().SignUp(gomock.Any(), domain.User{
		Email:    "123@qq.com",
		PassWord: "hello#world123",
	}).Return(nil)

	err := us.SignUp(context.Background(), domain.User{
		Email:    "123@qq.com",
		PassWord: "hello#world123",
	})
	t.Log(err)
}
