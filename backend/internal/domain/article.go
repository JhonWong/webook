package domain

type ArticleStatus uint8

const (
	ArticleStatusUnknow ArticleStatus = iota
	ArticleStatusUnpublished
	ArticleStatusPublished
	ArticleStatusPrivate
)

func (a ArticleStatus) ToUint8() uint8 {
	return uint8(a)
}

type Article struct {
	Id      int64
	Tittle  string
	Content string
	Author  Author
	Status  ArticleStatus
}

type Author struct {
	Id int64
}
