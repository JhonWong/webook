package domain

type Article struct {
	Id      int64
	Tittle  string
	Content string
	Author  Author
}

type Author struct {
	Id int64
}
