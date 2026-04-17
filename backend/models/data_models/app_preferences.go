package data_models

type AppLanguage string

const (
	AppLanguageZhCN AppLanguage = "zh-CN"
	AppLanguageEnUS AppLanguage = "en-US"
)

type AppRegion string

const (
	AppRegionAsia         AppRegion = "asia"
	AppRegionEurope       AppRegion = "europe"
	AppRegionNorthAmerica AppRegion = "north-america"
	AppRegionSouthAmerica AppRegion = "south-america"
	AppRegionAfrica       AppRegion = "africa"
	AppRegionOceania      AppRegion = "oceania"
	AppRegionAntarctica   AppRegion = "antarctica"
)

const AppPreferencesSingletonID uint = 1

type AppPreferences struct {
	OrmModel            `json:"-"`
	SingletonID         uint        `gorm:"uniqueIndex" json:"singleton_id"`
	Language            AppLanguage `gorm:"type:varchar(16);default:'zh-CN'" json:"language"`
	Region              AppRegion   `gorm:"type:varchar(32);default:'asia'" json:"region"`
	MemorySystemEnabled bool        `gorm:"default:false" json:"memory_system_enabled"`
	VectorSearchEnabled bool        `gorm:"default:false" json:"vector_search_enabled"`
	EmbeddingProvider   string      `gorm:"type:varchar(64)" json:"embedding_provider"`
	EmbeddingBaseURL    string      `gorm:"type:text" json:"embedding_base_url"`
	EmbeddingAPIKey     string      `gorm:"type:text" json:"embedding_api_key"`
	EmbeddingModel      string      `gorm:"type:varchar(255)" json:"embedding_model"`
}
