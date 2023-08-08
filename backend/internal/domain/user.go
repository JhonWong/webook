package domain

type User struct {
	Id       int64
	Email    []byte
	PassWord []byte
	CTime    int64
}
