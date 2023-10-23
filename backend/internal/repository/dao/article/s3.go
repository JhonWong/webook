package article

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ecodeclub/ekit"
	"github.com/johnwongx/webook/backend/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"time"
)

var _ ArticleDAO = &S3DAO{}

type S3DAO struct {
	oss    *s3.S3
	bucket *string
	GORMArticleDAO
}

func NewS3DAO(oss *s3.S3, db *gorm.DB) *S3DAO {
	return &S3DAO{
		oss:    oss,
		bucket: ekit.ToPtr[string]("webook-1314583317"),
		GORMArticleDAO: GORMArticleDAO{
			db: db,
		},
	}
}

func (s *S3DAO) Sync(ctx context.Context, art Article) (int64, error) {
	// 使用事物保证两张表同时成功或失败
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var (
			id  = art.Id
			err error
		)

		// 更新制作库，插入或删除
		if id > 0 {
			err = s.UpdateById(ctx, art)
		} else {
			id, err = s.Insert(ctx, art)
		}
		if err != nil {
			return err
		}
		// 更新数据到线上库
		art.Id = id

		pArt := PublishArticle(art)
		now := time.Now().UnixMilli()
		pArt.Ctime = now
		pArt.Utime = now

		return s.db.WithContext(ctx).
			Clauses(clause.OnConflict{
				DoUpdates: clause.Assignments(map[string]interface{}{
					"title":  art.Title,
					"utime":  art.Utime,
					"status": art.Status,
				}),
			}).Create(&pArt).Error
	})
	if err != nil {
		return 0, err
	}
	_, err = s.oss.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      s.bucket,
		Key:         ekit.ToPtr[string](strconv.FormatInt(art.Id, 10)),
		Body:        bytes.NewReader([]byte(art.Content)),
		ContentType: ekit.ToPtr[string]("text/plain;charset=utf-8"),
	})
	return art.Id, err
}

func (s *S3DAO) SyncStatus(ctx context.Context, id, usrId int64, status uint8) error {
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UnixMilli()
		res := tx.Model(&Article{}).
			Where("id = ? AND author_id = ?", id, usrId).
			Updates(map[string]any{
				"status": status,
				"utime":  now,
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return fmt.Errorf("可能有人在攻击系统，误操作非自己的文章, Uid:%d, authorId:", id, usrId)
		}

		return tx.Model(&PublishArticle{}).
			Where("id = ?", id).
			Updates(map[string]any{
				"status": status,
				"utime":  now,
			}).Error
	})
	if err != nil {
		return err
	}
	if status == domain.ArticleStatusPrivate.ToUint8() {
		_, err = s.oss.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
			Bucket: s.bucket,
			Key:    ekit.ToPtr[string](strconv.FormatInt(id, 10)),
		})
	}
	return err
}
