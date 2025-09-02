package core

import "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/rpc/service"

type IAuth interface {
	ICommon
	Login(request *service.LoginRequest) error
	Register(request *service.RegisterRequest) error
}
