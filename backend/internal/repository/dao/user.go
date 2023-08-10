package dao

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
	ErrUserNotFound       = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (dao *UserDAO) Insert(ctx *gin.Context, u User) error {
	now := time.Now().UnixMilli()
	u.CTime = now
	u.UTime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			return ErrUserDuplicateEmail
		}
	}
	return err
}

func (dao *UserDAO) FindByEmail(ctx *gin.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (dao *UserDAO) FindById(ctx *gin.Context, Id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id = ?", Id).First(&u).Error
	return u, err
}

func (dao *UserDAO) Update(ctx *gin.Context, u User) error {
	now := time.Now().UnixMicro()
	u.UTime = now
	err := dao.db.WithContext(ctx).Save(&u).Error
	return err
}

// 与表结构对应
type User struct {
	Id               int64  `gorm:"primaryKey,autoIncrement"`
	Email            string `gorm:"unique"`
	Password         string
	CTime            int64
	UTime            int64
	NickName         string
	Birthday         string
	SelfIntroduction string
}
