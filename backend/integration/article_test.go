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

func (s *ArticleHandlerTestSuite) SetupSuite() {
	s.s = gin.Default()
	s.db = startup.InitTestDB()
	s.s.Use(func(ctx *gin.Context) {
		ctx.Set("claims", myjwt.UserClaim{
			UserId: 123,
		})
		ctx.Next()
	})
	hdl := startup.InitArticleHandler()
	hdl.RegisterRutes(s.s)
}

func (s *ArticleHandlerTestSuite) TearDownTest() {
	err := s.db.Exec("TRUNCATE TABLE `articles`").Error
	assert.NoError(s.T(), err)
}

func TestArticle(t *testing.T) {
	suite.Run(t, new(ArticleHandlerTestSuite))
}

func (s *ArticleHandlerTestSuite) TestArticleHandler_Edit() {
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
				s.db.Where("author_id = ?", 123).First(&art)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Utime > 0)
				art.Ctime = 0
				art.Utime = 0
				assert.Equal(t, dao.Article{
					Id:       1,
					Tittle:   "A Tittle",
					Content:  "This is content",
					AuthorId: 123,
				}, art)
			},
			Article: Article{
				Tittle:  "A Tittle",
				Content: "This is content",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Data: 1,
			},
		},
		{
			name: "编辑成功",
			before: func(t *testing.T) {
				s.db.Create(dao.Article{
					Id:       2,
					Tittle:   "My tittle",
					Content:  "My Content",
					AuthorId: 123,
					Ctime:    123,
					Utime:    678,
				})
			},
			after: func(t *testing.T) {
				//检查数据库中是否有对应数据
				var art dao.Article
				s.db.Where("id = ?", 2).First(&art)
				assert.True(t, art.Utime > 678)
				art.Utime = 0
				assert.Equal(t, dao.Article{
					Id:       2,
					Tittle:   "New Tittle",
					Content:  "new content",
					AuthorId: 123,
					Ctime:    123,
				}, art)
			},
			Article: Article{
				Id:      2,
				Tittle:  "New Tittle",
				Content: "new content",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Data: 2,
			},
		},
		{
			name: "修改别人帖子",
			before: func(t *testing.T) {
				s.db.Create(dao.Article{
					Id:       3,
					Tittle:   "My tittle",
					Content:  "My Content",
					AuthorId: 233,
					Ctime:    123,
					Utime:    678,
				})
			},
			after: func(t *testing.T) {
				//检查数据库中是否有对应数据
				var art dao.Article
				s.db.Where("id = ?", 3).First(&art)
				assert.Equal(t, dao.Article{
					Id:       3,
					Tittle:   "My tittle",
					Content:  "My Content",
					AuthorId: 233,
					Ctime:    123,
					Utime:    678,
				}, art)
			},
			Article: Article{
				Id:      3,
				Tittle:  "New Tittle",
				Content: "new content",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Code: 5,
				Msg:  "系统错误",
			},
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
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
