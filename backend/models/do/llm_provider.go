package do

type Provider struct {
	OrmModel
	ProviderName string  `gorm:"type:varchar(255)" json:"provider_name"` // 提供方名称
	BaseUrl      string  `gorm:"type:varchar(255)" json:"base_url"`      // 基础url
	ApiKey       string  `gorm:"type:varchar(255)" json:"api_key"`       // api key
	Enable       bool    `gorm:"index;type:bool" json:"enable"`          // 启用
	Alias        *string `gorm:"type:varchar(255)" json:"alias"`         // 别名
}
