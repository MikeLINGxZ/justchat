package data_models

// Session stores one chat session, including normal chats and automated task chats.
type Session struct {
	OrmModel
	Title   string `gorm:"type:varchar(255)" json:"title"`
	Starred bool   `gorm:"type:bool;default:0;index" json:"starred"`
	Status  string `gorm:"type:varchar(50);default:'idle'" json:"status"`
	Kind    string `gorm:"type:varchar(20);default:'user';index" json:"kind"`
	Tags    string `gorm:"type:text" json:"tags"`
}
