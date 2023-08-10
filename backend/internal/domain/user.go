package domain

type User struct {
	Id               int64
	Email            string
	PassWord         string
	CTime            int64
	NickName         string
	Birthday         string
	SelfIntroduction string
}
