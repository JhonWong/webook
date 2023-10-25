package domain

import "time"

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
	Title   string
	Content string
	Author  Author
	Status  ArticleStatus
	Ctime   time.Time
	Utime   time.Time
}

func (a *Article) Abstract() string {
	cs := []rune(a.Content)
	if len(cs) < 100 {
		return a.Content
	}
	return string(a.Content[:100])
}

type Author struct {
	Id   int64
	Name string
}
