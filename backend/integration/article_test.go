package integration

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/johnwongx/webook/backend/integration/startup"
	"github.com/johnwongx/webook/backend/internal/repository/dao"
	myjwt "github.com/johnwongx/webook/backend/internal/web/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ArticleHandlerTestSuite struct {
	suite.Suite
	s  *gin.Engine
	db *gorm.DB
}

func (s *ArticleHandlerTestSuite) SetupSubTest() {
	s.s = gin.Default()
	s.db = startup.InitTestDB()
	s.s.Use(func(ctx *gin.Context) {
		ctx.Set("user", myjwt.UserClaim{
			UserId: 123,
		})
		ctx.Next()
	})
	hdl := startup.InitArticleHandler()
	hdl.RegisterRutes(s.s)
}

func (s *ArticleHandlerTestSuite) TearDownSubTest() {
	// TODO
}

func TestArticle(t *testing.T) {
	suite.Run(t, new(ArticleHandlerTestSuite))
}

func (s *ArticleHandlerTestSuite) TestArticleHandler_Edit(t *testing.T) {
	testCases := []struct {
		name     string
		before   func(t *testing.T)
		after    func(t *testing.T)
		Article  Article
		wantCode int
		wantRes  Result[int64]
	}{
		{
			name:   "发帖成功",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				//检查数据库中是否有对应数据
				var art dao.Article
				s.db.Where("where author_id = ?", 123).First(&art)
				assert.True(t, art.CTime > 0)
				assert.True(t, art.UTime > 0)
				art.CTime = 0
				art.UTime = 0
				assert.Equal(t, dao.Article{
					Id:       0,
					Tittle:   "A Tittle",
					Content:  "This is content",
					AuthorId: 123,
				}, art)
			},
			Article: Article{
				Id:      0,
				Tittle:  "A Tittle",
				Content: "This is content",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Data: 1,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)

			data, err := json.Marshal(tc.Article)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/articles/edit", bytes.NewReader(data))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			s.s.ServeHTTP(resp, req)
			assert.Equal(t, tc.wantCode, resp.Code)
			if resp.Code != http.StatusOK {
				return
			}

			var res Result[int64]
			err = json.NewDecoder(resp.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantRes, res)

			tc.after(t)
		})
	}
}

type Article struct {
	Id      int64  `json:"id"`
	Tittle  string `json:"tittle"`
	Content string `json:"content"`
}
