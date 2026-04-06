package data_models

// OPCPerson 表示 OPC 模式中的一个人员，每个人员对应一个 CustomAgentDef
type OPCPerson struct {
	OrmModel
	Uuid     string `gorm:"unique;index" json:"uuid"`
	Name     string `gorm:"type:varchar(255)" json:"name"`
	Role     string `gorm:"type:varchar(255)" json:"role"`
	AgentID  string `gorm:"type:varchar(255);index" json:"agent_id"`
	Avatar   string `gorm:"type:varchar(255)" json:"avatar"`
	IsPinned bool   `gorm:"default:false" json:"is_pinned"`
	ChatUuid string `gorm:"type:varchar(255);index" json:"chat_uuid"`
}
