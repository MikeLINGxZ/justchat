package view_models

import "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"

type AppPreferences struct {
	Language data_models.AppLanguage `json:"language"`
	Region   data_models.AppRegion   `json:"region"`
}
