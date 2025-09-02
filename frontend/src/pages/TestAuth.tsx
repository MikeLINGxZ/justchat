import React from 'react';
import { Button, Card, Typography, Space, Alert } from 'antd';
import { useAuthStore } from '@/stores/authStore';
import { authClient } from '@/api/authClient';
import { CommonVerificationCodeType } from '@/api/service/auth/Api';

const { Title, Text } = Typography;

const TestAuth: React.FC = () => {
  const { user, isAuthenticated, logout } = useAuthStore();

  const handleTestLogin = async () => {
    try {
      const response = await authClient.login({
        loginField: 'testuser',
        passwordMd5: '5f4dcc3b5aa765d61d8327deb882cf99', // password
      });
      console.log('Login response:', response);
      alert('登录成功！');
    } catch (error) {
      console.error('Login error:', error);
      alert('登录失败：' + (error as any).message);
    }
  };

  const handleTestRegister = async () => {
    try {
      const response = await authClient.register({
        username: 'testuser2',
        passwordMd5: '5f4dcc3b5aa765d61d8327deb882cf99', // password
        email: 'test2@example.com',
        emailVerificationCode: '123456',
      });
      console.log('Register response:', response);
      alert('注册成功！');
    } catch (error) {
      console.error('Register error:', error);
      alert('注册失败：' + (error as any).message);
    }
  };

  const handleTestSendCode = async () => {
    try {
      const response = await authClient.sendEmailVerificationCode({
        email: 'test@example.com',
        username: 'testuser',
        codeType: CommonVerificationCodeType.VERIFICATION_CODE_TYPE_REGISTER,
      });
      console.log('Send code response:', response);
      alert('验证码发送成功！');
    } catch (error) {
      console.error('Send code error:', error);
      alert('验证码发送失败：' + (error as any).message);
    }
  };

  return (
    <div style={{ padding: 24, maxWidth: 800, margin: '0 auto' }}>
      <Title level={2}>认证功能测试</Title>
      
      <Card style={{ marginBottom: 24 }}>
        <Title level={3}>当前状态</Title>
        <Space direction="vertical" style={{ width: '100%' }}>
          <Text>认证状态: {isAuthenticated ? '已登录' : '未登录'}</Text>
          {user && (
            <div>
              <Text>用户信息:</Text>
              <pre>{JSON.stringify(user, null, 2)}</pre>
            </div>
          )}
        </Space>
      </Card>

      <Card style={{ marginBottom: 24 }}>
        <Title level={3}>API 测试</Title>
        <Space direction="vertical" style={{ width: '100%' }}>
          <Button type="primary" onClick={handleTestLogin}>
            测试登录
          </Button>
          <Button onClick={handleTestRegister}>
            测试注册
          </Button>
          <Button onClick={handleTestSendCode}>
            测试发送验证码
          </Button>
          {isAuthenticated && (
            <Button danger onClick={logout}>
              登出
            </Button>
          )}
        </Space>
      </Card>

      <Alert
        message="测试说明"
        description="这些测试按钮会直接调用后端API，请确保后端服务正在运行。测试数据仅供参考，实际使用时请使用真实数据。"
        type="info"
        showIcon
      />
    </div>
  );
};

export default TestAuth;
