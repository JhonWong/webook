package dao

type Article struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Tittle   string `gorm:"type=varchar(4096)"`
	Content  string `gorm:"type=BLOB"`
	AuthorId int64  `gorm:"index"`
	CTime    int64
	UTime    int64
}
