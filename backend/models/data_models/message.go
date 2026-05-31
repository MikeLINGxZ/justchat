package data_models

type Message struct {
	OrmModel
	SessionID   uint   `gorm:"index" json:"session_id"`
	ParentID    *uint  `gorm:"index" json:"parent_id"`
	Role        string `gorm:"type:varchar(50)" json:"role"`
	ContentType string `gorm:"type:varchar(50)" json:"content_type"`
	Content     string `gorm:"type:text" json:"content"`
	ModelName   string `gorm:"type:varchar(255)" json:"model_name"`
	AgentName   string `gorm:"type:varchar(255)" json:"agent_name"`
	TokensIn    int    `json:"tokens_in"`
	TokensOut   int    `json:"tokens_out"`
	Extra       string `gorm:"type:text" json:"extra"`
	Attachments string `gorm:"type:text" json:"attachments"`
}
