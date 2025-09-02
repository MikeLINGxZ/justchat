import React, { useState, useRef, useEffect } from 'react';
import { Form, Input, Button, Card, Typography, Alert, Steps } from 'antd';
import { MailOutlined, LockOutlined, SafetyOutlined } from '@ant-design/icons';
import { Link, useNavigate } from 'react-router-dom';
import { authClient } from '@/api/authClient';
import { CommonVerificationCodeType } from '@/api/service/auth/Api';
import { hashPassword } from '@/utils/crypto';
import styles from './style.module.scss';

const { Title, Text } = Typography;
const { Step } = Steps;

interface ResetPasswordForm {
  loginField: string;
  emailVerificationCode: string;
  newPassword: string;
  confirmPassword: string;
}

const ForgotPassword: React.FC = () => {
  const navigate = useNavigate();
  const [form] = Form.useForm<ResetPasswordForm>();
  const [currentStep, setCurrentStep] = useState(0);
  const [loading, setLoading] = useState(false);
  const [sendingCode, setSendingCode] = useState(false);
  const [countdown, setCountdown] = useState(0);
  const [error, setError] = useState('');
  const [email, setEmail] = useState('');
  const timerRef = useRef<NodeJS.Timeout | null>(null);

  useEffect(() => {
    return () => {
      if (timerRef.current) {
        clearInterval(timerRef.current);
      }
    };
  }, []);

  // 发送验证码
  const handleSendVerificationCode = async () => {
    try {
      const loginField = form.getFieldValue('loginField');
      if (!loginField) {
        setError('请先输入用户名或邮箱');
        return;
      }

      setSendingCode(true);
      setError('');

      await authClient.sendEmailVerificationCode({
        email: loginField,
        codeType: CommonVerificationCodeType.VERIFICATION_CODE_TYPE_RESET_PASSWORD,
      });

      setEmail(loginField);
      setCurrentStep(1);
      
      // 启动倒计时
      setCountdown(60);
      timerRef.current = setInterval(() => {
        setCountdown((prev) => {
          if (prev <= 1) {
            if (timerRef.current) {
              clearInterval(timerRef.current);
              timerRef.current = null;
            }
            return 0;
          }
          return prev - 1;
        });
      }, 1000);
    } catch (error: any) {
      setError(error.message || '发送验证码失败，请重试');
    } finally {
      setSendingCode(false);
    }
  };

  // 重置密码
  const handleResetPassword = async (values: ResetPasswordForm) => {
    if (values.newPassword !== values.confirmPassword) {
      setError('两次输入的密码不一致');
      return;
    }

    try {
      setLoading(true);
      setError('');

      await authClient.resetPassword({
        loginField: values.loginField,
        emailVerificationCode: values.emailVerificationCode,
        newPasswordMd5: hashPassword(values.newPassword),
      });

      setCurrentStep(2);
    } catch (error: any) {
      setError(error.message || '重置密码失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  const handleErrorClose = () => {
    setError('');
  };

  const renderStepContent = () => {
    switch (currentStep) {
      case 0:
        return (
          <>
            <Form.Item
              name="loginField"
              rules={[
                { required: true, message: '请输入用户名或邮箱' },
                { min: 3, message: '用户名至少3位字符' },
              ]}
            >
              <Input
                prefix={<MailOutlined />}
                placeholder="用户名或邮箱"
                size="large"
              />
            </Form.Item>

            <Form.Item>
              <Button
                type="primary"
                onClick={handleSendVerificationCode}
                loading={sendingCode}
                block
                size="large"
              >
                发送验证码
              </Button>
            </Form.Item>
          </>
        );

      case 1:
        return (
          <>
            <Form.Item
              name="loginField"
              initialValue={email}
              rules={[
                { required: true, message: '请输入用户名或邮箱' },
              ]}
            >
              <Input
                prefix={<MailOutlined />}
                placeholder="用户名或邮箱"
                size="large"
                disabled
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
                prefix={<SafetyOutlined />}
                placeholder="邮箱验证码"
                size="large"
                suffix={
                  <Button
                    type="link"
                    size="small"
                    loading={sendingCode}
                    onClick={handleSendVerificationCode}
                    disabled={countdown > 0}
                  >
                    {countdown > 0 ? `${countdown}s后重发` : '重新发送'}
                  </Button>
                }
              />
            </Form.Item>

            <Form.Item
              name="newPassword"
              rules={[
                { required: true, message: '请输入新密码' },
                { min: 6, message: '密码至少6位字符' },
                {
                  pattern: /^(?=.*[a-zA-Z])(?=.*\d)[a-zA-Z\d@$!%*?&]{6,}$/,
                  message: '密码必须包含字母和数字',
                },
              ]}
            >
              <Input.Password
                prefix={<LockOutlined />}
                placeholder="新密码"
                size="large"
              />
            </Form.Item>

            <Form.Item
              name="confirmPassword"
              dependencies={['newPassword']}
              rules={[
                { required: true, message: '请确认新密码' },
                ({ getFieldValue }) => ({
                  validator(_, value) {
                    if (!value || getFieldValue('newPassword') === value) {
                      return Promise.resolve();
                    }
                    return Promise.reject(new Error('两次输入的密码不一致'));
                  },
                }),
              ]}
            >
              <Input.Password
                prefix={<LockOutlined />}
                placeholder="确认新密码"
                size="large"
              />
            </Form.Item>

            <Form.Item>
              <Button
                type="primary"
                htmlType="submit"
                loading={loading}
                block
                size="large"
              >
                重置密码
              </Button>
            </Form.Item>
          </>
        );

      case 2:
        return (
          <div style={{ textAlign: 'center', padding: '20px 0' }}>
            <div style={{ fontSize: '40px', color: '#52c41a', marginBottom: '12px' }}>
              ✓
            </div>
            <Title level={4} style={{ color: '#52c41a', marginBottom: '8px' }}>
              密码重置成功！
            </Title>
            <Text type="secondary" style={{ display: 'block', marginBottom: '20px' }}>
              您的密码已成功重置，请使用新密码登录
            </Text>
            <Button
              type="primary"
              size="large"
              onClick={() => navigate('/login')}
            >
              立即登录
            </Button>
          </div>
        );

      default:
        return null;
    }
  };

  return (
    <div className={styles.authContainer}>
      <div className={styles.authCard}>
        <Card>
          <div className={styles.authHeader}>
            <Title level={2}>忘记密码</Title>
            <Text type="secondary">
              {currentStep === 0 && '请输入您的用户名或邮箱地址'}
              {currentStep === 1 && '请输入验证码并设置新密码'}
              {currentStep === 2 && '密码重置完成'}
            </Text>
          </div>

          {currentStep < 2 && (
            <Steps 
              current={currentStep} 
              size="small"
              style={{ marginBottom: '24px' }}
            >
              <Step title="验证身份" />
              <Step title="重置密码" />
              <Step title="完成" />
            </Steps>
          )}

          {error && (
            <Alert
              message={error}
              type="error"
              closable
              onClose={handleErrorClose}
              style={{ marginBottom: 16 }}
            />
          )}

          {currentStep < 2 && (
            <Form
              form={form}
              name="forgotPassword"
              onFinish={handleResetPassword}
              autoComplete="off"
            >
              {renderStepContent()}
            </Form>
          )}

          {currentStep === 2 && renderStepContent()}

          {currentStep < 2 && (
            <div className={styles.authFooter}>
              <div className={styles.authSwitch}>
                <Text>想起密码了？</Text>
                <Link to="/login">立即登录</Link>
              </div>
            </div>
          )}
        </Card>
      </div>
    </div>
  );
};

export default ForgotPassword;