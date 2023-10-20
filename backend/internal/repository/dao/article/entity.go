package article

type PublishArticle struct {
	Article
}

type Article struct {
	Id       int64  `gorm:"primaryKey,autoIncrement" bson:"id,omitempty"`
	Tittle   string `gorm:"type=varchar(4096)" bson:"tittle,omitempty"`
	Content  string `gorm:"type=BLOB" bson:"content,omitempty"`
	AuthorId int64  `gorm:"index" bson:"author_id,omitempty"`
	Status   uint8  `bson:"status,omitempty"`
	Ctime    int64  `bson:"ctime,omitempty"`
	Utime    int64  `bson:"utime,omitempty"`
}
