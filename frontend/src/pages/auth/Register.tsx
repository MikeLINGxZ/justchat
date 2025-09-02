import React from 'react';
import { Form, Input, Button, Card, Typography, Alert } from 'antd';
import { UserOutlined, LockOutlined, MailOutlined } from '@ant-design/icons';
import { Link } from 'react-router-dom';
import { useRegister } from '@/hooks/useRegister';
import { useSendVerificationCode } from '@/hooks/useSendVerificationCode';
import styles from './style.module.scss';

const { Title, Text } = Typography;

const Register: React.FC = () => {
  const { register, isLoading, error, clearError } = useRegister();
  const { 
    isSending: sendingCode, 
    codeSent, 
    countdown, 
    error: sendCodeError, 
    sendCode, 
    clearError: clearSendCodeError 
  } = useSendVerificationCode();
  const [form] = Form.useForm();
  
  const handleSendVerificationCode = async () => {
    const email = form.getFieldValue('email');
    const username = form.getFieldValue('username');
    
    if (!email || !username) {
      // 这里可以使用form的验证来提示用户
      form.validateFields(['email', 'username']).catch(() => {
        // 验证失败时的处理
      });
      return;
    }

    await sendCode({ email, username });
  };

  const handleSubmit = async (values: {
    username: string;
    email: string;
    password: string;
    confirmPassword: string;
    emailVerificationCode: number;
  }) => {
    // 清除之前的错误状态
    clearError();
    
    // 调用注册逻辑
    await register({
      username: values.username,
      email: values.email,
      password: values.password,
      emailVerificationCode: values.emailVerificationCode.toString(),
    });
  };

  const handleErrorClose = () => {
    clearError();
  };

  return (
    <div className={styles.authContainer}>
      <div className={styles.authCard}>
        <Card>
          <div className={styles.authHeader}>
            <Title level={2}>注册</Title>
            <Text type="secondary">创建您的 Lemon Tea 账号</Text>
          </div>

          <Form
            form={form}
            name="register"
            onFinish={handleSubmit}
            autoComplete="off"
            size="large"
          >
            <Form.Item
              name="username"
              rules={[
                { required: true, message: '请输入用户名' },
                { min: 3, message: '用户名至少3位字符' },
                { max: 20, message: '用户名最多20位字符' },
                {
                  pattern: /^[a-zA-Z0-9_\u4e00-\u9fa5]+$/,
                  message: '用户名只能包含字母、数字、下划线和中文',
                },
              ]}
            >
              <Input
                prefix={<UserOutlined />}
                placeholder="用户名"
                autoComplete="username"
              />
            </Form.Item>

            <Form.Item
              name="email"
              rules={[
                { required: true, message: '请输入邮箱地址' },
                { type: 'email', message: '请输入有效的邮箱地址' },
              ]}
            >
              <Input
                prefix={<MailOutlined />}
                placeholder="邮箱地址"
                autoComplete="email"
              />
            </Form.Item>

            <Form.Item
              name="emailVerificationCode"
              rules={[
                { required: true, message: '请输入邮箱验证码' },
                { pattern: /^\d+$/, message: '验证码必须是数字' },
              ]}
            >
              <Input
                placeholder="邮箱验证码"
                suffix={
                  <Button
                    type="link"
                    size="small"
                    loading={sendingCode}
                    onClick={handleSendVerificationCode}
                    disabled={codeSent || countdown > 0}
                  >
                    {countdown > 0 ? `${countdown}s后重发` : '发送验证码'}
                  </Button>
                }
              />
            </Form.Item>

            <Form.Item
              name="password"
              rules={[
                { required: true, message: '请输入密码' },
                { min: 6, message: '密码至少6位字符' },
                {
                  pattern: /^(?=.*[a-zA-Z])(?=.*\d)[a-zA-Z\d@$!%*?&]{6,}$/,
                  message: '密码必须包含字母和数字',
                },
              ]}
            >
              <Input.Password
                prefix={<LockOutlined />}
                placeholder="密码"
                autoComplete="new-password"
              />
            </Form.Item>

            <Form.Item
              name="confirmPassword"
              dependencies={['password']}
              rules={[
                { required: true, message: '请确认密码' },
                ({ getFieldValue }) => ({
                  validator(_, value) {
                    if (!value || getFieldValue('password') === value) {
                      return Promise.resolve();
                    }
                    return Promise.reject(new Error('两次输入的密码不一致'));
                  },
                }),
              ]}
            >
              <Input.Password
                prefix={<LockOutlined />}
                placeholder="确认密码"
                autoComplete="new-password"
              />
            </Form.Item>

            {/* 验证码错误容器 - 固定高度 */}
            <div className={`${styles.codeErrorContainer} ${!sendCodeError ? styles.empty : ''}`}>
              {sendCodeError && (
                <Alert
                  message={sendCodeError}
                  type="error"
                  closable
                  onClose={clearSendCodeError}
                />
              )}
            </div>

            {/* 注册错误消息容器 - 放在注册按钮之前 */}
            <div className={`${styles.errorContainer} ${!error ? styles.empty : ''}`}>
              {error && (
                <Alert
                  message={error}
                  type="error"
                  closable
                  onClose={handleErrorClose}
                />
              )}
            </div>

            <Form.Item>
              <Button
                type="primary"
                htmlType="submit"
                loading={isLoading}
                block
              >
                注册
              </Button>
            </Form.Item>
          </Form>

          <div className={styles.authFooter}>
            <div className={styles.authSwitch}>
              <Text>已有账号？</Text>
              <Link to="/login">立即登录</Link>
            </div>
          </div>
        </Card>
      </div>
    </div>
  );
};

export default Register;