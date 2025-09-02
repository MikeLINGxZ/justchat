package cloud

import (
	"context"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/core"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/rpc/service"
)

type Auth struct {
	ctx context.Context
}

func (a *Auth) Startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *Auth) Login(request *service.LoginRequest) error {
	//TODO implement me
	panic("implement me")
}

func (a *Auth) Register(request *service.RegisterRequest) error {
	//TODO implement me
	panic("implement me")
}

func NewAuth() core.IAuth {
	return &Auth{}
}
