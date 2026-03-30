package ierror

import (
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/i18n"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"
)

type IError struct {
	errCode ErrorCode
}

func (e IError) Error() string {
	if key := MessageKey(e.errCode); key != "" {
		return i18n.TCurrent(key, nil)
	}
	return string(e.errCode)
}

func New(errCode ErrorCode) error {
	return IError{errCode: errCode}
}

func NewError(err error) error {
	logger.Error("internal error:", err)
	return IError{errCode: ErrCodeInternalError}
}
