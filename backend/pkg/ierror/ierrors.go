package ierror

import (
	"encoding/json"
	"errors"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/i18n"
)

// IError is the standard error type returned by all Wails-bound service methods.
// Detail carries the raw error string; Msg is the localized user-facing message.
type IError struct {
	Detail string `json:"detail"`
	Msg    string `json:"msg"`
}

// Error returns the JSON representation of the IError.
func (e *IError) Error() string {
	errBytes, err := json.Marshal(e)
	if err != nil {
		return i18n.TCurrent("ierror.unknown_error", nil)
	}
	return string(errBytes)
}

// Is enables errors.Is to compare IError values by Msg.
func (e *IError) Is(target error) bool {
	if t, ok := target.(*IError); ok {
		return e.Msg == t.Msg
	}
	return false
}

// Error wraps err in an IError using the given error code for the localized Msg.
// Returns nil when err is nil. Passes through an existing *IError unchanged.
func Error(code errorCode, err error) error {
	if err == nil {
		return nil
	}

	var iErr *IError
	if errors.As(err, &iErr) {
		return iErr
	}

	return &IError{
		Detail: err.Error(),
		Msg:    code.Msg(),
	}
}
