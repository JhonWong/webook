package integration

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/johnwongx/webook/backend/integration/startup"
	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/johnwongx/webook/backend/internal/repository/dao/article"
	myjwt "github.com/johnwongx/webook/backend/internal/web/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ArticleGORMHandlerTestSuite struct {
	suite.Suite
	s  *gin.Engine
	db *gorm.DB
}

func (s *ArticleGORMHandlerTestSuite) SetupSuite() {
	s.s = gin.Default()
	s.db = startup.InitTestDB()
	s.s.Use(func(ctx *gin.Context) {
		ctx.Set("claims", myjwt.UserClaim{
			UserId: 123,
		})
		ctx.Next()
	})
	hdl := startup.InitArticleHandler(article.NewGORMArticleDAO(s.db))
	hdl.RegisterRutes(s.s)
}

func (s *ArticleGORMHandlerTestSuite) TearDownTest() {
	err := s.db.Exec("TRUNCATE TABLE `articles`").Error
	assert.NoError(s.T(), err)
	err = s.db.Exec("TRUNCATE TABLE `publish_articles`").Error
	assert.NoError(s.T(), err)
}

func TestArticleGORM(t *testing.T) {
	suite.Run(t, new(ArticleGORMHandlerTestSuite))
}

func (s *ArticleGORMHandlerTestSuite) TestArticleHandler_Withdraw() {
	testCases := []struct {
		name     string
		before   func(t *testing.T)
		after    func(t *testing.T)
		req      string
		wantCode int
		wantRes  Result[int64]
	}{
		{
			name: "修改自己帖子",
			before: func(t *testing.T) {
				art := article.Article{
					Id:       3,
					Tittle:   "My tittle",
					Content:  "My Content",
					AuthorId: 123,
					Ctime:    123,
					Utime:    678,
					Status:   domain.ArticleStatusPublished.ToUint8(),
				}
				err := s.db.Create(art).Error
				assert.NoError(t, err)
				err = s.db.Create(article.PublishArticle{Article: art}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				//检查数据库中是否有对应数据
				var eArt article.Article
				err := s.db.Where("id = ?", 3).First(&eArt).Error
				assert.NoError(t, err)
				assert.True(t, eArt.Utime > 678)
				eArt.Utime = 0
				assert.Equal(t, article.Article{
					Id:       3,
					Tittle:   "My tittle",
					Content:  "My Content",
					AuthorId: 123,
					Ctime:    123,
					Status:   domain.ArticleStatusPrivate.ToUint8(),
				}, eArt)

				var art article.PublishArticle
				err = s.db.Where("id = ?", 3).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.Utime > 678)
				art.Utime = 0
				assert.Equal(t, article.PublishArticle{Article: article.Article{
					Id:       3,
					Tittle:   "My tittle",
					Content:  "My Content",
					AuthorId: 123,
					Ctime:    123,
					Status:   domain.ArticleStatusPrivate.ToUint8(),
				},
				}, art)
			},
			req: `
{
	"id":3
}
`,
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Data: 3,
			},
		},
		{
			name: "修改别人帖子，并失败",
			before: func(t *testing.T) {
				art := article.Article{
					Id:       4,
					Tittle:   "My tittle",
					Content:  "My Content",
					AuthorId: 233,
					Ctime:    123,
					Utime:    678,
					Status:   domain.ArticleStatusPublished.ToUint8(),
				}
				err := s.db.Create(art).Error
				assert.NoError(t, err)
				err = s.db.Create(article.PublishArticle{Article: art}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				//检查数据库中是否有对应数据
				var eArt article.Article
				err := s.db.Where("id = ?", 4).First(&eArt).Error
				assert.NoError(t, err)
				assert.Equal(t, article.Article{
					Id:       4,
					Tittle:   "My tittle",
					Content:  "My Content",
					AuthorId: 233,
					Ctime:    123,
					Utime:    678,
					Status:   domain.ArticleStatusPublished.ToUint8(),
				}, eArt)

				var art article.PublishArticle
				err = s.db.Where("id = ?", 4).First(&art).Error
				assert.NoError(t, err)
				assert.Equal(t, article.PublishArticle{Article: article.Article{
					Id:       4,
					Tittle:   "My tittle",
					Content:  "My Content",
					AuthorId: 233,
					Ctime:    123,
					Utime:    678,
					Status:   domain.ArticleStatusPublished.ToUint8(),
				},
				}, art)
			},
			req: `
{
	"id":4
}`,
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Code: 5,
				Msg:  "系统错误",
			},
		},
		{
			name:   "数据格式错误",
			before: func(t *testing.T) {},
			after:  func(t *testing.T) {},
			req: `
{
	"id":
}`,
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			tc.before(t)

			req, err := http.NewRequest(http.MethodPost, "/articles/withdraw", bytes.NewReader([]byte(tc.req)))
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

func (s *ArticleGORMHandlerTestSuite) TestArticleHandler_Publish() {
	testCases := []struct {
		name     string
		before   func(t *testing.T)
		after    func(t *testing.T)
		Article  Article
		wantCode int
		wantRes  Result[int64]
	}{
		{
			name:   "新建帖子，直接发表成功",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				//检查数据库中是否有对应数据
				var art article.PublishArticle
				s.db.Where("author_id = ?", 123).First(&art)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Utime > 0)
				art.Ctime = 0
				art.Utime = 0
				assert.Equal(t, article.PublishArticle{Article: article.Article{
					Id:       1,
					Tittle:   "A Tittle",
					Content:  "This is content",
					AuthorId: 123,
					Status:   domain.ArticleStatusPublished.ToUint8(),
				},
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
			name: "更新帖子，线上不存在",
			before: func(t *testing.T) {
				art := article.Article{
					Id:       2,
					Tittle:   "My tittle",
					Content:  "My Content",
					AuthorId: 123,
					Ctime:    123,
					Utime:    678,
					Status:   domain.ArticleStatusUnpublished.ToUint8(),
				}
				err := s.db.Create(art).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				//检查数据库中是否有对应数据
				var eArt article.Article
				err := s.db.Where("id = ?", 2).First(&eArt).Error
				assert.NoError(t, err)
				assert.True(t, eArt.Utime > 678)
				eArt.Utime = 0
				assert.Equal(t, article.Article{
					Id:       2,
					Tittle:   "New Tittle",
					Content:  "new content",
					AuthorId: 123,
					Ctime:    123,
					Status:   domain.ArticleStatusPublished.ToUint8(),
				}, eArt)

				var art article.PublishArticle
				err = s.db.Where("id = ?", 2).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.Utime > 678)
				assert.True(t, art.Ctime > 123)
				art.Utime = 0
				art.Ctime = 0
				assert.Equal(t, article.PublishArticle{Article: article.Article{
					Id:       2,
					Tittle:   "New Tittle",
					Content:  "new content",
					AuthorId: 123,
					Status:   domain.ArticleStatusPublished.ToUint8(),
				},
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
			name: "更新帖子，线上已存在",
			before: func(t *testing.T) {
				art := article.Article{
					Id:       3,
					Tittle:   "My tittle",
					Content:  "My Content",
					AuthorId: 123,
					Ctime:    123,
					Utime:    678,
					Status:   domain.ArticleStatusPublished.ToUint8(),
				}
				err := s.db.Create(art).Error
				assert.NoError(t, err)
				err = s.db.Create(article.PublishArticle{Article: art}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				//检查数据库中是否有对应数据
				var eArt article.Article
				err := s.db.Where("id = ?", 3).First(&eArt).Error
				assert.NoError(t, err)
				assert.True(t, eArt.Utime > 678)
				eArt.Utime = 0
				assert.Equal(t, article.Article{
					Id:       3,
					Tittle:   "New Tittle",
					Content:  "new content",
					AuthorId: 123,
					Ctime:    123,
					Status:   domain.ArticleStatusPublished.ToUint8(),
				}, eArt)

				var art article.PublishArticle
				err = s.db.Where("id = ?", 3).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.Utime > 678)
				art.Utime = 0
				assert.Equal(t, article.PublishArticle{Article: article.Article{
					Id:       3,
					Tittle:   "New Tittle",
					Content:  "new content",
					AuthorId: 123,
					Ctime:    123,
					Status:   domain.ArticleStatusPublished.ToUint8(),
				},
				}, art)
			},
			Article: Article{
				Id:      3,
				Tittle:  "New Tittle",
				Content: "new content",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Data: 3,
			},
		},
		{
			name: "修改别人帖子，并失败",
			before: func(t *testing.T) {
				art := article.Article{
					Id:       4,
					Tittle:   "My tittle",
					Content:  "My Content",
					AuthorId: 233,
					Ctime:    123,
					Utime:    678,
					Status:   domain.ArticleStatusPublished.ToUint8(),
				}
				err := s.db.Create(art).Error
				assert.NoError(t, err)
				err = s.db.Create(article.PublishArticle{Article: art}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				//检查数据库中是否有对应数据
				var eArt article.Article
				err := s.db.Where("id = ?", 4).First(&eArt).Error
				assert.NoError(t, err)
				assert.Equal(t, article.Article{
					Id:       4,
					Tittle:   "My tittle",
					Content:  "My Content",
					AuthorId: 233,
					Ctime:    123,
					Utime:    678,
					Status:   domain.ArticleStatusPublished.ToUint8(),
				}, eArt)

				var art article.PublishArticle
				err = s.db.Where("id = ?", 4).First(&art).Error
				assert.NoError(t, err)
				assert.Equal(t, article.PublishArticle{Article: article.Article{
					Id:       4,
					Tittle:   "My tittle",
					Content:  "My Content",
					AuthorId: 233,
					Ctime:    123,
					Utime:    678,
					Status:   domain.ArticleStatusPublished.ToUint8(),
				},
				}, art)
			},
			Article: Article{
				Id:      4,
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
			req, err := http.NewRequest(http.MethodPost, "/articles/publish", bytes.NewReader(data))
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

func (s *ArticleGORMHandlerTestSuite) TestArticleHandler_Edit() {
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
				var art article.Article
				s.db.Where("author_id = ?", 123).First(&art)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Utime > 0)
				art.Ctime = 0
				art.Utime = 0
				assert.Equal(t, article.Article{
					Id:       1,
					Tittle:   "A Tittle",
					Content:  "This is content",
					AuthorId: 123,
					Status:   domain.ArticleStatusUnpublished.ToUint8(),
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
				err := s.db.Create(article.Article{
					Id:       2,
					Tittle:   "My tittle",
					Content:  "My Content",
					AuthorId: 123,
					Ctime:    123,
					Utime:    678,
					Status:   domain.ArticleStatusPublished.ToUint8(),
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				//检查数据库中是否有对应数据
				var art article.Article
				err := s.db.Where("id = ?", 2).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.Utime > 678)
				art.Utime = 0
				assert.Equal(t, article.Article{
					Id:       2,
					Tittle:   "New Tittle",
					Content:  "new content",
					AuthorId: 123,
					Ctime:    123,
					Status:   domain.ArticleStatusUnpublished.ToUint8(),
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
				s.db.Create(article.Article{
					Id:       3,
					Tittle:   "My tittle",
					Content:  "My Content",
					AuthorId: 233,
					Ctime:    123,
					Utime:    678,
					Status:   domain.ArticleStatusPublished.ToUint8(),
				})
			},
			after: func(t *testing.T) {
				//检查数据库中是否有对应数据
				var art article.Article
				s.db.Where("id = ?", 3).First(&art)
				assert.Equal(t, article.Article{
					Id:       3,
					Tittle:   "My tittle",
					Content:  "My Content",
					AuthorId: 233,
					Ctime:    123,
					Utime:    678,
					Status:   domain.ArticleStatusPublished.ToUint8(),
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
