import React, {useCallback, useEffect} from 'react';
import {Alert, Button, Card, Form, Input, Space, Typography} from 'antd';
import {LockOutlined, UserOutlined} from '@ant-design/icons';
import {Link, useLocation} from 'react-router-dom';
import {useLogin} from '@/hooks/useLogin';
import styles from './style.module.scss';

const {Title, Text} = Typography;

const Login: React.FC = () => {
    const location = useLocation();
    const { login, isLoading, error, clearError } = useLogin();
    const [form] = Form.useForm();
    
    // 从注册页面传递的状态
    const registrationMessage = (location.state as any)?.message;
    const registrationUsername = (location.state as any)?.username;
    
    // 如果有注册传递的用户名，自动填充
    useEffect(() => {
        if (registrationUsername) {
            form.setFieldsValue({ username: registrationUsername });
        }
    }, [registrationUsername, form]);

    const handleSubmit = useCallback(async (values: { username: string; password: string }) => {
        // 清除之前的错误状态
        clearError();

        // 调用登录逻辑
        await login({
            username: values.username,
            password: values.password,
        });
    }, [login, clearError]);

    const handleErrorClose = useCallback(() => {
        clearError();
    }, [clearError]);

    const handleFinishFailed = useCallback((errorInfo: any) => {
        console.log('Form validation failed:', errorInfo);
        // 阻止表单的默认提交行为
    }, []);



    return (
        <div className={styles.authContainer}>
            <div className={styles.authCard}>
                <Card>
                    <div className={styles.authHeader}>
                        <Title level={2}>登录</Title>
                        <Text type="secondary">欢迎回到 Lemon Tea</Text>
                    </div>

                    {/* 成功消息容器 - 固定高度 */}
                    <div className={`${styles.successContainer} ${!registrationMessage ? styles.empty : ''}`}>
                        {registrationMessage && (
                            <Alert
                                message={registrationMessage}
                                type="success"
                                closable
                            />
                        )}
                    </div>

                    <Form
                        form={form}
                        name="login"
                        onFinish={handleSubmit}
                        onFinishFailed={handleFinishFailed}
                        autoComplete="off"
                        size="large"
                        preserve={false}
                    >
                        <Form.Item
                            name="username"
                            rules={[
                                {required: true, message: '请输入用户名或邮箱'},
                                {min: 3, message: '用户名至少3位字符'},
                            ]}
                        >
                            <Input
                                prefix={<UserOutlined/>}
                                placeholder="用户名或邮箱"
                                autoComplete="username"
                            />
                        </Form.Item>

                        <Form.Item
                            name="password"
                            rules={[
                                {required: true, message: '请输入密码'},
                                {min: 6, message: '密码至少6位字符'},
                            ]}
                        >
                            <Input.Password
                                prefix={<LockOutlined/>}
                                placeholder="密码"
                                autoComplete="current-password"
                            />
                        </Form.Item>

                        <Form.Item>
                            {/* 错误消息容器 - 放在密码输入框和登录按钮之间 */}
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
                            <Button
                                type="primary"
                                htmlType="submit"
                                loading={isLoading}
                                block
                            >
                                登录
                            </Button>
                        </Form.Item>
                    </Form>

                    <div className={styles.authFooter}>
                        <Space direction="vertical" size="small" style={{width: '100%'}}>
                            <div className={styles.authLinks}>
                                <Link to="/forgot-password">忘记密码？</Link>
                            </div>
                            <div className={styles.authSwitch}>
                                <Text>还没有账号？</Text>
                                <Link to="/register">立即注册</Link>
                            </div>
                        </Space>
                    </div>
                </Card>
            </div>
        </div>
    );
};

export default Login;