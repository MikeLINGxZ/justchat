package data_models

// OPCGroup 表示 OPC 模式中的一个群聊
type OPCGroup struct {
	OrmModel
	Uuid        string `gorm:"unique;index" json:"uuid"`
	ChatUuid    string `gorm:"unique;index" json:"chat_uuid"`
	Name        string `gorm:"type:varchar(255)" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	IsPinned    bool   `gorm:"default:false" json:"is_pinned"`
}

// OPCGroupMember 表示群聊成员关系
type OPCGroupMember struct {
	OrmModel
	GroupUuid  string `gorm:"index" json:"group_uuid"`
	PersonUuid string `gorm:"index" json:"person_uuid"`
}
