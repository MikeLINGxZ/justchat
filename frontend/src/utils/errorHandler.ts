import i18n from '@/i18n';

interface ApiError {
  code: number;
  message: string;
  details: any[];
}

const ERROR_MESSAGE_MAP: Record<string, string> = {
  ErrCodeInternalError: 'errors.internal',
  ErrCodeModelNotFound: 'errors.modelNotFound',
  ErrCodeProviderNotFound: 'errors.providerNotFound',
  ErrCodeCompletionsParams: 'errors.chatParams',
  ErrCodeUsernameExists: 'errors.usernameExists',
  ErrCodeEmailExists: 'errors.emailExists',
  ErrCodeLoginExpired: 'errors.loginExpired',
  ErrCodeInvalidAccountPassword: 'errors.invalidAccountPassword',
  ErrCodeInvalidEmailPassword: 'errors.invalidEmailPassword',
  ErrCodeLoginNotValidYet: 'errors.loginNotValidYet',
  ErrCodeInvalidVerifyCode: 'errors.invalidVerifyCode',
  ErrCodeAccountNotFound: 'errors.accountNotFound',
  ErrCodeChatNotFound: 'errors.chatNotFound',
  ErrCodeUnsupportedFileType: 'errors.unsupportedFileType',
};

const DEFAULT_ERROR_MESSAGE_KEY = 'common.operationFailed';

export const translateError = (error: any): string => {
  try {
    let apiError: ApiError | string | null = null;
    if (error?.error?.message != null) {
      apiError = error.error.message;
    } else {
      apiError = error;
    }

    if (apiError && typeof apiError === 'string') {
      const key = ERROR_MESSAGE_MAP[apiError];
      if (key) {
        return i18n.t(key);
      }
      return apiError;
    }

    if (error?.response?.status) {
      switch (error.response.status) {
        case 400:
          return i18n.t('errors.badRequest');
        case 401:
          return i18n.t('errors.unauthorized');
        case 403:
          return i18n.t('errors.forbidden');
        case 404:
          return i18n.t('errors.notFound');
        case 429:
          return i18n.t('errors.tooManyRequests');
        case 500:
          return i18n.t('errors.internal');
        case 502:
        case 503:
          return i18n.t('errors.serviceUnavailable');
        case 504:
          return i18n.t('errors.timeout');
        default:
          return i18n.t(DEFAULT_ERROR_MESSAGE_KEY);
      }
    }

    if (error?.message) {
      if (error.message.includes('Network Error') || error.message.includes('timeout')) {
        return i18n.t('errors.networkError');
      }

      if (error.message.includes('ECONNREFUSED') || error.message.includes('ERR_CONNECTION_REFUSED')) {
        return i18n.t('errors.connectionRefused');
      }
    }

    return apiError && typeof apiError !== 'string'
      ? apiError.message || i18n.t(DEFAULT_ERROR_MESSAGE_KEY)
      : i18n.t(DEFAULT_ERROR_MESSAGE_KEY);
  } catch (e) {
    console.error('Error while extracting error message:', e);
    return i18n.t(DEFAULT_ERROR_MESSAGE_KEY);
  }
};

export const extractErrorMessage = translateError;

export const addErrorMapping = (errorCode: string, messageKey: string): void => {
  ERROR_MESSAGE_MAP[errorCode] = messageKey;
};

export const getErrorMappings = (): Record<string, string> => {
  return { ...ERROR_MESSAGE_MAP };
};

export const isApiError = (error: any): error is ApiError => {
  return error &&
    typeof error.code === 'number' &&
    typeof error.message === 'string' &&
    Array.isArray(error.details);
};

export const extractErrorCode = (error: any): number | null => {
  try {
    if (error?.response?.data?.code !== undefined) {
      return error.response.data.code;
    }
    if (error?.data?.code !== undefined) {
      return error.data.code;
    }
    if (error?.code !== undefined) {
      return error.code;
    }
    return null;
  } catch {
    return null;
  }
};
