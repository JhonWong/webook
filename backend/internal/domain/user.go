package domain

type User struct {
	Id               int64
	Email            string
	Phone            string
	WechatInfo       WechatInfo
	PassWord         string
	CTime            int64
	NickName         string
	Birthday         string
	SelfIntroduction string
}
