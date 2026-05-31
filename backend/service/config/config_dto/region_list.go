package config_dto

import "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_model"

type RegionListInput struct {
}

type RegionListOutput struct {
	Regions []view_model.Region `json:"regions"`
}
