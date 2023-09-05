package dao

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/assert/v2"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGORMUserDAO_Insert(t *testing.T) {
	testCase := []struct {
		name    string
		mock    func(t *testing.T) *sql.DB
		user    User
		wantErr error
	}{
		{
			name: "success",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				res := sqlmock.NewResult(3, 1)
				mock.ExpectExec("INSERT INTO `users` .*").
					WillReturnResult(res)
				require.NoError(t, err)
				return mockDB
			},
			user: User{
				Email: sql.NullString{
					String: "123@qq.com",
					Valid:  true,
				},
			},
		},
		{
			name: "email duplicate",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				mock.ExpectExec("INSERT INTO `users` .*").
					WillReturnError(&mysql.MySQLError{
						Number: 1062,
					})
				require.NoError(t, err)
				return mockDB
			},
			user:    User{},
			wantErr: ErrUserDuplicateEmail,
		},
		{
			name: "database error",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				mock.ExpectExec("INSERT INTO `users` .*").
					WillReturnError(errors.New("database error"))
				require.NoError(t, err)
				return mockDB
			},
			user:    User{},
			wantErr: errors.New("database error"),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			db, err := gorm.Open(gormMysql.New(gormMysql.Config{
				Conn:                      tc.mock(t),
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				DisableAutomaticPing:   true,
				SkipDefaultTransaction: true,
			})
			d := NewUserDAO(db)
			u := tc.user
			err = d.Insert(nil, u)
			assert.Equal(t, err, tc.wantErr)
		})
	}
}
