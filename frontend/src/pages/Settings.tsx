import React from 'react';
import { Card, Form, Input, Button, Switch, Select, Divider, Space, message, Modal } from 'antd';
import { LockOutlined, BellOutlined, EyeOutlined, GlobalOutlined, DeleteOutlined } from '@ant-design/icons';
import { useAuthStore } from '@/stores';

const { Option } = Select;
const { confirm } = Modal;

const Settings: React.FC = () => {
  const { logout } = useAuthStore();
  const [passwordForm] = Form.useForm();
  const [loading, setLoading] = React.useState(false);
  
  // 本地状态管理
  const [themeMode, setThemeMode] = React.useState<'light' | 'dark' | 'auto'>('light');
  const [primaryColor, setPrimaryColor] = React.useState('#1890ff');
  const [locale, setLocale] = React.useState('zh-CN');
  const [notificationDuration, setNotificationDuration] = React.useState(4.5);

  React.useEffect(() => {
    // 设置页面标题
    document.title = '系统设置 - Lemon Tea';
  }, []);

  const handlePasswordChange = async (values: any) => {
    if (values.newPassword !== values.confirmPassword) {
      message.error('两次输入的密码不一致');
      return;
    }

    try {
      setLoading(true);
      // 这里应该调用修改密码的API
      console.log('修改密码:', values);
      message.success('密码修改成功！');
      passwordForm.resetFields();
    } catch (error) {
      message.error('密码修改失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteAccount = () => {
    confirm({
      title: '确认删除账户',
      content: '删除账户后，所有数据将无法恢复。确定要继续吗？',
      okText: '确认删除',
      okType: 'danger',
      cancelText: '取消',
      onOk: async () => {
        try {
          // 这里应该调用删除账户的API
          console.log('删除账户');
          message.success('账户删除成功');
          await logout();
        } catch (error) {
          message.error('账户删除失败');
        }
      },
    });
  };

  const themeOptions = [
    { label: '浅色模式', value: 'light' },
    { label: '深色模式', value: 'dark' },
    { label: '跟随系统', value: 'auto' },
  ];

  const languageOptions = [
    { label: '简体中文', value: 'zh-CN' },
    { label: 'English', value: 'en-US' },
  ];

  const colorOptions = [
    { label: '蓝色', value: '#1890ff' },
    { label: '绿色', value: '#52c41a' },
    { label: '橙色', value: '#fa8c16' },
    { label: '红色', value: '#f5222d' },
    { label: '紫色', value: '#722ed1' },
    { label: '青色', value: '#13c2c2' },
  ];

  return (
    <div style={{ padding: '24px', maxWidth: '800px', margin: '0 auto' }}>
      <Space direction="vertical" size="large" style={{ width: '100%' }}>
        {/* 外观设置 */}
        <Card title={<><EyeOutlined /> 外观设置</>}>
          <Space direction="vertical" size="middle" style={{ width: '100%' }}>
            <div>
              <label style={{ display: 'block', marginBottom: '8px', fontWeight: 500 }}>
                主题模式
              </label>
              <Select
                value={themeMode}
                onChange={setThemeMode}
                style={{ width: '200px' }}
                options={themeOptions}
              />
            </div>
            
            <div>
              <label style={{ display: 'block', marginBottom: '8px', fontWeight: 500 }}>
                主题色
              </label>
              <Select
                value={primaryColor}
                onChange={setPrimaryColor}
                style={{ width: '200px' }}
              >
                {colorOptions.map(option => (
                  <Option key={option.value} value={option.value}>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                      <div
                        style={{
                          width: '16px',
                          height: '16px',
                          backgroundColor: option.value,
                          borderRadius: '2px',
                        }}
                      />
                      {option.label}
                    </div>
                  </Option>
                ))}
              </Select>
            </div>
            
            <div>
              <Button onClick={() => {
                setThemeMode('light');
                setPrimaryColor('#1890ff');
              }}>重置主题</Button>
            </div>
          </Space>
        </Card>

        {/* 语言设置 */}
        <Card title={<><GlobalOutlined /> 语言设置</>}>
          <div>
            <label style={{ display: 'block', marginBottom: '8px', fontWeight: 500 }}>
              界面语言
            </label>
            <Select
              value={locale}
              onChange={setLocale}
              style={{ width: '200px' }}
              options={languageOptions}
            />
          </div>
        </Card>

        {/* 通知设置 */}
        <Card title={<><BellOutlined /> 通知设置</>}>
          <Space direction="vertical" size="middle" style={{ width: '100%' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <span>桌面通知</span>
              <Switch defaultChecked />
            </div>
            
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <span>邮件通知</span>
              <Switch defaultChecked />
            </div>
            
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <span>短信通知</span>
              <Switch />
            </div>
            
            <div>
              <label style={{ display: 'block', marginBottom: '8px', fontWeight: 500 }}>
                通知显示时长（秒）
              </label>
              <Select
                value={notificationDuration}
                onChange={setNotificationDuration}
                style={{ width: '200px' }}
              >
                <Option value={3}>3秒</Option>
                <Option value={4.5}>4.5秒</Option>
                <Option value={6}>6秒</Option>
                <Option value={0}>不自动关闭</Option>
              </Select>
            </div>
          </Space>
        </Card>

        {/* 安全设置 */}
        <Card title={<><LockOutlined /> 安全设置</>}>
          <Space direction="vertical" size="large" style={{ width: '100%' }}>
            <div>
              <h4>修改密码</h4>
              <Form
                form={passwordForm}
                layout="vertical"
                onFinish={handlePasswordChange}
                style={{ maxWidth: '400px' }}
              >
                <Form.Item
                  label="当前密码"
                  name="currentPassword"
                  rules={[{ required: true, message: '请输入当前密码' }]}
                >
                  <Input.Password placeholder="请输入当前密码" />
                </Form.Item>
                
                <Form.Item
                  label="新密码"
                  name="newPassword"
                  rules={[
                    { required: true, message: '请输入新密码' },
                    { min: 6, message: '密码至少6个字符' },
                  ]}
                >
                  <Input.Password placeholder="请输入新密码" />
                </Form.Item>
                
                <Form.Item
                  label="确认新密码"
                  name="confirmPassword"
                  rules={[{ required: true, message: '请确认新密码' }]}
                >
                  <Input.Password placeholder="请再次输入新密码" />
                </Form.Item>
                
                <Form.Item>
                  <Button type="primary" htmlType="submit" loading={loading}>
                    修改密码
                  </Button>
                </Form.Item>
              </Form>
            </div>
            
            <Divider />
            
            <div>
              <h4 style={{ color: '#ff4d4f' }}>危险操作</h4>
              <p style={{ color: '#8c8c8c', marginBottom: '16px' }}>
                删除账户后，所有数据将无法恢复，请谨慎操作。
              </p>
              <Button
                type="primary"
                danger
                icon={<DeleteOutlined />}
                onClick={handleDeleteAccount}
              >
                删除账户
              </Button>
            </div>
          </Space>
        </Card>
      </Space>
    </div>
  );
};

export default Settings;