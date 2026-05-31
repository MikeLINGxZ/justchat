package config_dto

import "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_model"

type LanguageListInput struct {
}

type LanguageListOutput struct {
	Languages []view_model.Language `json:"languages"`
}
