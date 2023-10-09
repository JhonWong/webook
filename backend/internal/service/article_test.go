package service

import (
	"context"
	"errors"
	"github.com/go-playground/assert/v2"
	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/johnwongx/webook/backend/internal/repository"
	repomocks "github.com/johnwongx/webook/backend/internal/repository/mocks"
	"github.com/johnwongx/webook/backend/pkg/logger"
	"go.uber.org/mock/gomock"
	"testing"
)

func Test_articleService_Publish(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) (repository.AuthorArticleRepository, repository.ReaderArticleRepository)
		art     domain.Article
		wantId  int64
		wantErr error
	}{
		{
			name: "直接发表成功",
			mock: func(ctrl *gomock.Controller) (repository.AuthorArticleRepository, repository.ReaderArticleRepository) {
				authorRepo := repomocks.NewMockAuthorArticleRepository(ctrl)
				authorRepo.EXPECT().Create(gomock.Any(), domain.Article{
					Tittle:  "tittle",
					Content: "content",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)

				readerRepo := repomocks.NewMockReaderArticleRepository(ctrl)
				readerRepo.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      1,
					Tittle:  "tittle",
					Content: "content",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)
				return authorRepo, readerRepo
			},
			art: domain.Article{
				Tittle:  "tittle",
				Content: "content",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId:  1,
			wantErr: nil,
		},
		{
			name: "草稿箱已存在，未发表",
			mock: func(ctrl *gomock.Controller) (repository.AuthorArticleRepository, repository.ReaderArticleRepository) {
				authorRepo := repomocks.NewMockAuthorArticleRepository(ctrl)
				authorRepo.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      2,
					Tittle:  "tittle",
					Content: "content",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)

				readerRepo := repomocks.NewMockReaderArticleRepository(ctrl)
				readerRepo.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      2,
					Tittle:  "tittle",
					Content: "content",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)
				return authorRepo, readerRepo
			},
			art: domain.Article{
				Id:      2,
				Tittle:  "tittle",
				Content: "content",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId:  2,
			wantErr: nil,
		},
		{
			name: "草稿箱更新失败",
			mock: func(ctrl *gomock.Controller) (repository.AuthorArticleRepository, repository.ReaderArticleRepository) {
				authorRepo := repomocks.NewMockAuthorArticleRepository(ctrl)
				authorRepo.EXPECT().Create(gomock.Any(), domain.Article{
					Tittle:  "tittle",
					Content: "content",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(0), errors.New("update failed!"))

				return authorRepo, nil
			},
			art: domain.Article{
				Tittle:  "tittle",
				Content: "content",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantErr: errors.New("update failed!"),
		},
		{
			name: "线上库重试后成功",
			mock: func(ctrl *gomock.Controller) (repository.AuthorArticleRepository, repository.ReaderArticleRepository) {
				authorRepo := repomocks.NewMockAuthorArticleRepository(ctrl)
				authorRepo.EXPECT().Create(gomock.Any(), domain.Article{
					Tittle:  "tittle",
					Content: "content",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(4), nil)

				readerRepo := repomocks.NewMockReaderArticleRepository(ctrl)
				readerRepo.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      4,
					Tittle:  "tittle",
					Content: "content",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(errors.New("update failed!")).Return(nil)
				return authorRepo, readerRepo
			},
			art: domain.Article{
				Tittle:  "tittle",
				Content: "content",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantId:  4,
			wantErr: nil,
		},
		{
			name: "线上库重试后失败",
			mock: func(ctrl *gomock.Controller) (repository.AuthorArticleRepository, repository.ReaderArticleRepository) {
				authorRepo := repomocks.NewMockAuthorArticleRepository(ctrl)
				authorRepo.EXPECT().Create(gomock.Any(), domain.Article{
					Tittle:  "tittle",
					Content: "content",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(4), nil)

				readerRepo := repomocks.NewMockReaderArticleRepository(ctrl)
				readerRepo.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      4,
					Tittle:  "tittle",
					Content: "content",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(errors.New("update failed!")).Times(3)
				return authorRepo, readerRepo
			},
			art: domain.Article{
				Tittle:  "tittle",
				Content: "content",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantErr: errors.New("update failed!"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			authorRepo, readerRepo := tc.mock(ctrl)
			svc := NewArticleService(authorRepo, readerRepo, &logger.NopLogger{})
			id, err := svc.Publish(context.Background(), tc.art)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantId, id)
		})
	}
}
