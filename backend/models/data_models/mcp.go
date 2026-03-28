package data_models

type CustomMCPServer struct {
	OrmModel
	Name        string `gorm:"type:varchar(255)" json:"name"`
	SourcePath  string `gorm:"uniqueIndex;type:text" json:"source_path"`
	ConfigPath  string `gorm:"type:text" json:"config_path"`
	ToolID      string `gorm:"uniqueIndex;type:varchar(255)" json:"tool_id"`
	Description string `gorm:"type:text" json:"description"`
	Enabled     bool   `gorm:"index;type:bool;default:1" json:"enabled"`
}
