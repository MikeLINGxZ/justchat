package file_uploader

import (
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/wrapper_models"
)

type Aliyuns struct {
	providerModel *wrapper_models.ProviderModel
}

func (a *Aliyuns) Upload(paths []string) (map[string]string, error) {
	//TODO implement me
	panic("implement me")
}
