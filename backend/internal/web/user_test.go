package web

import (
	"bytes"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/johnwongx/webook/backend/internal/service"
	svcmocks "github.com/johnwongx/webook/backend/internal/service/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_SignUps(t *testing.T) {
	testCase := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserService
		reqBody  string
		wantCode int
		wantBody string
	}{
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) service.UserService {
				us := svcmocks.NewMockUserService(ctrl)
				us.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					PassWord: "hello@world123",
				}).Return(errors.New("邮箱冲突"))
				return us
			},
			reqBody: `
{
	"email":"123@qq.com",
	"passWord":"hello@world123",
	"confirmPassWord":"hello@world123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "邮箱已存在",
		},
		{
			name: "signup success",
			mock: func(ctrl *gomock.Controller) service.UserService {
				us := svcmocks.NewMockUserService(ctrl)
				us.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					PassWord: "hello@world123",
				}).Return(nil)
				return us
			},
			reqBody: `
{
	"email":"123@qq.com",
	"passWord":"hello@world123",
	"confirmPassWord":"hello@world123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "注册成功",
		},
		{
			name: "参数错误",
			mock: func(ctrl *gomock.Controller) service.UserService {
				us := svcmocks.NewMockUserService(ctrl)
				return us
			},
			reqBody: `
{
	"email":"123@qq.com",
	"passWord":"hello@world123",
}
`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "邮箱格式错误",
			mock: func(ctrl *gomock.Controller) service.UserService {
				us := svcmocks.NewMockUserService(ctrl)
				return us
			},
			reqBody: `
{
	"email":"123@qqcom",
	"passWord":"hello@world123",
	"confirmPassWord":"hello@world123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "邮箱格式错误",
		},
		{
			name: "密码不一致",
			mock: func(ctrl *gomock.Controller) service.UserService {
				us := svcmocks.NewMockUserService(ctrl)
				return us
			},
			reqBody: `
{
	"email":"123@qq.com",
	"passWord":"hello@world123",
	"confirmPassWord":"Hello@world123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "两次输入密码不一致",
		},
		{
			name: "密码格式错误",
			mock: func(ctrl *gomock.Controller) service.UserService {
				us := svcmocks.NewMockUserService(ctrl)
				return us
			},
			reqBody: `
{
	"email":"123@qq.com",
	"passWord":"helloworld123",
	"confirmPassWord":"helloworld123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "密码格式错误",
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) service.UserService {
				us := svcmocks.NewMockUserService(ctrl)
				us.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					PassWord: "hello@world123",
				}).Return(errors.New("系统错误"))
				return us
			},
			reqBody: `
{
	"email":"123@qq.com",
	"passWord":"hello@world123",
	"confirmPassWord":"hello@world123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "系统错误",
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			server := gin.Default()
			h := NewUserHandler(tc.mock(ctrl), nil)
			h.RegisterRoutes(server)

			req, err := http.NewRequest(http.MethodPost,
				"/users/signup", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)

			assert.Equal(t, resp.Code, tc.wantCode)
			assert.Equal(t, resp.Body.String(), tc.wantBody)
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
	}).Return(errors.New("fxxk u"))

	err := us.SignUp(context.Background(), domain.User{
		Email:    "123@qq.com",
		PassWord: "hello#world123",
	})
	t.Log(err)
}
