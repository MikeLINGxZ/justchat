package core

import "context"

type ICommon interface {
	Startup(ctx context.Context)
}
