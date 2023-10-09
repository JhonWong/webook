package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/johnwongx/webook/backend/internal/domain"
	svcmocks "github.com/johnwongx/webook/backend/internal/service/mocks"
	myjwt "github.com/johnwongx/webook/backend/internal/web/jwt"
	"github.com/johnwongx/webook/backend/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestArticleHandler_Publish(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) *svcmocks.MockArticleService
		reqBody  string
		wantCode int
		wantMsg  Result
	}{
		{
			name: "发表成功",
			mock: func(ctrl *gomock.Controller) *svcmocks.MockArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Tittle:  "my tittle",
					Content: "my content",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				return svc
			},
			reqBody: `
{
	"tittle":"my tittle",
	"content":"my content"
}
`,
			wantCode: http.StatusOK,
			wantMsg: Result{
				Data: float64(1),
			},
		},
		{
			name: "已有帖子发表成功",
			mock: func(ctrl *gomock.Controller) *svcmocks.MockArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Id:      2,
					Tittle:  "my tittle",
					Content: "my content",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(2), nil)
				return svc
			},
			reqBody: `
{
	"id":2,
	"tittle":"my tittle",
	"content":"my content"
}
`,
			wantCode: http.StatusOK,
			wantMsg: Result{
				Data: float64(2),
			},
		},
		{
			name: "发表帖子失败",
			mock: func(ctrl *gomock.Controller) *svcmocks.MockArticleService {
				svc := svcmocks.NewMockArticleService(ctrl)
				svc.EXPECT().Publish(gomock.Any(), domain.Article{
					Id:      2,
					Tittle:  "my tittle",
					Content: "my content",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(2), errors.New("system error"))
				return svc
			},
			reqBody: `
{
	"id":2,
	"tittle":"my tittle",
	"content":"my content"
}
`,
			wantCode: http.StatusOK,
			wantMsg: Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
		{
			name: "bind错误",
			mock: func(ctrl *gomock.Controller) *svcmocks.MockArticleService {
				return nil
			},
			reqBody: `
{
	"id":2
	"tittle":"my tittle",
	"content":"my content"
}
`,
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			svc := tc.mock(ctrl)
			hdl := NewArticleHandler(svc, &logger.NopLogger{})

			server := gin.Default()
			server.Use(func(ctx *gin.Context) {
				ctx.Set("claims", myjwt.UserClaim{
					UserId: 123,
				})
				ctx.Next()
			})
			hdl.RegisterRutes(server)

			req, err := http.NewRequest(http.MethodPost, "/articles/publish",
				bytes.NewBufferString(tc.reqBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)

			assert.Equal(t, resp.Code, tc.wantCode)
			if resp.Code != http.StatusOK {
				return
			}

			var res Result
			err = json.Unmarshal(resp.Body.Bytes(), &res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantMsg, res)
		})
	}
}
