package settings_dto

import "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_model"

type LoadBootstrapInput struct {
}

type LoadBootstrapOutput struct {
	Bootstrap view_model.SettingsBootstrap `json:"bootstrap"`
}
