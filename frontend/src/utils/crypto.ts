import MD5 from 'crypto-js/md5';

/**
 * 将密码转换为MD5哈希
 * @param password 原始密码
 * @returns MD5哈希字符串
 */
export const hashPassword = (password: string): string => {
  return MD5(password).toString();
};

/**
 * 验证密码强度
 * @param password 密码
 * @returns 验证结果
 */
export const validatePassword = (password: string): {
  isValid: boolean;
  message?: string;
} => {
  if (password.length < 6) {
    return { isValid: false, message: '密码至少6位字符' };
  }
  
  if (!/^(?=.*[a-zA-Z])(?=.*\d)/.test(password)) {
    return { isValid: false, message: '密码必须包含字母和数字' };
  }
  
  return { isValid: true };
};
