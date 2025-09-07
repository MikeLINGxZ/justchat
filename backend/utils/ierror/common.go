package ierror

import "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"

type IError struct {
	errCode ErrorCode
}

func (e IError) Error() string {
	return string(e.errCode)
}

func New(errCode ErrorCode) error {
	return IError{errCode: errCode}
}

func NewError(err error) error {
	logger.Error(err)
	return IError{errCode: ErrCodeInternalError}
}
