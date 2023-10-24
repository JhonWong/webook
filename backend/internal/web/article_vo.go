package web

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
