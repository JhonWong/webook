package web

import "github.com/johnwongx/webook/backend/internal/domain"

// VO: view object 对标前端
type ArticleVO struct {
	Id       int64  `json:"id"`
	Title    string `json:"title"`
	Abstract string `json:"abstract"`
	Content  string `json:"content"`
	Author   string `json:"author"`
	Status   uint8  `json:"status"`
	Ctime    string `json:"ctime"`
	Utime    string `json:"utime"`
}

type CollectReq struct {
	Id  int64 `json:"id"`
	CId int64 `json:"c_id"`
}

type LikeReq struct {
	Id     int64 `json:"id"`
	IsLike bool  `json:"is_like"`
}

type ListReq struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type WithdrawReq struct {
	Id int64 `json:"id"`
}

type ArticleReq struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (a *ArticleReq) toDomain(uid int64) domain.Article {
	return domain.Article{
		Id:      a.Id,
		Title:   a.Title,
		Content: a.Content,
		Author: domain.Author{
			Id: uid,
		},
	}
}
