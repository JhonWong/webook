package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/johnwongx/webook/backend/integration/startup"
	"github.com/johnwongx/webook/backend/internal/web"
	"github.com/johnwongx/webook/backend/ioc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUserHandler_e2e_SendLoginSMSCode(t *testing.T) {
	server := startup.InitWebServer()
	rdb := ioc.InitRedis()
	testCase := []struct {
		name     string
		before   func(t *testing.T)
		after    func(t *testing.T)
		reqBody  string
		wantCode int
		wantBody web.Result
	}{
		{
			name: "success",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				val, err := rdb.GetDel(ctx, "phone_code:login:15212345678").Result()
				cancel()
				assert.NoError(t, err)
				assert.True(t, len(val) == 6)
			},
			reqBody: `
{
	"phoneNumber": "15212345678"
}
`,
			wantCode: 200,
			wantBody: web.Result{
				Msg: "发送成功",
			},
		},
		{
			name: "send too many",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				_, err := rdb.Set(ctx, "phone_code:login:15212345678", "123456",
					time.Minute*9+time.Second*30).Result()
				cancel()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				val, err := rdb.GetDel(ctx, "phone_code:login:15212345678").Result()
				cancel()
				assert.NoError(t, err)
				assert.True(t, len(val) == 6)
			},
			reqBody: `
{
	"phoneNumber": "15212345678"
}
`,
			wantCode: 200,
			wantBody: web.Result{
				Msg: "短信发送太频繁，请稍后再试",
			},
		},
		{
			name: "code service error",
			before: func(t *testing.T) {
				//该手机已存在验证码但无过期时间
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				_, err := rdb.Set(ctx, "phone_code:login:15212345678", "123456", 0).Result()
				cancel()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				val, err := rdb.GetDel(ctx, "phone_code:login:15212345678").Result()
				cancel()
				assert.NoError(t, err)
				assert.True(t, len(val) == 6)
			},
			reqBody: `
{
	"phoneNumber": "15212345678"
}
`,
			wantCode: 200,
			wantBody: web.Result{
				Msg:  "系统错误",
				Code: 5,
			},
		},
		{
			name: "empty phone number",
			before: func(t *testing.T) {
			},
			after: func(t *testing.T) {
			},
			reqBody: `
{
	"phoneNumber": ""
}
`,
			wantCode: 200,
			wantBody: web.Result{
				Msg:  "请输入手机号码",
				Code: 4,
			},
		},
		{
			name: "request data format error",
			before: func(t *testing.T) {
			},
			after: func(t *testing.T) {
			},
			reqBody: `
{
	"phone", "15212345678"
}
`,
			wantCode: 400,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			req, err := http.NewRequest(http.MethodPost,
				"/users/login_sms/code/send", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)
			assert.Equal(t, tc.wantCode, resp.Code)
			if resp.Code != 200 {
				return
			}

			var webRes web.Result
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			require.NoError(t, err)
			assert.Equal(t, tc.wantBody, webRes)
			tc.after(t)
		})
	}
}
