package ierror

type ErrorCode string

const (
	// ErrCodeInternalError 系统内部错误
	ErrCodeInternalError ErrorCode = "ErrCodeInternalError"
	// ErrCodeUsernameExists 用户名已被使用
	ErrCodeUsernameExists = "ErrCodeUsernameExists"
	// ErrCodeEmailExists 邮箱已使用
	ErrCodeEmailExists = "ErrCodeEmailExists"
	// ErrCodeLoginExpired 登陆过期
	ErrCodeLoginExpired = "ErrCodeLoginExpired"
	// ErrCodeInvalidAccountPassword 账号或密码错误
	ErrCodeInvalidAccountPassword = "ErrCodeInvalidAccountPassword"
	// ErrCodeInvalidEmailPassword 邮箱或密码错误
	ErrCodeInvalidEmailPassword = "ErrCodeInvalidEmailPassword"
	// ErrCodeLoginNotValidYet 登陆未生效
	ErrCodeLoginNotValidYet = "ErrCodeLoginNotValidYet"
	// ErrCodeInvalidVerifyCode 验证码错误
	ErrCodeInvalidVerifyCode = "ErrCodeInvalidVerifyCode"
	// ErrCodeAccountNotFound 账号不存在
	ErrCodeAccountNotFound = "ErrCodeAccountNotFound"
	// ErrCodeModelNotFound 模型不存在
	ErrCodeModelNotFound = "ErrCodeModelNotFound"
	// ErrCodeChatNotFound 对话不存在
	ErrCodeChatNotFound = "ErrCodeChatNotFound"
	// ErrCodeUnsupportedFileType 不支持的文件类型
	ErrCodeUnsupportedFileType = "ErrCodeUnsupportedFileType"
	// ErrCodeProviderNotFound 供应商不存在
	ErrCodeProviderNotFound = "ErrCodeProviderNotFound"
	// ErrCodeCompletionsParams 对话参数错误
	ErrCodeCompletionsParams = "ErrCodeCompletionsParams"
)

func MessageKey(errCode ErrorCode) string {
	switch errCode {
	case ErrCodeInternalError:
		return "errors.internal"
	case ErrCodeUsernameExists:
		return "errors.username_exists"
	case ErrCodeEmailExists:
		return "errors.email_exists"
	case ErrCodeLoginExpired:
		return "errors.login_expired"
	case ErrCodeInvalidAccountPassword:
		return "errors.invalid_account_password"
	case ErrCodeInvalidEmailPassword:
		return "errors.invalid_email_password"
	case ErrCodeLoginNotValidYet:
		return "errors.login_not_valid_yet"
	case ErrCodeInvalidVerifyCode:
		return "errors.invalid_verify_code"
	case ErrCodeAccountNotFound:
		return "errors.account_not_found"
	case ErrCodeModelNotFound:
		return "errors.model_not_found"
	case ErrCodeChatNotFound:
		return "errors.chat_not_found"
	case ErrCodeUnsupportedFileType:
		return "errors.unsupported_file_type"
	case ErrCodeProviderNotFound:
		return "errors.provider_not_found"
	case ErrCodeCompletionsParams:
		return "errors.chat_params"
	default:
		return ""
	}
}
