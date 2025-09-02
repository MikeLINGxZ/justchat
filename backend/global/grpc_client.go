package global

import (
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/rpc/service"
)

var GRPC *GrpcClient

func init() {
	auth := service.NewAuthClient(utils.NewGrpcConn())
	GRPC = &GrpcClient{
		Auth: auth,
	}
}

type GrpcClient struct {
	Auth service.AuthClient
}
